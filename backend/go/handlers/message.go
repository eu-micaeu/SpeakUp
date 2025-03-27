package handlers

import (
	"net/http"
	"time"

	"speakup/models"
	"speakup/config"

	"go.mongodb.org/mongo-driver/bson"
	

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CRUD operations for messages

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

// GetMessageById is a handler function that gets a message by ID
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

// GetMessages is a handler function that gets all messages
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

// UpdateMessage is a handler function that updates a message
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

// DeleteMessage is a handler function that deletes a message
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

// GetMessagesByChatId is a handler function that gets all messages by chat ID
func GetMessagesByChatId(c *gin.Context) {
	chatId := c.Param("id")
	db := config.GetMongoClient()
	collection := db.Database("speakup").Collection("messages")
	cursor, err := collection.Find(c, map[string]string{"chatid": chatId})
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