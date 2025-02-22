package handlers

import (
	"net/http"
	"time"

	"speakup/models"
	"speakup/ai"

	"github.com/gin-gonic/gin"

	"speakup/config"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// Send handles the creation of a new message
func Send(c *gin.Context) {
	
	client := config.GetMongoClient()

	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := client.Database("speakup").Collection("messages")

	message.ID = uuid.ClockSequence()

	message.Timestamp = time.Now().String()

	_, err := collection.InsertOne(c, bson.M{
		"id":      message.ID,
		"sender_id": message.SenderID,
		"chat_id": message.ChatID,
		"content": message.Content,
		"timestamp": message.Timestamp,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create message"})
		return
	}

	response, err := ai.GetAIResponse(message.Content)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get AI response"})
        return
    }

	c.JSON(http.StatusOK, gin.H{"message": message, "response": response})

}