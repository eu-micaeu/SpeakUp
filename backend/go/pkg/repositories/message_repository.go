package repositories

import (
	"context"
	"speakup/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageRepository interface {
	Create(ctx context.Context, message models.Message) error
	FindByID(ctx context.Context, id string) (models.Message, error)
	FindAll(ctx context.Context) ([]models.Message, error)
	FindAllByChatID(ctx context.Context, chatID string) ([]models.Message, error)
	FindRecentByChatID(ctx context.Context, chatID string, limit int) ([]models.Message, error)
	Update(ctx context.Context, id string, updateDoc bson.M) error
	Delete(ctx context.Context, id string) error
	DeleteByChatID(ctx context.Context, chatID string) error
}

type MongoMessageRepository struct {
	Collection *mongo.Collection
}

func NewMongoMessageRepository(db *mongo.Database) *MongoMessageRepository {
	return &MongoMessageRepository{
		Collection: db.Collection("messages"),
	}
}

func messageIDFilter(id string) bson.M {
	or := []bson.M{
		{"_id": id},
		{"id": id},
	}
	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		or = append(or, bson.M{"_id": oid})
	}
	return bson.M{"$or": or}
}

func normalizeMessages(messages []models.Message) []models.Message {
	if messages == nil {
		return []models.Message{}
	}
	for i := range messages {
		messages[i].Normalize()
	}
	return messages
}

func (r *MongoMessageRepository) Create(ctx context.Context, message models.Message) error {
	_, err := r.Collection.InsertOne(ctx, message)
	return err
}

func (r *MongoMessageRepository) FindByID(ctx context.Context, id string) (models.Message, error) {
	var message models.Message
	err := r.Collection.FindOne(ctx, messageIDFilter(id)).Decode(&message)
	if err == nil {
		message.Normalize()
	}
	return message, err
}

func (r *MongoMessageRepository) FindAll(ctx context.Context) ([]models.Message, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}})
	cursor, err := r.Collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []models.Message
	err = cursor.All(ctx, &messages)
	return normalizeMessages(messages), err
}

func (r *MongoMessageRepository) FindAllByChatID(ctx context.Context, chatID string) ([]models.Message, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}})
	cursor, err := r.Collection.Find(ctx, bson.M{"chat_id": chatID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []models.Message
	err = cursor.All(ctx, &messages)
	return normalizeMessages(messages), err
}

func (r *MongoMessageRepository) FindRecentByChatID(ctx context.Context, chatID string, limit int) ([]models.Message, error) {
	if limit <= 0 {
		limit = 10
	}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.Collection.Find(ctx, bson.M{"chat_id": chatID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []models.Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	// Reverse back into chronological order (oldest -> newest)
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return normalizeMessages(messages), nil
}

func (r *MongoMessageRepository) Update(ctx context.Context, id string, updateDoc bson.M) error {
	_, err := r.Collection.UpdateOne(ctx, messageIDFilter(id), bson.M{"$set": updateDoc})
	return err
}

func (r *MongoMessageRepository) Delete(ctx context.Context, id string) error {
	_, err := r.Collection.DeleteOne(ctx, messageIDFilter(id))
	return err
}

func (r *MongoMessageRepository) DeleteByChatID(ctx context.Context, chatID string) error {
	_, err := r.Collection.DeleteMany(ctx, bson.M{"chat_id": chatID})
	return err
}

