package routes

import (

	"github.com/gin-gonic/gin"
	"speakup/pkg/handlers"

	"speakup/pkg/middlewares"

)

func MessageRoutes(router *gin.Engine, messageHandler *handlers.MessageHandler) {

	messageRoutes := router.Group("api/message")

	{
		// CRUD operations for messages
		messageRoutes.POST("", middlewares.AuthMiddleware(), messageHandler.CreateMessage)
		messageRoutes.GET("/:id", middlewares.AuthMiddleware(), messageHandler.GetMessageById)
		messageRoutes.GET("", middlewares.AuthMiddleware(), messageHandler.GetMessages)
		messageRoutes.PUT("/:id", middlewares.AuthMiddleware(), messageHandler.UpdateMessage)
		messageRoutes.DELETE("/:id", middlewares.AuthMiddleware(), messageHandler.DeleteMessage)

		// Get messages by chat ID
		messageRoutes.GET("/chat/:id", middlewares.AuthMiddleware(), messageHandler.GetMessagesByChatId)
	}

}