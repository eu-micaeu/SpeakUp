package handlers

import (
	"net/http"
	"time"

	"speakup/models"

	"speakup/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// StartChat handles the creation of a new chat
func StartChat(c *gin.Context) {
    client := config.GetMongoClient()

    var chat models.Chat
    if err := c.ShouldBindJSON(&chat); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    collection := client.Database("speakup").Collection("chats")

    chat.ID = uuid.ClockSequence()

    chat.StartTime = time.Now().String()

    _, err := collection.InsertOne(c, bson.M{
        "id":              chat.ID,
        "user_id":         chat.UserID,
        "start_time":      chat.StartTime,
        "end_time":        chat.EndTime,
        "difficulty_level": chat.DifficultyLevel,
        "topic":           chat.Topic,
    })
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Chat created successfully", "chat": chat})
}

// EndChat handles the deletion of a chat
func EndChat(c *gin.Context) {
	
	client := config.GetMongoClient()

	var chat models.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := client.Database("speakup").Collection("chats")

	_, err := collection.DeleteOne(c, bson.M{
		"id": chat.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete chat"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat deleted successfully"})

}