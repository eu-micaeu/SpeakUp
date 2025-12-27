package repositories

import (
	"context"
	"speakup/pkg/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
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
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var chats []models.Chat
	err = cursor.All(ctx, &chats)
	return chats, err
}

func (r *MongoChatRepository) FindAllByUserID(ctx context.Context, userID string) ([]models.Chat, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	var chats []models.Chat
	err = cursor.All(ctx, &chats)
	return chats, err
}

func (r *MongoChatRepository) Update(ctx context.Context, id string, chat models.Chat) error {
	// Note: We might want to use $set here instead of replacing the whole document to avoid overwriting fields unintentionally,
	// but sticking to original logic for now which seemed to replace.
	// However, usually UpdateOne takes an update document. The original code used UpdateOne with 'chat' as the update?
	// Let's double check the original code logic.
	// "collection.UpdateOne(c, map[string]string{"_id": id}, chat)" in original code might be wrong if 'chat' is a struct?
	// UpdateOne usually expects a '$set' operator or similar if passing a struct directly it might fail or replace.
	// Actually, ReplaceOne is for replacing. UpdateOne expects atomic operators.
	// If the original code was working with UpdateOne(..., chat), maybe the driver handles it if 'chat' has bson tags?
	// But usually UpdateOne expects {$set: ...}.
	// Let's assume the intention is to update fields. I will use $set.
	update := bson.M{"$set": chat}
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (r *MongoChatRepository) Delete(ctx context.Context, id string) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
