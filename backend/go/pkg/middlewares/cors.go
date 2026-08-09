package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()

	config.AllowOriginFunc = func(origin string) bool { return true }
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	config.AllowHeaders = []string{
		"Authorization",
		"Content-Type",
		"X-AI-Provider",
		"X-User-ID",
		"X-User-Id",
		"x-user-id",
		"Origin",
		"Accept",
	}

	return cors.New(config)
}
