package routes

import (
    "github.com/gin-gonic/gin"

    "speakup/pkg/handlers"
    "speakup/pkg/middlewares"
)

func ChatRoutes(router *gin.Engine, chatHandler *handlers.ChatHandler) {
    chatRoutes := router.Group("api/chat")
    {
        // CRUD operations for chat
        chatRoutes.POST("", middlewares.AuthMiddleware(), chatHandler.CreateChat)
        chatRoutes.GET("/:id", middlewares.AuthMiddleware(), chatHandler.GetChatById)
        chatRoutes.GET("/", middlewares.AuthMiddleware(), chatHandler.GetChats)
        chatRoutes.PUT("/:id", middlewares.AuthMiddleware(), chatHandler.UpdateChat)
        chatRoutes.DELETE("/:id", middlewares.AuthMiddleware(), chatHandler.DeleteChat)

        // special routes
        chatRoutes.GET("/user", middlewares.AuthMiddleware(), chatHandler.GetChatsByUserId)
    }
}