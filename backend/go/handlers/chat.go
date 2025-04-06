package handlers

import (
	"net/http"
	"time"

	"speakup/config"
	"speakup/middlewares"
	"speakup/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary CRUD operations for chat

// CreateChat godoc
// @Summary Create a new chat
// @Description Creates a new chat for a user
// @Tags Chat
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param chat body models.Chat true "Chat object"
// @Success 201 {object} models.Chat "Created chat"
// @Failure 400 {object} object "Bad request"
// @Failure 500 {object} object "Internal server error"
// @Router /chat [post]
func CreateChat(c *gin.Context) {
	var chat models.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chat.ID = uuid.New().String()
	chat.UserID = middlewares.GetUserIDFromContext(c)
	chat.StartTime = time.Now().Format(time.RFC3339)

	db := config.GetMongoClient()
	collection := db.Database("speakup").Collection("chats")
	_, err := collection.InsertOne(c, chat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
		return
	}

	c.JSON(http.StatusCreated, chat)
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
func GetChatById(c *gin.Context) {
	id := c.Param("id")
	db := config.GetMongoClient()
	collection := db.Database("speakup").Collection("chats")
	var chat models.Chat
	err := collection.FindOne(c, map[string]string{"id": id}).Decode(&chat)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat not found"})
		return
	}

	c.JSON(http.StatusOK, chat)
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
func GetChats(c *gin.Context) {
	db := config.GetMongoClient()
	collection := db.Database("speakup").Collection("chats")
	cursor, err := collection.Find(c, map[string]string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chats"})
		return
	}

	var chats []models.Chat
	if err := cursor.All(c, &chats); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chats"})
		return
	}

	c.JSON(http.StatusOK, chats)
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
func UpdateChat(c *gin.Context) {
	id := c.Param("id")
	var chat models.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.GetMongoClient()
	collection := db.Database("speakup").Collection("chats")
	_, err := collection.UpdateOne(c, map[string]string{"id": id}, chat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update chat"})
		return
	}

	c.JSON(http.StatusOK, chat)
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
func DeleteChat(c *gin.Context) {
	id := c.Param("id")
	db := config.GetMongoClient()
	collection := db.Database("speakup").Collection("chats")
	_, err := collection.DeleteOne(c, map[string]string{"id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete chat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat deleted successfully"})
}

// Special operations for chat

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
func GetChatsByUserId(c *gin.Context) {
	id := middlewares.GetUserIDFromContext(c)
	db := config.GetMongoClient()
	collection := db.Database("speakup").Collection("chats")
	cursor, err := collection.Find(c, map[string]string{"user_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chats"})
		return
	}

	var chats []models.Chat
	if err := cursor.All(c, &chats); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chats": chats})
}
