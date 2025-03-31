package handlers

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"speakup/connectors"

	"github.com/gin-gonic/gin"
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

	connector := connectors.NewGeminiConnector()

	// Generate a response for the dialogue
	dialogueResp, err := connector.GenerateResponse(context.Background(), prompt+"Please pretend to be a friend and answer:"+request.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": dialogueResp})
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

	connector := connectors.NewGeminiConnector()

	// Generate a correction for the dialogue
	correctionResp, err := connector.GenerateResponse(context.Background(), "Please correct the spelling and grammar of the following text, without prompts, correct the sentence/word and that's it!: "+request.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": correctionResp})
}
