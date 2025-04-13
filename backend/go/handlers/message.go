package handlers

import (
	"net/http"
	"time"

	"speakup/config"
	"speakup/models"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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

// CreateMessage is a handler function that creates a new message
func CreateMessage(c *gin.Context) {
	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message.ID = uuid.New().String()
	message.CreatedAt = time.Now().Format(time.RFC3339)

	db := config.GetMongoClient()
	collection := db.Database("speakup").Collection("messages")
	_, err := collection.InsertOne(c, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create message"})
		return
	}

	c.JSON(http.StatusCreated, message)

}

// GetMessageById gets a message by ID
// @Summary Get a message by ID
// @Description Retrieve a message using its ID
// @Tags message
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Message ID"
// @Success 200 {object} models.Message
// @Failure 401 "Não autorizado"
// @Failure 404 "Mensagem não encontrada"
// @Router /message/{id} [get]
func GetMessageById(c *gin.Context) {
	id := c.Param("id")
	db := config.GetMongoClient()
	collection := db.Database("speakup").Collection("messages")
	var message models.Message
	err := collection.FindOne(c, map[string]string{"id": id}).Decode(&message)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	c.JSON(http.StatusOK, message)
}

// GetMessages gets all messages
// @Summary Get all messages
// @Description Retrieve a list of all messages
// @Tags message
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Success 200 {array} models.Message
// @Failure 401 "Não autorizado"
// @Failure 500 "Erro ao buscar mensagens"
// @Router /message [get]
func GetMessages(c *gin.Context) {
	db := config.GetMongoClient()
	collection := db.Database("speakup").Collection("messages")
	cursor, err := collection.Find(c, map[string]string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
		return
	}

	var messages []models.Message
	if err := cursor.All(c, &messages); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// UpdateMessage updates a message by ID
// @Summary Update a message
// @Description Update an existing message with new information
// @Tags message
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
func UpdateMessage(c *gin.Context) {
	id := c.Param("id")
	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db := config.GetMongoClient()
	collection := db.Database("speakup").Collection("messages")
	_, err := collection.UpdateOne(c, map[string]string{"id": id}, bson.M{"$set": message})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message updated successfully"})
}

// DeleteMessage deletes a message by ID
// @Summary Delete a message
// @Description Remove an existing message from the system
// @Tags message
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Message ID"
// @Success 200 "Mensagem deletada com sucesso"
// @Failure 401 "Não autorizado"
// @Failure 500 "Erro ao deletar mensagem"
// @Router /message/{id} [delete]
func DeleteMessage(c *gin.Context) {
	id := c.Param("id")
	db := config.GetMongoClient()
	collection := db.Database("speakup").Collection("messages")
	_, err := collection.DeleteOne(c, map[string]string{"id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
}

// Special handler functions

// GetMessagesByChatId gets all messages from a specific chat
// @Summary Get messages by chat ID
// @Description Retrieve all messages belonging to a specific chat
// @Tags message
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Chat ID"
// @Success 200 {array} models.Message
// @Failure 401 "Não autorizado"
// @Failure 500 "Erro ao buscar mensagens"
// @Router /message/chat/{id} [get]
func GetMessagesByChatId(c *gin.Context) {
	chatId := c.Param("id")
	db := config.GetMongoClient()
	collection := db.Database("speakup").Collection("messages")
	cursor, err := collection.Find(c, map[string]string{"chat_id": chatId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
		return
	}

	var messages []models.Message
	if err := cursor.All(c, &messages); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}
