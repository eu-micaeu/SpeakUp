package routes

import (

	"github.com/gin-gonic/gin"
	"speakup/handlers"

)

func UserRoutes(router *gin.Engine) {

	userRoutes := router.Group("/user")

	{

		userRoutes.POST("/login", handlers.Login)

		userRoutes.POST("/create", handlers.CreateUser)

		userRoutes.DELETE("/delete", handlers.DeleteUser)

	}

}