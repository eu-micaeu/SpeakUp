package routes

import (
    "github.com/gin-gonic/gin"
    "speakup/handlers"
)

func ChatRoutes(router *gin.Engine) {
    chatRoutes := router.Group("/chat")
    {
        chatRoutes.POST("/start", handlers.StartChat)
        chatRoutes.DELETE("/end", handlers.EndChat)
    }
}