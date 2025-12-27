package routes

import (

	"github.com/gin-gonic/gin"
	"speakup/pkg/handlers"

)

// UserRoutes sets up the routes for the user
func UserRoutes(router *gin.Engine, userHandler *handlers.UserHandler) {
	userRoutes := router.Group("api/user")
	{
		// RESTful routes
		userRoutes.POST("/", userHandler.CreateUser)
		userRoutes.GET("/:id", userHandler.GetUsers)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)

		// Custom routes
		userRoutes.POST("/login", userHandler.Login)
	}
}