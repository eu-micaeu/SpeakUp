package repositories

import (
	"context"
	"speakup/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageRepository interface {
	Create(ctx context.Context, message models.Message) error
	FindByID(ctx context.Context, id string) (models.Message, error)
	FindAll(ctx context.Context) ([]models.Message, error)
	FindAllByChatID(ctx context.Context, chatID string) ([]models.Message, error)
	Update(ctx context.Context, id string, updateDoc bson.M) error
	Delete(ctx context.Context, id string) error
}

type MongoMessageRepository struct {
	Collection *mongo.Collection
}

func NewMongoMessageRepository(db *mongo.Database) *MongoMessageRepository {
	return &MongoMessageRepository{
		Collection: db.Collection("messages"),
	}
}

func (r *MongoMessageRepository) Create(ctx context.Context, message models.Message) error {
	_, err := r.Collection.InsertOne(ctx, message)
	return err
}

func (r *MongoMessageRepository) FindByID(ctx context.Context, id string) (models.Message, error) {
	var message models.Message
	err := r.Collection.FindOne(ctx, bson.M{"id": id}).Decode(&message)
	return message, err
}

func (r *MongoMessageRepository) FindAll(ctx context.Context) ([]models.Message, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var messages []models.Message
	err = cursor.All(ctx, &messages)
	return messages, err
}

func (r *MongoMessageRepository) FindAllByChatID(ctx context.Context, chatID string) ([]models.Message, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{"chat_id": chatID})
	if err != nil {
		return nil, err
	}
	var messages []models.Message
	err = cursor.All(ctx, &messages)
	return messages, err
}

func (r *MongoMessageRepository) Update(ctx context.Context, id string, updateDoc bson.M) error {
	_, err := r.Collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": updateDoc})
	return err
}

func (r *MongoMessageRepository) Delete(ctx context.Context, id string) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}
