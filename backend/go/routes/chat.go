package routes

import (
    "github.com/gin-gonic/gin"
    "speakup/handlers"
)

func ChatRoutes(router *gin.Engine) {
    chatRoutes := router.Group("/chat")
    {
        // CRUD operations for chat
        chatRoutes.POST("/", handlers.CreateChat)
        chatRoutes.GET("/:id", handlers.GetChatById)
        chatRoutes.GET("/", handlers.GetChats)
        chatRoutes.PUT("/:id", handlers.UpdateChat)
        chatRoutes.DELETE("/:id", handlers.DeleteChat)

        // special routes
        chatRoutes.GET("/user/:id", handlers.GetChatsByUserId)
    }
}