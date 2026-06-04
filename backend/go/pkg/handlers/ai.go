package handlers

import (
	"context"
	"fmt"
	"net/http"

	"speakup/pkg/adapters/ai"
	"speakup/pkg/config"
	"speakup/pkg/middlewares"
	"speakup/pkg/models"
	"speakup/pkg/planlimits"

	"github.com/gin-gonic/gin"
)

func GetAIUsage(c *gin.Context) {
	userID := middlewares.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	pro, err := planlimits.IsProUser(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load subscription status"})
		return
	}

	limit := planlimits.GetFreeDailyLimit()
	used := int64(0)
	remaining := int64(0)
	if !pro {
		used, err = planlimits.GetUsageCount(c, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load usage"})
			return
		}
		remaining = limit - used
		if remaining < 0 {
			remaining = 0
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"is_pro":      pro,
		"daily_limit": limit,
		"used_today":  used,
		"remaining":   remaining,
	})
}

func EnforcePlanLimits(c *gin.Context) bool {
	userID := middlewares.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return false
	}

	pro, err := planlimits.IsProUser(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load subscription status"})
		return false
	}

	if pro {
		return true
	}

	limit := planlimits.GetFreeDailyLimit()
	used, err := planlimits.GetUsageCount(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enforce plan limits"})
		return false
	}
	if used >= limit {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Limite diário do plano Free atingido. Faça upgrade para continuar.",
		})
		return false
	}

	return true
}

func GenerateResponseDialog(c *gin.Context) {
	if !EnforcePlanLimits(c) {
		return
	}

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
	if !EnforcePlanLimits(c) {
		return
	}

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
	if !EnforcePlanLimits(c) {
		return
	}

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
	if !EnforcePlanLimits(c) {
		return
	}

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
