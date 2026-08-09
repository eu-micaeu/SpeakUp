package routes

import (
	"speakup/pkg/handlers"

	"github.com/gin-gonic/gin"
)

func FlashcardRoutes(r *gin.Engine, handler *handlers.FlashcardHandler) {
	api := r.Group("/api/flashcards")
	{
		api.POST("/generate", handler.GenerateFlashcard)
		api.POST("", handler.CreateFlashcard)
		api.GET("", handler.GetFlashcards)
		api.POST("/:id/review", handler.ReviewFlashcard)
		api.DELETE("/:id", handler.DeleteFlashcard)
	}
}
