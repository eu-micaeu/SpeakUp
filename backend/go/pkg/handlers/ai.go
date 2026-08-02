package handlers

import (
	"context"
	"fmt"
	"net/http"

	"speakup/pkg/adapters/ai"
	"speakup/pkg/config"
	"speakup/pkg/middlewares"
	"speakup/pkg/models"

	"github.com/gin-gonic/gin"
)

func GetAIUsage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"is_pro":      true,
		"daily_limit": -1,
		"used_today":  0,
		"remaining":   -1,
	})
}

func GenerateResponseDialog(c *gin.Context) {
	var request struct {
		Message string `json:"message"`
		ChatID  string `json:"chat_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.GetMongoClient()
	collection := db.Database("speakup").Collection("messages")
	cursor, err := collection.Find(c, map[string]string{"chat_id": request.ChatID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat history"})
		return
	}

	var messages []models.Message
	if err := cursor.All(c, &messages); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat history"})
		return
	}

	language := middlewares.GetLanguageFromContext(c)
	level := middlewares.GetLevelFromContext(c)
	aiProvider := c.GetHeader("X-AI-Provider")
	if aiProvider == "" {
		aiProvider = "gemini"
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
		aiProvider = "gemini"
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
		aiProvider = "gemini"
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
		aiProvider = "gemini"
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

