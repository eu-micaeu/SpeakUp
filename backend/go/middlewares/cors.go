package middlewares

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config := cors.DefaultConfig()

	config.AllowOrigins = []string{"http://localhost:3000"}                   // Permitir apenas a origem do frontend
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"} // Métodos permitidos
	config.AllowHeaders = []string{"Authorization", "Content-Type"}           // Cabeçalhos permitidos

	// Permitir credenciais (opcional se necessário)
	config.AllowCredentials = true
	config.AllowOrigins = []string{"http://localhost:3000"}                   // Permitir apenas a origem do frontend
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"} // Métodos permitidos
	config.AllowHeaders = []string{"Authorization", "Content-Type"}           // Cabeçalhos permitidos

	// Permitir credenciais (opcional se necessário)
	config.AllowCredentials = true

	return cors.New(config)
	return cors.New(config)
}
