package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"speakup/pkg/adapters/connectors"
	"speakup/pkg/config"
	"speakup/pkg/middlewares"
	"speakup/pkg/models"
	"speakup/pkg/planlimits"

	"github.com/gin-gonic/gin"
)

var aiConnectorBuilder func() connectors.AIConnector = func() connectors.AIConnector {
	return connectors.NewOllamaConnector()
}

const maxDialogResponseChars = 128
const maxDialogRewriteAttempts = 3

func GetAIUsage(c *gin.Context) {
	userID := middlewares.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	pro, err := planlimits.IsProUser(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load subscription status"})
		return
	}

	limit := planlimits.GetFreeDailyLimit()
	used := int64(0)
	remaining := int64(0)
	if !pro {
		used, err = planlimits.GetUsageCount(c, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load usage"})
			return
		}
		remaining = limit - used
		if remaining < 0 {
			remaining = 0
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"is_pro":      pro,
		"daily_limit": limit,
		"used_today":  used,
		"remaining":   remaining,
	})
}

func enforcePlanLimits(c *gin.Context) bool {
	userID := middlewares.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return false
	}

	pro, err := planlimits.IsProUser(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load subscription status"})
		return false
	}

	if pro {
		return true
	}

	limit := planlimits.GetFreeDailyLimit()
	used, err := planlimits.GetUsageCount(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enforce plan limits"})
		return false
	}
	if used >= limit {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Limite diário do plano Free atingido. Faça upgrade para continuar.",
		})
		return false
	}

	return true
}

