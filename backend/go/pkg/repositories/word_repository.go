package repositories

import (
	"context"
	"speakup/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WordRepository interface {
	FindAllByUserID(ctx context.Context, userID string) ([]models.Word, error)
}

type MongoWordRepository struct {
	Collection *mongo.Collection
}

func NewMongoWordRepository(db *mongo.Database) *MongoWordRepository {
	return &MongoWordRepository{
		Collection: db.Collection("words"),
	}
}

func (r *MongoWordRepository) FindAllByUserID(ctx context.Context, userID string) ([]models.Word, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{
		"user_id": userID,
	}, options.Find().SetSort(bson.M{"created_at": -1}))
	
	if err != nil {
		return nil, err
	}
	
	var words []models.Word
	err = cursor.All(ctx, &words)
	return words, err
}
