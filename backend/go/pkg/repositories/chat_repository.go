package repositories

import (
	"context"
	"speakup/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatRepository interface {
	Create(ctx context.Context, chat models.Chat) error
	FindByID(ctx context.Context, id string) (models.Chat, error)
	FindAll(ctx context.Context) ([]models.Chat, error)
	FindAllByUserID(ctx context.Context, userID string) ([]models.Chat, error)
	Update(ctx context.Context, id string, chat models.Chat) error
	Delete(ctx context.Context, id string) error
}

type MongoChatRepository struct {
	Collection *mongo.Collection
}

func NewMongoChatRepository(db *mongo.Database) *MongoChatRepository {
	return &MongoChatRepository{
		Collection: db.Collection("chats"),
	}
}

func (r *MongoChatRepository) Create(ctx context.Context, chat models.Chat) error {
	_, err := r.Collection.InsertOne(ctx, chat)
	return err
}

func (r *MongoChatRepository) FindByID(ctx context.Context, id string) (models.Chat, error) {
	var chat models.Chat
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&chat)
	return chat, err
}

func (r *MongoChatRepository) FindAll(ctx context.Context) ([]models.Chat, error) {
	opts := options.Find().SetSort(bson.D{{Key: "start_time", Value: -1}})
	cursor, err := r.Collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var chats []models.Chat
	err = cursor.All(ctx, &chats)
	if chats == nil {
		chats = []models.Chat{}
	}
	return chats, err
}

func (r *MongoChatRepository) FindAllByUserID(ctx context.Context, userID string) ([]models.Chat, error) {
	opts := options.Find().SetSort(bson.D{{Key: "start_time", Value: -1}})
	cursor, err := r.Collection.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var chats []models.Chat
	err = cursor.All(ctx, &chats)
	if chats == nil {
		chats = []models.Chat{}
	}
	return chats, err
}

func (r *MongoChatRepository) Update(ctx context.Context, id string, chat models.Chat) error {
	update := bson.M{"$set": chat}
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *MongoChatRepository) Delete(ctx context.Context, id string) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

