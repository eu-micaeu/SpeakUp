package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"speakup/config"
	"speakup/connectors"
	"speakup/middlewares"
	"speakup/models"

	"github.com/gin-gonic/gin"
)

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
	promptPath := filepath.Join("prompts", "promptDialog.txt")
	promptBytes, err := os.ReadFile(promptPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load prompt: " + err.Error()})
		return
	}
	basePrompt := string(promptBytes)

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

	// Gerar resposta diretamente com histórico
	fullPrompt := fmt.Sprintf(
		"%s\n\nChat:\n%s\n\nUsuário: %s\n\nResponda nesse idioma: %s",
		basePrompt,
		chatHistory.String(),
		request.Message,
		middlewares.GetLanguageFromContext(c),
	)

	connector := connectors.NewGeminiConnector()
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
	promptPath := filepath.Join("prompts", "promptCorrection.txt")
	promptBytes, err := os.ReadFile(promptPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load prompt: " + err.Error()})
		return
	}
	basePrompt := string(promptBytes)

	var request struct {
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	language := middlewares.GetLanguageFromContext(c)

	finalPrompt := fmt.Sprintf(
		"%s\n\nAnswer me in this language: %s\nInput: %s",
		basePrompt, language, request.Message,
	)

	connector := connectors.NewGeminiConnector()
	correctionResp, err := connector.GenerateResponse(context.Background(), finalPrompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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
	promptPath := filepath.Join("prompts", "promptTranslate.txt")

	promptBytes, err := os.ReadFile(promptPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao carregar o prompt: " + err.Error()})
		return
	}
	basePrompt := string(promptBytes)

	var req struct {
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	connector := connectors.NewGeminiConnector()

	// Garante que o modelo entenda claramente o input como parte do prompt
	finalPrompt := fmt.Sprintf("%s\n\nTexto para traduzir:\n%s", basePrompt, req.Message)

	response, err := connector.GenerateResponse(context.Background(), finalPrompt)
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
	var request struct {
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	connector := connectors.NewGeminiConnector()

	// Generate a response for the topic
	topicResp, err := connector.GenerateResponse(context.Background(), "Please generate a topic for the following text, return only words ,generate with 2 words only: "+request.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": topicResp})
}
