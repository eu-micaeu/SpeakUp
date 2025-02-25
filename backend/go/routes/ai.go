package routes

import (
    "github.com/gin-gonic/gin"
    "speakup/handlers"
    "speakup/middlewares"
)

func AIRoutes(router *gin.Engine) {
    aiRoutes := router.Group("/ai")
    {
        aiRoutes.POST("/generate-response", handlers.GenerateResponse)
        aiRoutes.POST("/detect-spelling-errors", middlewares.AuthMiddleware(), handlers.DetectSpellingErrors)
    }
}