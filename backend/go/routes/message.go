package routes

import (

	"github.com/gin-gonic/gin"
	"speakup/handlers"

	"speakup/middlewares"

)

func MessageRoutes(router *gin.Engine) {

	messageRoutes := router.Group("api/message")

	{
		// CRUD operations for messages
		messageRoutes.POST("", middlewares.AuthMiddleware(), handlers.CreateMessage)
		messageRoutes.GET("/:id", middlewares.AuthMiddleware(), handlers.GetMessageById)
		messageRoutes.GET("", middlewares.AuthMiddleware(), handlers.GetMessages)
		messageRoutes.PUT("/:id", middlewares.AuthMiddleware(), handlers.UpdateMessage)
		messageRoutes.DELETE("/:id", middlewares.AuthMiddleware(), handlers.DeleteMessage)

		// Get messages by chat ID
		messageRoutes.GET("/chat/:id", middlewares.AuthMiddleware(), handlers.GetMessagesByChatId)
	}

}