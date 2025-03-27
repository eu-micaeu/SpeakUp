package routes

import (

	"github.com/gin-gonic/gin"
	"speakup/handlers"

	"speakup/middlewares"

)

func MessageRoutes(router *gin.Engine) {

	messageRoutes := router.Group("/message")

	{
		// CRUD operations for messages
		messageRoutes.POST("", middlewares.AuthMiddleware(), handlers.CreateChat)
		messageRoutes.GET("/:id", middlewares.AuthMiddleware(), handlers.GetChatById)
		messageRoutes.GET("", middlewares.AuthMiddleware(), handlers.GetChats)
		messageRoutes.PUT("/:id", middlewares.AuthMiddleware(), handlers.UpdateChat)
		messageRoutes.DELETE("/:id", middlewares.AuthMiddleware(), handlers.DeleteChat)

		// Get messages by chat ID
		messageRoutes.GET("/chat/:id", middlewares.AuthMiddleware(), handlers.GetMessagesByChatId)
	}

}