package handlers

import (
	"net/http"

	"speakup/config"
	"speakup/middlewares"
	"speakup/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Lista todas as palavras de um usuário
// @Description Retorna todas as palavras geradas para um usuário específico
// @Tags Word
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Token de autenticação"
// @Success 200 {array} models.Word "Lista de palavras"
// @Failure 401 {object} map[string]string "Erro de autenticação" example({"error":"Usuário não autenticado"})
// @Failure 500 {object} map[string]string "Erro interno do servidor" example({"error":"Falha ao buscar palavras"})
// @Router /word/user [get]
func ListUserWords(c *gin.Context) {
	// Obter ID do usuário do contexto
	userID := middlewares.GetUserIDFromContext(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return
	}

	// Buscar palavras do usuário
	db := config.GetMongoClient()
	collection := db.Database("speakup").Collection("words")

	cursor, err := collection.Find(c, bson.M{
		"user_id": userID,
	}, options.Find().SetSort(bson.M{"created_at": -1}))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao buscar palavras"})
		return
	}

	var words []models.Word
	if err = cursor.All(c, &words); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao processar palavras"})
		return
	}

	c.JSON(http.StatusOK, words)
}