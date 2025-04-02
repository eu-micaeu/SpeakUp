package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"speakup/connectors"
	"speakup/middlewares"

	"github.com/gin-gonic/gin"
)

// GenerateResponse generates a response for the AI
func GenerateResponseDialog(c *gin.Context) {

	promptPath := filepath.Join("prompts", "promptDialog.txt")
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

	connector := connectors.NewOpenAIConnector()

	// Generate a response for the dialogue
	dialogueResp, err := connector.GenerateResponse(context.Background(), prompt + request.Message + "Answer me in this language: " + middlewares.GetLanguageFromContext(c))
	fmt.Println(dialogueResp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": dialogueResp})
}

// GenerateResponseCorrection generates a correction for the AI
func GenerateResponseCorrection(c *gin.Context) {

	promptPath := filepath.Join("prompts", "promptCorrectionEN.txt")
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

	connector := connectors.NewOpenAIConnector()

	fmt.Println(middlewares.GetLanguageFromContext(c))
	fmt.Println(middlewares.GetUserIDFromContext(c))

	// Generate a correction for the dialogue
	correctionResp, err := connector.GenerateResponse(context.Background(), prompt + request.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": correctionResp})
}

// GenerateResponseTranslate generates a translation for the AI
func GenerateResponseTranslate(c *gin.Context) {
	var request struct {
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	connector := connectors.NewOpenAIConnector()

	// Generate a translation for the dialogue
	translateResp, err := connector.GenerateResponse(context.Background(), "Please translate the following text to Portuguese: "+request.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": translateResp})
}

// GenerateResponseTopic generates a response for the AI
func GenerateResponseTopic(c *gin.Context) {
	var request struct {
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	connector := connectors.NewOpenAIConnector()

	// Generate a response for the topic
	topicResp, err := connector.GenerateResponse(context.Background(), "Please generate a topic for the following text, generate with 2 words only: " + request.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": topicResp})
}