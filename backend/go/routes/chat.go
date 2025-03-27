package routes

import (
    "github.com/gin-gonic/gin"

    "speakup/handlers"
    "speakup/middlewares"
)

func ChatRoutes(router *gin.Engine) {
    chatRoutes := router.Group("/chat")
    {
        // CRUD operations for chat
        chatRoutes.POST("", middlewares.AuthMiddleware(), handlers.CreateChat)
        chatRoutes.GET("/:id", middlewares.AuthMiddleware(), handlers.GetChatById)
        chatRoutes.GET("/", middlewares.AuthMiddleware(), handlers.GetChats)
        chatRoutes.PUT("/:id", middlewares.AuthMiddleware(), handlers.UpdateChat)
        chatRoutes.DELETE("/:id", middlewares.AuthMiddleware(), handlers.DeleteChat)

        // special routes
        chatRoutes.GET("/user", middlewares.AuthMiddleware(), handlers.GetChatsByUserId)
    }
}