package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"speakup/config"
	"speakup/middlewares"
	"speakup/routes"
)

// @title           SpeakUp API
// @version         1.0
// @description     API para o projeto SpeakUp
// @host            localhost:8080
// @BasePath        /

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

	// Load Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run()

}
