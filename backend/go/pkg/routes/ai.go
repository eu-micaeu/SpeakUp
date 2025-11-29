package routes

import (
    "github.com/gin-gonic/gin"
    aiHandlers "speakup/pkg/adapters/ai"
    "speakup/pkg/handlers"
    "speakup/pkg/middlewares"
)

func AIRoutes(router *gin.Engine) {

    aiRoutes := router.Group("api/ai")

    {

        aiRoutes.POST("/generate-response-dialog", middlewares.AuthMiddleware(), aiHandlers.GenerateResponseDialog)
        aiRoutes.POST("/generate-response-correction", middlewares.AuthMiddleware(), aiHandlers.GenerateResponseCorrection)
        aiRoutes.POST("/generate-response-translation", middlewares.AuthMiddleware(), aiHandlers.GenerateResponseTranslate)
        aiRoutes.POST("/generate-response-topic", middlewares.AuthMiddleware(), aiHandlers.GenerateResponseTopic)
        aiRoutes.POST("/generate-random-word", middlewares.AuthMiddleware(), aiHandlers.GenerateRandomWord)
        aiRoutes.POST("/transcribe-audio", middlewares.AuthMiddleware(), handlers.TranscribeAudio)
    
    }

}