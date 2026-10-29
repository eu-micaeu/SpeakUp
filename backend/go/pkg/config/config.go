package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

}

func EnsureIndexes(db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Index messages by (chat_id, created_at)
	_, err := db.Collection("messages").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "chat_id", Value: 1},
			{Key: "created_at", Value: 1},
		},
	})
	if err != nil {
		log.Printf("Aviso ao criar índice em messages: %v", err)
	}

	// Index chats by (user_id, start_time)
	_, err = db.Collection("chats").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "start_time", Value: -1},
		},
	})
	if err != nil {
		log.Printf("Aviso ao criar índice em chats: %v", err)
	}

	// Unique index on users(email)
	_, err = db.Collection("users").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		log.Printf("Aviso ao criar índice em users: %v", err)
	}

	// Index flashcards by (user_id, next_review)
	_, err = db.Collection("flashcards").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "next_review", Value: 1},
		},
	})
	if err != nil {
		log.Printf("Aviso ao criar índice em flashcards: %v", err)
	}
}

func GetMongoClient() *mongo.Client {
	return client
}

func SetMongoClient(c *mongo.Client) {
	client = c
}

func GetDBName() string {
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		return "speakup"
	}
	return dbName
}
