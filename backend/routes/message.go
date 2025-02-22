package routes

import (

	"github.com/gin-gonic/gin"
	"speakup/handlers"

)

func MessageRoutes(router *gin.Engine) {

	messageRoutes := router.Group("/message")

	{

		messageRoutes.POST("/send", handlers.Send)
	}

}