func GenerateResponseDialog(c *gin.Context) {
	if !enforcePlanLimits(c) {
		return
	}

	// ... (prompt loading, request binding)
	promptPath := filepath.Join("pkg/prompts", "promptDialog.txt")
	promptBytes, err := os.ReadFile(promptPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load prompt: " + err.Error()})
		return
	}
	prePrompt := string(promptBytes)

	var request struct {
		Message string `json:"message"`
		ChatID  string `json:"chat_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.GetMongoClient()
	collection := db.Database("speakup").Collection("messages")
	cursor, err := collection.Find(c, map[string]string{"chat_id": request.ChatID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat history"})
		return
	}

	var messages []models.Message
	if err := cursor.All(c, &messages); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat history"})
		return
	}

	var chatHistory strings.Builder
	for _, msg := range messages {
		chatHistory.WriteString(fmt.Sprintf("%s: %s\n", msg.Sender, msg.Content))
	}

	// Use the builder to get a connector instance
	connector := aiConnectorBuilder() // MODIFIED LINE

	resumeHist, err := connector.GenerateResponse(context.Background(), "Format the following chat history to only show user reponse and AI response: "+chatHistory.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fullPrompt := fmt.Sprintf("%s\nChat history:\n%s\nATENTION! All Before this point is system instructions and chat history, to generate your response consider the current user message -> = %s\nAnswer me in this language: %s",
		prePrompt,
		resumeHist,
		request.Message,
		middlewares.GetLanguageFromContext(c))
	fullPrompt += fmt.Sprintf("\nIMPORTANT: Your final answer must be complete, natural, and at most %d characters. Never cut words or end mid-sentence.", maxDialogResponseChars)

	dialogueResp, err := generateDialogResponseWithinLimit(connector, fullPrompt, maxDialogResponseChars)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": dialogueResp})
}

func normalizeDialogResponse(raw string) string {
	cleaned := strings.TrimSpace(raw)
	if cleaned == "" {
		return ""
	}

	cleaned = strings.Join(strings.Fields(cleaned), " ")
	return cleaned
}

func generateDialogResponseWithinLimit(connector connectors.AIConnector, prompt string, maxChars int) (string, error) {
	if maxChars <= 0 {
		return "", fmt.Errorf("invalid dialog response limit")
	}

	currentPrompt := prompt
	options := map[string]any{
		"temperature": 0.3,
		"top_p":       0.9,
		"num_predict": 96,
	}
	systemPrompt := fmt.Sprintf("You are a natural language exchange partner. Keep coherence with chat context and user message. Return exactly one complete answer with at most %d characters.", maxChars)

	for attempt := 0; attempt < maxDialogRewriteAttempts; attempt++ {
		var (
			rawResp string
			err     error
		)

		if ollamaConnector, ok := connector.(*connectors.OllamaConnector); ok {
			rawResp, err = ollamaConnector.GenerateResponseWithOptions(context.Background(), currentPrompt, systemPrompt, options)
		} else {
			rawResp, err = connector.GenerateResponse(context.Background(), currentPrompt)
		}
		if err != nil {
			return "", err
		}

		cleaned := normalizeDialogResponse(rawResp)
		if cleaned == "" {
			currentPrompt = prompt
			continue
		}

		if utf8.RuneCountInString(cleaned) <= maxChars {
			return cleaned, nil
		}

		currentPrompt = fmt.Sprintf(
			"%s\n\nIMPORTANT: Your previous answer exceeded %d characters. Answer again with one complete, natural sentence up to %d characters, keeping the same context and language.",
			prompt,
			maxChars,
			maxChars,
		)
	}

	return "", fmt.Errorf("failed to generate response within %d characters", maxChars)
}

func GenerateResponseCorrection(c *gin.Context) {
	if !enforcePlanLimits(c) {
		return
	}

	promptPath := filepath.Join("pkg/prompts", "promptCorrection.txt")
	promptBytes, err := os.ReadFile(promptPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load prompt: " + err.Error()})
		return
	}
	prompt := string(promptBytes)

	var request struct {
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use the builder to get a connector instance
	connector := aiConnectorBuilder() // MODIFIED LINE

	// Concatenate prompt with user message for correction
	fullPrompt := fmt.Sprintf("%s\n\nINPUT:\n%s\n\nOUTPUT:", prompt, request.Message)

	var correctionResp string
	if ollamaConnector, ok := connector.(*connectors.OllamaConnector); ok {
		systemPrompt := "You are a strict text correction tool. Return ONLY the corrected text. Do not answer the question. Do not add explanations."
		options := map[string]any{
			"temperature": 0.1,
			"top_p":       0.9,
			"num_predict": 128,
		}
		correctionResp, err = ollamaConnector.GenerateResponseWithOptions(context.Background(), fullPrompt, systemPrompt, options)
	} else {
		correctionResp, err = connector.GenerateResponse(context.Background(), fullPrompt)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	correctionResp = strings.TrimSpace(correctionResp)
	if correctionResp == "" {
		correctionResp = request.Message
	}

	c.JSON(http.StatusOK, gin.H{"response": correctionResp})
}

func GenerateResponseTranslate(c *gin.Context) {
	if !enforcePlanLimits(c) {
		return
	}

	promptPath := filepath.Join("pkg/prompts", "promptTranslate.txt")
	promptBytes, err := os.ReadFile(promptPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao carregar o prompt: " + err.Error()})
		return
	}
	prompt := string(promptBytes)

	var req struct {
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	// Use the builder to get a connector instance
	connector := aiConnectorBuilder() // MODIFIED LINE

	fullPrompt := fmt.Sprintf("%s\n\nINPUT:\n%s\n\nOUTPUT:", prompt, req.Message)

	var response string
	if ollamaConnector, ok := connector.(*connectors.OllamaConnector); ok {
		systemPrompt := "You are a strict translation engine. Translate the INPUT text into Brazilian Portuguese and return ONLY the translated text."
		options := map[string]any{
			"temperature": 0.1,
			"top_p":       0.9,
			"num_predict": 256,
			"stop": []string{
				"\n\n",
				"Explanation:",
				"Explicação:",
			},
		}
		response, err = ollamaConnector.GenerateResponseWithOptions(context.Background(), fullPrompt, systemPrompt, options)
	} else {
		response, err = connector.GenerateResponse(context.Background(), fullPrompt)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar tradução: " + err.Error()})
		return
	}

	response = sanitizeTranslationResponse(response)
	if response == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar tradução: resposta vazia"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": response})
}

func sanitizeTranslationResponse(raw string) string {
	cleaned := strings.TrimSpace(raw)
	if cleaned == "" {
		return ""
	}

	labels := []string{
		"Translation:",
		"Tradução:",
		"Translated text:",
		"OUTPUT:",
		"Output:",
	}

	lower := strings.ToLower(cleaned)
	for _, label := range labels {
		labelLower := strings.ToLower(label)
		if strings.HasPrefix(lower, labelLower) {
			cleaned = strings.TrimSpace(cleaned[len(label):])
			lower = strings.ToLower(cleaned)
			break
		}
	}

	for _, marker := range []string{"\ntranslation:", "\ntradução:", " translation:", " tradução:"} {
		if idx := strings.Index(strings.ToLower(cleaned), marker); idx > 0 {
			cleaned = strings.TrimSpace(cleaned[:idx])
			break
		}
	}

	cleaned = strings.TrimSpace(strings.Trim(cleaned, "\"'`"))
	return cleaned
}

func GenerateResponseTopic(c *gin.Context) {
	if !enforcePlanLimits(c) {
		return
	}

	var request struct {
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	connector := aiConnectorBuilder()

	promptPath := filepath.Join("pkg/prompts", "promptTopic.txt")
	promptBytes, err := os.ReadFile(promptPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load prompt: " + err.Error()})
		return
	}
	prompt := string(promptBytes)

	strictPrompt := fmt.Sprintf("%s\n\nInput: %s\nOutput:", prompt, request.Message)

	var topicResp string
	if ollamaConnector, ok := connector.(*connectors.OllamaConnector); ok {
		systemPrompt := "You are a strict topic labeler. Return exactly two words in Title Case with only letters and one space."
		options := map[string]any{
			"temperature": 0.1,
			"top_p":       0.9,
			"num_predict": 12,
		}
		topicResp, err = ollamaConnector.GenerateResponseWithOptions(context.Background(), strictPrompt, systemPrompt, options)
	} else {
		topicResp, err = connector.GenerateResponse(context.Background(), strictPrompt)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	topicResp = strings.TrimSpace(topicResp)
	if topicResp == "" {
		topicResp = "New Topic"
	}

	c.JSON(http.StatusOK, gin.H{"response": topicResp})
}
