package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"speakup/config"
	"speakup/connectors"
	"speakup/middlewares"
	"speakup/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	prePrompt := string(promptBytes)

	var request struct {
		Message string `json:"message"`
		ChatID  string `json:"chat_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Buscar histórico de mensagens do chat
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

	// Formatar histórico de mensagens
	var chatHistory strings.Builder
	for _, msg := range messages {
		chatHistory.WriteString(fmt.Sprintf("%s: %s\n", msg.Sender, msg.Content))
	}
	connector := connectors.NewGeminiConnector()

	resumeHist, err := connector.GenerateResponse(context.Background(), "Format the following chat history to only show user reponse and AI response: "+chatHistory.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Gerar resposta com histórico
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
	prompt := string(promptBytes)

	var req struct {
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	connector := connectors.NewGeminiConnector()

	response, err := connector.GenerateResponse(context.Background(), prompt+req.Message)
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

// @Summary Gera uma palavra aleatória usando IA
// @Description Gera uma palavra aleatória baseada no nível e idioma do usuário
// @Tags AI
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Token de autenticação"
// @Success 200 {object} models.Word "Palavra gerada com sucesso"
// @Failure 400 {object} map[string]string "Erro na requisição" example({"error":"Invalid request"})
// @Failure 500 {object} map[string]string "Erro interno do servidor" example({"error":"Internal server error"})
// @Router /ai/generate-random-word [post]
func GenerateRandomWord(c *gin.Context) {
	// Carregar o prompt
	promptPath := filepath.Join("prompts", "promptRandomWord.txt")
	promptBytes, err := os.ReadFile(promptPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao carregar o prompt: " + err.Error()})
		return
	}
	prompt := string(promptBytes)

	// Obter ID do usuário do contexto
	userID := middlewares.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	// Buscar dados completos do usuário
	db := config.GetMongoClient()
	collection := db.Database("speakup").Collection("users")
	
	var user models.User
	err = collection.FindOne(c, bson.M{"id": userID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao buscar dados do usuário"})
		return
	}

	// Preparar o prompt com os dados do usuário
	fullPrompt := fmt.Sprintf(prompt, user.Level, user.Language)

	// Gerar palavra usando IA
	connector := connectors.NewGeminiConnector()
	response, err := connector.GenerateResponse(context.Background(), fullPrompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar palavra: " + err.Error()})
		return
	}

	// Limpar a resposta - remover caracteres indesejados
	response = strings.TrimSpace(response)
	if strings.HasPrefix(response, "```json") {
		response = strings.TrimPrefix(response, "```json")
	}
	if strings.HasSuffix(response, "```") {
		response = strings.TrimSuffix(response, "```")
	}
	response = strings.TrimSpace(response)

	// Converter resposta para struct
	var word models.Word
	err = json.Unmarshal([]byte(response), &word)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar resposta: " + err.Error()})
		return
	}

	// Preencher campos adicionais
	word.ID = primitive.NewObjectID().Hex()
	word.UserID = userID
	word.Language = user.Language
	word.Level = user.Level
	word.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	// Salvar no MongoDB
	db = config.GetMongoClient()
	collection = db.Database("speakup").Collection("words")
	_, err = collection.InsertOne(c, word)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar palavra: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, word)
}
