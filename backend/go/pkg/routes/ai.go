package routes

import (
	"speakup/pkg/handlers"
	"speakup/pkg/middlewares"

	"github.com/gin-gonic/gin"
)

func AIRoutes(router *gin.Engine) {

	aiRoutes := router.Group("api/ai")

	{

		aiRoutes.POST("/generate-response-dialog", middlewares.AuthMiddleware(), handlers.GenerateResponseDialog)
		aiRoutes.POST("/generate-response-correction", middlewares.AuthMiddleware(), handlers.GenerateResponseCorrection)
		aiRoutes.POST("/generate-response-translation", middlewares.AuthMiddleware(), handlers.GenerateResponseTranslate)
		aiRoutes.POST("/generate-response-topic", middlewares.AuthMiddleware(), handlers.GenerateResponseTopic)
		aiRoutes.POST("/transcribe-audio", middlewares.AuthMiddleware(), handlers.TranscribeAudio)
		aiRoutes.GET("/usage", middlewares.AuthMiddleware(), handlers.GetAIUsage)

	}

}
