package handlers

import (
	"net/http"
	"time"

	"speakup/models"
	"speakup/middlewares"
	"speakup/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CRUD operations for chat

// CreateChat creates a new chat
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

// GetChat gets a chat by ID
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

// GetChats gets all chats
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

// UpdateChat updates a chat by ID
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

// DeleteChat deletes a chat by ID
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

// GetChatsByUserId gets all chats by user ID
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