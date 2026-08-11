package handlers

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"speakup/pkg/adapters/ai"
	"speakup/pkg/models"
	"speakup/pkg/repositories"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type FlashcardHandler struct {
	Repo repositories.FlashcardRepository
}

func NewFlashcardHandler(repo repositories.FlashcardRepository) *FlashcardHandler {
	return &FlashcardHandler{Repo: repo}
}

func (h *FlashcardHandler) GenerateFlashcard(c *gin.Context) {
	var req struct {
		Term            string `json:"term" binding:"required"`
		ContextSentence string `json:"context_sentence"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	aiProvider := c.GetHeader("X-AI-Provider")
	if aiProvider == "" {
		aiProvider = os.Getenv("AI_PROVIDER")
	}
	if aiProvider == "" {
		aiProvider = "ollama"
	}
	ctx := context.WithValue(c.Request.Context(), "aiProvider", aiProvider)

	genData, err := ai.GetFlashcardGenerateResponse(ctx, req.Term, req.ContextSentence)
	if err != nil {
		// Fallback to translation engine if prompt generation fails
		trans, errTrans := ai.GetTranslationResponse(ctx, req.Term)
		if errTrans == nil && trans != "" {
			genData.Back = trans
			genData.Explanation = "Tradução automática."
		} else {
			genData.Back = req.Term
			genData.Explanation = "Tradução/Explicação automática."
		}
	}
	genData.Back = ai.FormatSingleWordTranslation(genData.Back)

	contextSent := req.ContextSentence
	if genData.ContextSentence != "" {
		contextSent = genData.ContextSentence
	}

	c.JSON(http.StatusOK, gin.H{
		"front":            req.Term,
		"back":             genData.Back,
		"context_sentence": contextSent,
		"explanation":      genData.Explanation,
	})
}

func (h *FlashcardHandler) CreateFlashcard(c *gin.Context) {
	var fc models.Flashcard
	if err := c.ShouldBindJSON(&fc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		userID = fc.UserID
	}
	if userID == "" {
		userID = "default_user"
	}
	fc.UserID = userID
	fc.Front = strings.TrimSpace(fc.Front)
	fc.Back = ai.FormatSingleWordTranslation(fc.Back)

	// Check if a flashcard for this word already exists for the user
	if existing, err := h.Repo.FindByUserAndFront(c.Request.Context(), userID, fc.Front); err == nil && existing.ID != "" {
		c.JSON(http.StatusConflict, gin.H{"error": "Já existe um flashcard cadastrado para a palavra '" + fc.Front + "'."})
		return
	}

	created, err := h.Repo.Create(c.Request.Context(), fc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar flashcard: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *FlashcardHandler) GetFlashcards(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		userID = c.Query("user_id")
	}
	if userID == "" {
		userID = "default_user"
	}

	dueOnly := c.Query("due") == "true"

	cards, err := h.Repo.FindByUser(c.Request.Context(), userID, dueOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar flashcards: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, cards)
}

func (h *FlashcardHandler) ReviewFlashcard(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID é obrigatório"})
		return
	}

	var req struct {
		Rating int `json:"rating" binding:"required"` // 1: Again, 2: Hard, 3: Good, 4: Easy
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Avaliação inválida (use 1 a 4): " + err.Error()})
		return
	}

	fc, err := h.Repo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Flashcard não encontrado"})
		return
	}

	// Algorithm SM-2 (SuperMemo-2)
	// Map rating 1..4 to score q 1..5: 1 -> 1.0, 2 -> 3.0, 3 -> 4.0, 4 -> 5.0
	var q float64
	switch req.Rating {
	case 1:
		q = 1.0 // Again / Errei
	case 2:
		q = 3.0 // Hard / Difícil
	case 3:
		q = 4.0 // Good / Bom
	case 4:
		q = 5.0 // Easy / Fácil
	default:
		q = 4.0
	}

	easeFactor := fc.EaseFactor
	if easeFactor < 1.3 {
		easeFactor = 2.5
	}
	repetitions := fc.Repetitions
	interval := fc.Interval

	var newInterval int
	var newRepetitions int

	if q >= 3.0 {
		if repetitions == 0 {
			newInterval = 1
		} else if repetitions == 1 {
			newInterval = 6
		} else {
			newInterval = int(math.Round(float64(interval) * easeFactor))
		}
		newRepetitions = repetitions + 1
	} else {
		newRepetitions = 0
		newInterval = 1
	}

	newEaseFactor := easeFactor + (0.1 - (5.0-q)*(0.08+(5.0-q)*0.02))
	if newEaseFactor < 1.3 {
		newEaseFactor = 1.3
	}

	now := time.Now()
	nextReview := now.AddDate(0, 0, newInterval)

	updateDoc := bson.M{
		"ease_factor":      newEaseFactor,
		"interval":         newInterval,
		"repetitions":      newRepetitions,
		"next_review":      nextReview,
		"last_reviewed_at": now,
	}

	if err := h.Repo.Update(c.Request.Context(), id, updateDoc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar revisão: " + err.Error()})
		return
	}

	fc.EaseFactor = newEaseFactor
	fc.Interval = newInterval
	fc.Repetitions = newRepetitions
	fc.NextReview = nextReview
	fc.LastReviewedAt = &now

	c.JSON(http.StatusOK, fc)
}

func (h *FlashcardHandler) DeleteFlashcard(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID é obrigatório"})
		return
	}

	if err := h.Repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar flashcard: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Flashcard removido com sucesso"})
}

func (h *FlashcardHandler) GenerateDailyBatch(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		userID = "default_user"
	}

	aiProvider := c.GetHeader("X-AI-Provider")
	if aiProvider == "" {
		aiProvider = os.Getenv("AI_PROVIDER")
	}
	if aiProvider == "" {
		aiProvider = "ollama"
	}
	ctx := context.WithValue(c.Request.Context(), "aiProvider", aiProvider)

	// Fetch existing flashcards for user to avoid duplicate terms
	existingCards, err := h.Repo.FindByUser(ctx, userID, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao carregar flashcards existentes: " + err.Error()})
		return
	}

	existingMap := make(map[string]bool)
	var existingTerms []string
	for _, card := range existingCards {
		termLower := strings.ToLower(strings.TrimSpace(card.Front))
		existingMap[termLower] = true
		existingTerms = append(existingTerms, card.Front)
	}

	batchItems, err := ai.GetBatchFlashcardsResponse(ctx, existingTerms, 20)
	if err != nil || len(batchItems) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar lote de flashcards com IA: " + err.Error()})
		return
	}

	var createdList []models.Flashcard
	now := time.Now()
	for _, item := range batchItems {
		termClean := strings.TrimSpace(item.Front)
		if termClean == "" {
			continue
		}
		termLower := strings.ToLower(termClean)
		if existingMap[termLower] {
			continue
		}
		existingMap[termLower] = true

		fc := models.Flashcard{
			UserID:          userID,
			Front:           termClean,
			Back:            ai.FormatSingleWordTranslation(item.Back),
			ContextSentence: item.ContextSentence,
			Explanation:     item.Explanation,
			EaseFactor:      2.5,
			Interval:        1,
			Repetitions:     0,
			CreatedAt:       now,
			NextReview:      now,
		}

		created, err := h.Repo.Create(ctx, fc)
		if err == nil {
			createdList = append(createdList, created)
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("%d novos flashcards gerados com sucesso!", len(createdList)),
		"count":   len(createdList),
		"cards":   createdList,
	})
}
