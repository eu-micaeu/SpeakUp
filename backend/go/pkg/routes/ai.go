package routes

import (
	aiHandlers "speakup/pkg/adapters/ai"
	"speakup/pkg/handlers"
	"speakup/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func AIRoutes(router *gin.Engine) {

	aiRoutes := router.Group("api/ai")

	{

		aiRoutes.POST("/generate-response-dialog", middlewares.AuthMiddleware(), aiHandlers.GenerateResponseDialog)
		aiRoutes.POST("/generate-response-correction", middlewares.AuthMiddleware(), aiHandlers.GenerateResponseCorrection)
		aiRoutes.POST("/generate-response-translation", middlewares.AuthMiddleware(), aiHandlers.GenerateResponseTranslate)
		aiRoutes.POST("/generate-response-topic", middlewares.AuthMiddleware(), aiHandlers.GenerateResponseTopic)
		aiRoutes.POST("/transcribe-audio", middlewares.AuthMiddleware(), handlers.TranscribeAudio)
		aiRoutes.GET("/usage", middlewares.AuthMiddleware(), aiHandlers.GetAIUsage)

	}

}
