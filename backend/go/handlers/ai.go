package handlers

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	openai "github.com/sashabaranov/go-openai"
)

// GenerateResponse generates a response for the AI
func GenerateResponseDialog(c *gin.Context) {
	// Ler o prompt de um arquivo externo
	promptPath := filepath.Join("prompts", "prompt1.txt")
	promptBytes, err := os.ReadFile(promptPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load prompt: " + err.Error()})
		return
	}
	prompt := string(promptBytes)

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
				Content: prompt + "Please pretend to be a friend and answer:" + request.Message,
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
				Content: "Please correct the spelling and grammar of the following text: " + request.Message,
			},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": correctionResp.Choices[0].Message.Content})
}
