package handlers

import (
	"net/http"

	"speakup/pkg/middlewares"
	"speakup/pkg/repositories"
	"speakup/pkg/utils"

	"github.com/gin-gonic/gin"
)

type WordHandler struct {
	Repo repositories.WordRepository
}

func NewWordHandler(repo repositories.WordRepository) *WordHandler {
	return &WordHandler{Repo: repo}
}

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
func (h *WordHandler) ListUserWords(c *gin.Context) {
	// Obter ID do usuário do contexto
	userID := middlewares.GetUserIDFromContext(c)
	if userID == "" {
		utils.RespondWithError(c, http.StatusUnauthorized, "Usuário não autenticado")
		return
	}

	words, err := h.Repo.FindAllByUserID(c, userID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Falha ao buscar palavras")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, words)
}