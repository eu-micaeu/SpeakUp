package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"speakup/pkg/adapters/ai"
	"speakup/pkg/config"
	"speakup/pkg/middlewares"
	"speakup/pkg/repositories"

	"github.com/gin-gonic/gin"
)

func GenerateResponseDialog(c *gin.Context) {
	var request struct {
		Message string `json:"message"`
		ChatID  string `json:"chat_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.GetMongoClient().Database(config.GetDBName())
	msgRepo := repositories.NewMongoMessageRepository(db)
	messages, err := msgRepo.FindRecentByChatID(c, request.ChatID, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat history"})
		return
	}

	language := middlewares.GetLanguageFromContext(c)
	level := middlewares.GetLevelFromContext(c)
	aiProvider := c.GetHeader("X-AI-Provider")
	if aiProvider == "" {
		aiProvider = os.Getenv("AI_PROVIDER")
	}
	if aiProvider == "" {
		aiProvider = "ollama"
	}
	ctx := context.WithValue(c.Request.Context(), "aiProvider", aiProvider)
	response, err := ai.GetDialogResponse(ctx, request.Message, messages, language, level)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": response})
}

func GenerateResponseCorrection(c *gin.Context) {
	var request struct {
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aiProvider := c.GetHeader("X-AI-Provider")
	if aiProvider == "" {
		aiProvider = os.Getenv("AI_PROVIDER")
	}
	if aiProvider == "" {
		aiProvider = "ollama"
	}
	ctx := context.WithValue(c.Request.Context(), "aiProvider", aiProvider)
	response, explanation, err := ai.GetCorrectionResponse(ctx, request.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response":    response,
		"explanation": explanation,
	})
}

func GenerateResponseTranslate(c *gin.Context) {
	var req struct {
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	aiProvider := c.GetHeader("X-AI-Provider")
	if aiProvider == "" {
		aiProvider = os.Getenv("AI_PROVIDER")
	}
	if aiProvider == "" {
		aiProvider = "ollama"
	}
	ctx := context.WithValue(c.Request.Context(), "aiProvider", aiProvider)
	response, err := ai.GetTranslationResponse(ctx, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar tradução: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": response})
}

func GenerateResponseTopic(c *gin.Context) {
	var request struct {
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aiProvider := c.GetHeader("X-AI-Provider")
	if aiProvider == "" {
		aiProvider = os.Getenv("AI_PROVIDER")
	}
	if aiProvider == "" {
		aiProvider = "ollama"
	}
	ctx := context.WithValue(c.Request.Context(), "aiProvider", aiProvider)
	response, err := ai.GetTopicResponse(ctx, request.Message)
	if err != nil {
		fmt.Printf("Error generating topic response: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": response})
}

