package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"speakup/config"
	"speakup/routes"
	"speakup/middlewares"
)

func main() {

	r := gin.Default()

	// CORS
	r.Use(middlewares.CorsMiddleware())

	// Carregar variáveis de ambiente do arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	// Connect to MongoDB
	config.ConnectMongoDB()

	// Load routes
	routes.UserRoutes(r)
	routes.ChatRoutes(r)
	routes.MessageRoutes(r)
	routes.AIRoutes(r)

	r.Run()

}
