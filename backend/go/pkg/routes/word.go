package routes

import (
	"github.com/gin-gonic/gin"
	"speakup/pkg/handlers"
	"speakup/pkg/middlewares"
)

func WordRoutes(router *gin.Engine, wordHandler *handlers.WordHandler) {
	wordRoutes := router.Group("api/word")
	{
		wordRoutes.GET("/user", middlewares.AuthMiddleware(), wordHandler.ListUserWords)
	}
}