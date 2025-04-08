package routes

import (

	"github.com/gin-gonic/gin"
	"speakup/handlers"

)

// UserRoutes sets up the routes for the user
func UserRoutes(router *gin.Engine) {
	userRoutes := router.Group("api/user")
	{
		// RESTful routes
		userRoutes.POST("/", handlers.CreateUser)
		userRoutes.GET("/:id", handlers.GetUsers)
		userRoutes.PUT("/:id", handlers.UpdateUser)
		userRoutes.DELETE("/:id", handlers.DeleteUser)

		// Custom routes
		userRoutes.POST("/login", handlers.Login)
	}
}