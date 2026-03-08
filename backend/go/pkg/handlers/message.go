package handlers

import (
	"net/http"
	"strings"
	"time"

	"speakup/pkg/middlewares"
	"speakup/pkg/models"
	"speakup/pkg/planlimits"
	"speakup/pkg/repositories"
	"speakup/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MessageHandler struct {
	Repo repositories.MessageRepository
}

func NewMessageHandler(repo repositories.MessageRepository) *MessageHandler {
	return &MessageHandler{Repo: repo}
}

// CreateMessage creates a new message
// @Summary Create a new message
// @Description Create a new message with the provided information
// @Tags Message
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param message body models.Message true "Message object containing chat_id, content, sender, and type"
// @Success 201 {object} models.Message
// @Failure 400 "Campo obrigatório faltando"
// @Failure 401 "Não autorizado"
// @Failure 500 "Erro ao criar mensagem"
// @Router /message [post]
func (h *MessageHandler) CreateMessage(c *gin.Context) {
	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if strings.ToLower(message.Sender) == "user" {
		userID := middlewares.GetUserIDFromContext(c)
		if userID == "" {
			utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
			return
		}

		pro, err := planlimits.IsProUser(c, userID)
		if err != nil {
			utils.RespondWithError(c, http.StatusInternalServerError, "Failed to load subscription status")
			return
		}

		if !pro {
			limit := planlimits.GetFreeDailyLimit()
			allowed, err := planlimits.CheckAndIncrementUsage(c, userID, limit)
			if err != nil {
				utils.RespondWithError(c, http.StatusInternalServerError, "Failed to enforce plan limits")
				return
			}
			if !allowed {
				utils.RespondWithError(c, http.StatusTooManyRequests, "Limite diário do plano Free atingido. Faça upgrade para continuar.")
				return
			}
		}
	}

	message.ID = uuid.New().String()
	message.CreatedAt = time.Now().Format(time.RFC3339)

	if err := h.Repo.Create(c, message); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create message")
		return
	}

	utils.RespondWithJSON(c, http.StatusCreated, message)
}

// GetMessageById gets a message by ID
// @Summary Get a message by ID
// @Description Retrieve a message using its ID
// @Tags Message
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Message ID"
// @Success 200 {object} models.Message
// @Failure 401 "Não autorizado"
// @Failure 404 "Mensagem não encontrada"
// @Router /message/{id} [get]
func (h *MessageHandler) GetMessageById(c *gin.Context) {
	id := c.Param("id")
	message, err := h.Repo.FindByID(c, id)
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Message not found")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, message)
}

// GetMessages gets all messages
// @Summary Get all messages
// @Description Retrieve a list of all messages
// @Tags Message
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} models.Message
// @Failure 401 "Não autorizado"
// @Failure 500 "Erro ao buscar mensagens"
// @Router /message [get]
func (h *MessageHandler) GetMessages(c *gin.Context) {
	messages, err := h.Repo.FindAll(c)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get messages")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, messages)
}

// UpdateMessage updates a message by ID
// @Summary Update a message
// @Description Update an existing message with new information
// @Tags Message
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Message ID"
// @Param message body models.Message true "Message object containing content and type"
// @Success 200 "Mensagem atualizada com sucesso"
// @Failure 400 "Campo obrigatório faltando"
// @Failure 401 "Não autorizado"
// @Failure 500 "Erro ao atualizar mensagem"
// @Router /message/{id} [put]
func (h *MessageHandler) UpdateMessage(c *gin.Context) {
	id := c.Param("id")
	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	updateDoc := bson.M{}
	if message.Content != "" {
		updateDoc["content"] = message.Content
	}
	if message.Type != "" {
		updateDoc["type"] = message.Type
	}

	if err := h.Repo.Update(c, id, updateDoc); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update message")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{"message": "Message updated successfully"})
}

// DeleteMessage deletes a message by ID
// @Summary Delete a message
// @Description Remove an existing message from the system
// @Tags Message
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Message ID"
// @Success 200 "Mensagem deletada com sucesso"
// @Failure 401 "Não autorizado"
// @Failure 500 "Erro ao deletar mensagem"
// @Router /message/{id} [delete]
func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	id := c.Param("id")
	if err := h.Repo.Delete(c, id); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete message")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, gin.H{"message": "Message deleted successfully"})
}

// GetMessagesByChatId gets all messages from a specific chat
// @Summary Get messages by chat ID
// @Description Retrieve all messages belonging to a specific chat
// @Tags Message
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Chat ID"
// @Success 200 {array} models.Message
// @Failure 401 "Não autorizado"
// @Failure 500 "Erro ao buscar mensagens"
// @Router /message/chat/{id} [get]
func (h *MessageHandler) GetMessagesByChatId(c *gin.Context) {
	chatId := c.Param("id")
	messages, err := h.Repo.FindAllByChatID(c, chatId)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get messages")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, messages)
}
