package repositories

import (
	"context"
	"regexp"
	"strings"
	"time"

	"speakup/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FlashcardRepository interface {
	Create(ctx context.Context, flashcard models.Flashcard) (models.Flashcard, error)
	FindByID(ctx context.Context, id string) (models.Flashcard, error)
	FindByUser(ctx context.Context, userID string, dueOnly bool) ([]models.Flashcard, error)
	FindByUserAndFront(ctx context.Context, userID string, front string) (models.Flashcard, error)
	Update(ctx context.Context, id string, updateDoc bson.M) error
	Delete(ctx context.Context, id string) error
}

type MongoFlashcardRepository struct {
	Collection *mongo.Collection
}

func NewMongoFlashcardRepository(db *mongo.Database) *MongoFlashcardRepository {
	return &MongoFlashcardRepository{
		Collection: db.Collection("flashcards"),
	}
}

func (r *MongoFlashcardRepository) Create(ctx context.Context, flashcard models.Flashcard) (models.Flashcard, error) {
	if flashcard.ID == "" {
		flashcard.ID = primitive.NewObjectID().Hex()
	}
	if flashcard.EaseFactor == 0 {
		flashcard.EaseFactor = 2.5
	}
	if flashcard.CreatedAt.IsZero() {
		flashcard.CreatedAt = time.Now()
	}
	if flashcard.NextReview.IsZero() {
		flashcard.NextReview = time.Now()
	}

	_, err := r.Collection.InsertOne(ctx, flashcard)
	return flashcard, err
}

func (r *MongoFlashcardRepository) FindByID(ctx context.Context, id string) (models.Flashcard, error) {
	var fc models.Flashcard
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&fc)
	return fc, err
}

func (r *MongoFlashcardRepository) FindByUser(ctx context.Context, userID string, dueOnly bool) ([]models.Flashcard, error) {
	filter := bson.M{}
	if userID != "" {
		filter["user_id"] = userID
	}
	if dueOnly {
		filter["next_review"] = bson.M{"$lte": time.Now()}
	}

	opts := options.Find().SetSort(bson.D{{Key: "next_review", Value: 1}})
	cursor, err := r.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var flashcards []models.Flashcard
	if err := cursor.All(ctx, &flashcards); err != nil {
		return nil, err
	}
	if flashcards == nil {
		flashcards = []models.Flashcard{}
	}
	return flashcards, nil
}

func (r *MongoFlashcardRepository) FindByUserAndFront(ctx context.Context, userID string, front string) (models.Flashcard, error) {
	var fc models.Flashcard
	frontTrimmed := strings.TrimSpace(front)
	filter := bson.M{
		"user_id": userID,
		"front":   primitive.Regex{Pattern: "^" + regexp.QuoteMeta(frontTrimmed) + "$", Options: "i"},
	}
	err := r.Collection.FindOne(ctx, filter).Decode(&fc)
	return fc, err
}

func (r *MongoFlashcardRepository) Update(ctx context.Context, id string, updateDoc bson.M) error {
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updateDoc})
	return err
}

func (r *MongoFlashcardRepository) Delete(ctx context.Context, id string) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
