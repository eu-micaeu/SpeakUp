package repositories

import (
	"context"
	"speakup/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) error
	FindByID(ctx context.Context, id string) (models.User, error)
	FindByEmail(ctx context.Context, email string) (models.User, error)
	Update(ctx context.Context, id string, updateDoc bson.M) error
	Delete(ctx context.Context, id string) error
}

type MongoUserRepository struct {
	Collection *mongo.Collection
}

func NewMongoUserRepository(db *mongo.Database) *MongoUserRepository {
	return &MongoUserRepository{
		Collection: db.Collection("users"),
	}
}

func userIDFilter(id string) bson.M {
	return bson.M{
		"$or": []bson.M{
			{"_id": id},
			{"id": id},
		},
	}
}

func (r *MongoUserRepository) Create(ctx context.Context, user models.User) error {
	_, err := r.Collection.InsertOne(ctx, user)
	return err
}

func (r *MongoUserRepository) FindByID(ctx context.Context, id string) (models.User, error) {
	var user models.User
	err := r.Collection.FindOne(ctx, userIDFilter(id)).Decode(&user)
	if err == nil {
		user.Normalize()
	}
	return user, err
}

func (r *MongoUserRepository) FindByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := r.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err == nil {
		user.Normalize()
	}
	return user, err
}

func (r *MongoUserRepository) Update(ctx context.Context, id string, updateDoc bson.M) error {
	_, err := r.Collection.UpdateOne(ctx, userIDFilter(id), bson.M{"$set": updateDoc})
	return err
}

func (r *MongoUserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.Collection.DeleteOne(ctx, userIDFilter(id))
	return err
}

