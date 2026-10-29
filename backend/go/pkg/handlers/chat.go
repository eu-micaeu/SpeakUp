package handlers

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"speakup/pkg/adapters/ai"
	"speakup/pkg/middlewares"
	"speakup/pkg/models"
	"speakup/pkg/repositories"
	"speakup/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatHandler struct {
	Repo        repositories.ChatRepository
	MessageRepo repositories.MessageRepository
}

func NewChatHandler(repo repositories.ChatRepository, messageRepo repositories.MessageRepository) *ChatHandler {
	return &ChatHandler{
		Repo:        repo,
		MessageRepo: messageRepo,
	}
}


// CreateChat godoc
// @Summary Create a new chat
// @Description Creates a new chat for a user
// @Tags Chat
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param chat body object{user_id=string,topic=string} true "Chat object"
// @Success 201 {object} models.Chat "Created chat"
// @Failure 400 {object} object "Bad request"
// @Failure 500 {object} object "Internal server error"
// @Router /chat [post]
func (h *ChatHandler) CreateChat(c *gin.Context) {
	var chat models.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	chat.ID = uuid.New().String()
	chat.UserID = middlewares.GetUserIDFromContext(c)
	chat.StartTime = time.Now().Format(time.RFC3339)

	if err := h.Repo.Create(c, chat); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create chat")
		return
	}

	utils.RespondWithJSON(c, http.StatusCreated, chat)
}

// GetChatById godoc
// @Summary Get a chat by ID
// @Description Retrieves a chat by its ID
// @Tags Chat
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Chat ID"
// @Success 200 {object} models.Chat "Chat found"
// @Failure 404 {object} object "Chat not found"
// @Router /chat/{id} [get]
func (h *ChatHandler) GetChatById(c *gin.Context) {
	id := c.Param("id")
	chat, err := h.Repo.FindByID(c, id)
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Chat not found")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, chat)
}

// GetChats godoc
// @Summary Get all chats
// @Description Retrieves all chats
// @Tags Chat
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} models.Chat "List of chats"
// @Failure 500 {object} object "Internal server error"
// @Router /chat [get]
func (h *ChatHandler) GetChats(c *gin.Context) {
	chats, err := h.Repo.FindAll(c)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get chats")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, chats)
}

// UpdateChat godoc
// @Summary Update a chat
// @Description Updates an existing chat by ID
// @Tags Chat
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Chat ID"
// @Param chat body models.Chat true "Updated chat object"
// @Success 200 {object} models.Chat "Updated chat"
// @Failure 400 {object} object "Bad request"
// @Failure 500 {object} object "Internal server error"
// @Router /chat/{id} [put]
func (h *ChatHandler) UpdateChat(c *gin.Context) {
	id := c.Param("id")
	var chat models.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.Repo.Update(c, id, chat); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update chat")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, chat)
}

// DeleteChat godoc
// @Summary Delete a chat
// @Description Deletes a chat by ID
// @Tags Chat
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Chat ID"
// @Success 200 {object} object "Chat deleted successfully"
// @Failure 500 {object} object "Internal server error"
// @Router /chat/{id} [delete]
func (h *ChatHandler) DeleteChat(c *gin.Context) {
	id := c.Param("id")
	if err := h.Repo.Delete(c, id); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete chat")
		return
	}

	if h.MessageRepo != nil {
		_ = h.MessageRepo.DeleteByChatID(c, id)
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{"message": "Chat deleted successfully"})
}

