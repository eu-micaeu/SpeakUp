package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func ConnectMongoDB() {
	// Pegar a string de conexão das variáveis de ambiente
	mongoURI := os.Getenv("MONGO_CONNECTION_STRING")
	if mongoURI == "" {
		log.Fatal("MONGO_CONNECTION_STRING não está definida")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Falha ao conectar ao MongoDB: %v", err)
	}

	// Ping no banco para verificar a conexão
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Falha ao pingar o MongoDB: %v", err)
	}

	fmt.Println("Conectado ao MongoDB!")
}

func GetMongoClient() *mongo.Client {
	return client
}