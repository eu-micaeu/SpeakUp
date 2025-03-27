package handlers

import (
	"context"
	"net/http"
	"os"

	"speakup/models"

	"github.com/gin-gonic/gin"

	"speakup/config"

	"github.com/google/uuid"
	openai "github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson"
)

// GenerateResponse generates a response for the AI
func GenerateResponse(c *gin.Context) {
	var request struct {
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "user",
				Content: request.Message,
			},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": resp.Choices[0].Message.Content})
}

// DetectSpellingErrors detects spelling errors in a message
func DetectSpellingErrors(c *gin.Context) {
	client := config.GetMongoClient()

	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := client.Database("speakup").Collection("messages")

	message.ID = uuid.New().String()


	_, err := collection.InsertOne(c, bson.M{
		"id":        message.ID,
		"chat_id":   message.ChatID,
		"content":   message.Content,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}
