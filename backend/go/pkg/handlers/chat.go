package handlers

import (
	"net/http"
	"time"

	"speakup/pkg/middlewares"
	"speakup/pkg/models"
	"speakup/pkg/repositories"
	"speakup/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatHandler struct {
	Repo repositories.ChatRepository
}

func NewChatHandler(repo repositories.ChatRepository) *ChatHandler {
	return &ChatHandler{Repo: repo}
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

	utils.RespondWithJSON(c, http.StatusOK, gin.H{"message": "Chat deleted successfully"})
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
