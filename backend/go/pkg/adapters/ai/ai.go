package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"speakup/pkg/adapters/connectors"
	"speakup/pkg/config"
	"speakup/pkg/middlewares"
	"speakup/pkg/models"
	"speakup/pkg/planlimits"

	"github.com/gin-gonic/gin"
)

var aiConnectorBuilder func() connectors.AIConnector = func() connectors.AIConnector {
	// Usa sempre o Ollama (local)
	return connectors.NewOllamaConnector()
}

// @Summary Retorna uso diário de IA
// @Description Retorna o consumo diário de créditos do usuário no plano Free
// @Tags AI
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Token de autenticação"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /ai/usage [get]
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

func enforcePlanLimits(c *gin.Context) bool {
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

// @Summary Gera uma resposta de diálogo usando IA
// @Description Gera uma resposta de diálogo contextual baseada na mensagem fornecida
// @Tags AI
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Token de autenticação"
// @Param request body object{message=string,chat_id=string} true "Chat object"
// @Success 200 {object} map[string]string "Resposta gerada com sucesso" example({"response":"Hi! I'm doing great, thank you for asking. How are you?"})
// @Failure 400 {object} map[string]string "Erro na requisição" example({"error":"Invalid request"})
// @Failure 500 {object} map[string]string "Erro interno do servidor" example({"error":"Internal server error"})
// @Router /ai/generate-response-dialog [post]
func GenerateResponseDialog(c *gin.Context) {
	if !enforcePlanLimits(c) {
		return
	}

	// ... (prompt loading, request binding)
	promptPath := filepath.Join("pkg/prompts", "promptDialog.txt")
	promptBytes, err := os.ReadFile(promptPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load prompt: " + err.Error()})
		return
	}
	prePrompt := string(promptBytes)

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

	var chatHistory strings.Builder
	for _, msg := range messages {
		chatHistory.WriteString(fmt.Sprintf("%s: %s\n", msg.Sender, msg.Content))
	}

	// Use the builder to get a connector instance
	connector := aiConnectorBuilder() // MODIFIED LINE

	resumeHist, err := connector.GenerateResponse(context.Background(), "Format the following chat history to only show user reponse and AI response: "+chatHistory.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fullPrompt := fmt.Sprintf("%s\nChat history:\n%s\nATENTION! All Before this point is system instructions and chat history, to generate your response consider the current user message -> = %s\nAnswer me in this language: %s",
		prePrompt,
		resumeHist,
		request.Message,
		middlewares.GetLanguageFromContext(c))

	dialogueResp, err := connector.GenerateResponse(context.Background(), fullPrompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": dialogueResp})
}

// @Summary Gera uma correção de texto usando IA
// @Description Analisa e corrige erros gramaticais no texto fornecido
// @Tags AI
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Token de autenticação"
// @Param message body object{message=string} true "Chat object"
// @Success 200 {object} map[string]string "Correção gerada com sucesso" example({"response":"I went to school yesterday"})
// @Failure 400 {object} map[string]string "Erro na requisição" example({"error":"Invalid request"})
// @Failure 500 {object} map[string]string "Erro interno do servidor" example({"error":"Internal server error"})
// @Router /ai/generate-response-correction [post]
func GenerateResponseCorrection(c *gin.Context) {
	if !enforcePlanLimits(c) {
		return
	}

	// ... (prompt loading, request binding)
	promptPath := filepath.Join("pkg/prompts", "promptCorrection.txt")
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

	// Use the builder to get a connector instance
	connector := aiConnectorBuilder() // MODIFIED LINE

	// Concatenate prompt with user message for correction
	fullPrompt := fmt.Sprintf("%s\n\nINPUT:\n%s\n\nOUTPUT:", prompt, request.Message)

	var correctionResp string
	if ollamaConnector, ok := connector.(*connectors.OllamaConnector); ok {
		systemPrompt := "You are a strict text correction tool. Return ONLY the corrected text. Do not answer the question. Do not add explanations."
		options := map[string]any{
			"temperature": 0.1,
			"top_p":       0.9,
			"num_predict": 128,
		}
		correctionResp, err = ollamaConnector.GenerateResponseWithOptions(context.Background(), fullPrompt, systemPrompt, options)
	} else {
		correctionResp, err = connector.GenerateResponse(context.Background(), fullPrompt)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	correctionResp = strings.TrimSpace(correctionResp)
	if correctionResp == "" {
		correctionResp = request.Message
	}

	c.JSON(http.StatusOK, gin.H{"response": correctionResp})
}

// @Summary Traduz um texto usando IA
// @Description Recebe um texto e retorna sua tradução para o idioma especificado
// @Tags AI
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Token de autenticação"
// @Param message body object{message=string,target_language=string} true "Texto e idioma de destino"
// @Success 200 {object} map[string]string "Tradução gerada com sucesso" example({"response":"Olá, como vai você?"})
// @Failure 400 {object} map[string]string "Requisição inválida" example({"error":"Invalid request"})
// @Failure 500 {object} map[string]string "Erro interno" example({"error":"Internal server error"})
// @Router /ai/generate-response-translation [post]
func GenerateResponseTranslate(c *gin.Context) {
	if !enforcePlanLimits(c) {
		return
	}

	promptPath := filepath.Join("pkg/prompts", "promptTranslate.txt")
	promptBytes, err := os.ReadFile(promptPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao carregar o prompt: " + err.Error()})
		return
	}
	prompt := string(promptBytes)

	var req struct {
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	// Use the builder to get a connector instance
	connector := aiConnectorBuilder() // MODIFIED LINE

	// Concatenate prompt with user message for translation
	fullPrompt := prompt + "\n\nText: " + req.Message + "\nTranslation:"

	response, err := connector.GenerateResponse(context.Background(), fullPrompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar tradução: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": response})
}

// @Summary Gera um tópico para uma conversa usando IA
// @Description Analisa o texto fornecido e gera um tópico relevante de duas palavras
// @Tags AI
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Token de autenticação"
// @Param message body object{topic=string} true "Chat object"
// @Success 200 {object} map[string]string "Tópico gerado com sucesso" example({"response":"World Travel"})
// @Failure 400 {object} map[string]string "Erro na requisição" example({"error":"Invalid request"})
// @Failure 500 {object} map[string]string "Erro interno do servidor" example({"error":"Internal server error"})
// @Router /ai/generate-response-topic [post]
func GenerateResponseTopic(c *gin.Context) {
	if !enforcePlanLimits(c) {
		return
	}

	var request struct {
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use the builder to get a connector instance
	connector := aiConnectorBuilder() // MODIFIED LINE

	strictPrompt := fmt.Sprintf(
		"You are a topic generator. Return ONLY a topic with exactly 2 words. No punctuation, no quotes, no extra text.\n\nExamples:\nInput: I want to learn English for travel.\nOutput: Travel English\nInput: We discussed ordering food in restaurants.\nOutput: Food Ordering\n\nInput: %s\nOutput:",
		request.Message,
	)

	var topicResp string
	var err error
	if ollamaConnector, ok := connector.(*connectors.OllamaConnector); ok {
		systemPrompt := "Return ONLY a two-word topic. No explanations."
		options := map[string]any{
			"temperature": 0.1,
			"top_p":       0.9,
			"num_predict": 12,
		}
		topicResp, err = ollamaConnector.GenerateResponseWithOptions(context.Background(), strictPrompt, systemPrompt, options)
	} else {
		topicResp, err = connector.GenerateResponse(context.Background(), strictPrompt)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	topicResp = strings.TrimSpace(topicResp)
	if topicResp == "" {
		topicResp = "New Topic"
	}

	c.JSON(http.StatusOK, gin.H{"response": topicResp})
}
