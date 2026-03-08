package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"speakup/pkg/config"
	"speakup/pkg/handlers"
	"speakup/pkg/middlewares"
	"speakup/pkg/repositories"
	"speakup/pkg/routes"
)

// @title           SpeakUp API
// @version         1.0
// @description     API para o projeto SpeakUp
// @host            localhost:8080
// @BasePath        /

func main() {

	// Carregar variáveis de ambiente do arquivo .env PRIMEIRO
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	r := gin.Default()

	// CORS
	r.Use(middlewares.CorsMiddleware())

	// Connect to MongoDB
	config.ConnectMongoDB()

	// Initialize Dependencies
	dbClient := config.GetMongoClient()
	dbName := config.GetDBName()
	db := dbClient.Database(dbName)

	// Repositories
	chatRepo := repositories.NewMongoChatRepository(db)
	userRepo := repositories.NewMongoUserRepository(db)
	messageRepo := repositories.NewMongoMessageRepository(db)

	// Handlers
	chatHandler := handlers.NewChatHandler(chatRepo)
	userHandler := handlers.NewUserHandler(userRepo)
	messageHandler := handlers.NewMessageHandler(messageRepo)
	billingHandler := handlers.NewBillingHandler(userRepo)

	// Load routes
	routes.UserRoutes(r, userHandler)
	routes.ChatRoutes(r, chatHandler)
	routes.MessageRoutes(r, messageHandler)
	routes.AIRoutes(r)
	routes.BillingRoutes(r, billingHandler)

	// Load Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8082")

}
