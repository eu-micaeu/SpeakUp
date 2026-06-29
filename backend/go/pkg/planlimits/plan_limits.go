package planlimits

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"speakup/pkg/config"
	// "speakup/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultFreeDailyLimit int64 = 10

func GetFreeDailyLimit() int64 {
	if raw := strings.TrimSpace(os.Getenv("FREE_DAILY_AI_LIMIT")); raw != "" {
		if value, err := parseInt64(raw); err == nil && value > 0 {
			return value
		}
	}
	return defaultFreeDailyLimit
}

func parseInt64(value string) (int64, error) {
	var result int64
	_, err := fmt.Sscanf(value, "%d", &result)
	return result, err
}

func IsProUser(ctx context.Context, userID string) (bool, error) {
	/*
	db := config.GetMongoClient()
	collection := db.Database("speakup").Collection("users")

	var user models.User
	if err := collection.FindOne(ctx, bson.M{"id": userID}).Decode(&user); err != nil {
		return false, err
	}

	status := strings.ToLower(user.StripeStatus)
	return status == "active" || status == "trialing", nil
	*/
	return true, nil
}

func GetUsageCount(ctx context.Context, userID string) (int64, error) {
	db := config.GetMongoClient()
	collection := db.Database("speakup").Collection("ai_usage")

	dateKey := time.Now().UTC().Format("2006-01-02")
	filter := bson.M{"user_id": userID, "date": dateKey}

	var usage struct {
		Count int64 `bson:"count"`
	}
	if err := collection.FindOne(ctx, filter).Decode(&usage); err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, err
	}

	return usage.Count, nil
}

func CheckAndIncrementUsage(ctx context.Context, userID string, limit int64) (bool, error) {
	db := config.GetMongoClient()
	collection := db.Database("speakup").Collection("ai_usage")

	dateKey := time.Now().UTC().Format("2006-01-02")
	filter := bson.M{"user_id": userID, "date": dateKey}

	var usage struct {
		Count int64 `bson:"count"`
	}
	if err := collection.FindOne(ctx, filter).Decode(&usage); err != nil {
		if err != mongo.ErrNoDocuments {
			return false, err
		}
	}

	if usage.Count >= limit {
		return false, nil
	}

	update := bson.M{
		"$inc": bson.M{"count": 1},
		"$setOnInsert": bson.M{
			"user_id": userID,
			"date":    dateKey,
		},
	}
	_, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return false, err
	}

	return true, nil
}