// SendTurn orchestrates a complete conversation turn on the backend:
// 1. Creates the chat (and generates topic) if chat_id is empty
// 2. Proofreads and explains grammar corrections for the user's message
// 3. Saves the structured user message
// 4. Generates the AI tutor response using recent chat history
// 5. Translates the AI response into Brazilian Portuguese
// 6. Saves the structured AI message and returns everything in a single payload
func (h *ChatHandler) SendTurn(c *gin.Context) {
	var req struct {
		ChatID  string `json:"chat_id"`
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Mensagem inválida: "+err.Error())
		return
	}

	rawMessage := strings.TrimSpace(req.Message)
	if rawMessage == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "A mensagem não pode estar vazia")
		return
	}

	userID := middlewares.GetUserIDFromContext(c)
	if userID == "" {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	language := middlewares.GetLanguageFromContext(c)
	level := middlewares.GetLevelFromContext(c)

	aiProvider := c.GetHeader("X-AI-Provider")
	if aiProvider == "" {
		aiProvider = os.Getenv("AI_PROVIDER")
	}
	if aiProvider == "" {
		aiProvider = "ollama"
	}
	ctx := context.WithValue(c.Request.Context(), "aiProvider", aiProvider)

	var chat models.Chat
	var history []models.Message
	chatID := strings.TrimSpace(req.ChatID)

	if chatID == "" {
		topic, err := ai.GetTopicResponse(ctx, rawMessage)
		if err != nil || strings.TrimSpace(topic) == "" {
			topic = "New Conversation"
		}
		chat = models.Chat{
			ID:        uuid.New().String(),
			UserID:    userID,
			Topic:     topic,
			StartTime: time.Now().UTC().Format(time.RFC3339),
		}
		if err := h.Repo.Create(c, chat); err != nil {
			utils.RespondWithError(c, http.StatusInternalServerError, "Falha ao criar conversa")
			return
		}
		chatID = chat.ID
	} else {
		existingChat, err := h.Repo.FindByID(c, chatID)
		if err == nil {
			chat = existingChat
		} else {
			chat = models.Chat{ID: chatID, UserID: userID}
		}
		if h.MessageRepo != nil {
			history, _ = h.MessageRepo.FindRecentByChatID(c, chatID, 10)
		}
	}

	correctedText, explanationText, err := ai.GetCorrectionResponse(ctx, rawMessage)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Falha ao gerar correção: "+err.Error())
		return
	}

	userMsg := models.Message{
		ID:               uuid.New().String(),
		ChatID:           chatID,
		Content:          rawMessage,
		CorrectedContent: correctedText,
		Explanation:      explanationText,
		Sender:           "user",
		Type:             "request",
		CreatedAt:        time.Now().UTC().Format(time.RFC3339Nano),
	}

	if h.MessageRepo != nil {
		if err := h.MessageRepo.Create(c, userMsg); err != nil {
			utils.RespondWithError(c, http.StatusInternalServerError, "Falha ao salvar mensagem do usuário")
			return
		}
	}

	dialogInput := correctedText
	if strings.TrimSpace(dialogInput) == "" {
		dialogInput = rawMessage
	}

	dialogResp, err := ai.GetDialogResponse(ctx, dialogInput, history, language, level)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Falha ao gerar resposta do diálogo: "+err.Error())
		return
	}

	translationResp, err := ai.GetTranslationResponse(ctx, dialogResp)
	if err != nil {
		translationResp = ""
	}

	aiMsg := models.Message{
		ID:          uuid.New().String(),
		ChatID:      chatID,
		Content:     dialogResp,
		Translation: translationResp,
		Sender:      "ai",
		Type:        "response",
		CreatedAt:   time.Now().UTC().Format(time.RFC3339Nano),
	}

	if h.MessageRepo != nil {
		if err := h.MessageRepo.Create(c, aiMsg); err != nil {
			utils.RespondWithError(c, http.StatusInternalServerError, "Falha ao salvar resposta da IA")
			return
		}
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{
		"chat":         chat,
		"user_message": userMsg,
		"ai_message":   aiMsg,
	})
}


// GetChatsByUserId godoc
// @Summary Get user's chats
// @Description Retrieves all chats for the authenticated user
// @Tags Chat
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} object{chats=[]models.Chat} "List of user's chats"
// @Failure 500 {object} object "Internal server error"
// @Router /chat/user [get]
func (h *ChatHandler) GetChatsByUserId(c *gin.Context) {
	id := middlewares.GetUserIDFromContext(c)
	chats, err := h.Repo.FindAllByUserID(c, id)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get chats")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{"chats": chats})
}
