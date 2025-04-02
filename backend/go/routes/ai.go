package routes

import (
    "github.com/gin-gonic/gin"
    "speakup/handlers"
    "speakup/middlewares"
)

func AIRoutes(router *gin.Engine) {

    aiRoutes := router.Group("/ai")

    {

        aiRoutes.POST("/generate-response-dialog", middlewares.AuthMiddleware(), handlers.GenerateResponseDialog)
        aiRoutes.POST("/generate-response-correction", middlewares.AuthMiddleware(), handlers.GenerateResponseCorrection)
        aiRoutes.POST("/generate-response-translation", middlewares.AuthMiddleware(), handlers.GenerateResponseTranslate)
    
    }

}