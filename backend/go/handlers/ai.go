package handlers

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"speakup/connectors"
	"speakup/middlewares"

	"github.com/gin-gonic/gin"
)

// @Summary Gera uma resposta de diálogo usando IA
// @Description Gera uma resposta de diálogo contextual baseada na mensagem fornecida
// @Tags AI
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Token de autenticação"
// @Param message body object{message=string} true "Chat object"
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
	dialogueResp, err := connector.GenerateResponse(context.Background(), prompt+request.Message+"Answer me in this language: "+middlewares.GetLanguageFromContext(c))
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
	prompt := string(promptBytes)

	var request struct {
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	connector := connectors.NewGeminiConnector()

	// Generate a correction for the dialogue
	correctionResp, err := connector.GenerateResponse(context.Background(), "Answer me in this language: " + middlewares.GetLanguageFromContext(c) + prompt + request.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": correctionResp})
}

// @Summary Gera uma tradução de texto usando IA
// @Description Traduz o texto fornecido para o idioma especificado
// @Tags AI
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Token de autenticação"
// @Param message body object{message=string,target_language=string} true "Chat object"
// @Success 200 {object} map[string]string "Tradução gerada com sucesso" example({"response":"Olá, como vai você?"})
// @Failure 400 {object} map[string]string "Erro na requisição" example({"error":"Invalid request"})
// @Failure 500 {object} map[string]string "Erro interno do servidor" example({"error":"Internal server error"})
// @Router /ai/generate-response-translation [post]
func GenerateResponseTranslate(c *gin.Context) {

	promptPath := filepath.Join("prompts", "promptTranslate.txt")
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

	// Generate a translation for the dialogue
	translateResp, err := connector.GenerateResponse(context.Background(), prompt+request.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": translateResp})
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
