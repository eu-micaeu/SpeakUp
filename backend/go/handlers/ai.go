package handlers

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	openai "github.com/sashabaranov/go-openai"
)

// GenerateResponse generates a response for the AI
func GenerateResponseDialog(c *gin.Context) {
    var request struct {
        Message string `json:"message"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    apiKey := os.Getenv("OPENAI_API_KEY")
    client := openai.NewClient(apiKey)

    // Generate a response for the dialogue
    dialogueResp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
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

    c.JSON(http.StatusOK, gin.H{"response": dialogueResp.Choices[0].Message.Content})
}

// GenerateResponseCorrection generates a correction for the AI
func GenerateResponseCorrection(c *gin.Context) {
    var request struct {
        Message string `json:"message"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    apiKey := os.Getenv("OPENAI_API_KEY")
    client := openai.NewClient(apiKey)

    // Generate a correction for the dialogue
    correctionResp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
        Model: "gpt-3.5-turbo",
        Messages: []openai.ChatCompletionMessage{
            {
                Role:    "user",
                Content: "(CORRIJA DE FORMA ORTOGRÁFICA) " + request.Message,
            },
        },
    })
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"response": correctionResp.Choices[0].Message.Content})
}