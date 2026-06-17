package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"user/internal/models"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id string) (*models.User, error)
	Create(ctx context.Context, name, email, image string) (*models.User, error)
	Update(ctx context.Context, id string, name, instagram, linkedin, facebook, bio string) (*models.User, error)
	UpdateProfilePicture(ctx context.Context, id string, imageURL string) (*models.User, error)
}

type MongoUserRepository struct {
	collection *mongo.Collection
}

// NewMongoUserRepository returns a new user repository instance
func NewMongoUserRepository() UserRepository {
	return &MongoUserRepository{
		collection: Database.Collection("users"),
	}
}

// FindByEmail searches MongoDB for a user matching the email address
func (r *MongoUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	user.SyncIDs()
	return &user, nil
}

// FindByID searches MongoDB for a user matching the string ObjectID
func (r *MongoUserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectID format: %w", err)
	}

	var user models.User
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}
	user.SyncIDs()
	return &user, nil
}

// Create inserts a new user record into MongoDB
func (r *MongoUserRepository) Create(ctx context.Context, name, email, image string) (*models.User, error) {
	now := time.Now()
	user := &models.User{
		BsonID:    primitive.NewObjectID(),
		Name:      name,
		Email:     email,
		Image:     image,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user record: %w", err)
	}

	user.SyncIDs()
	return user, nil
}

// Update updates name, bio, and social links in MongoDB, returning the modified document
func (r *MongoUserRepository) Update(ctx context.Context, id string, name, instagram, linkedin, facebook, bio string) (*models.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectID format: %w", err)
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"name":      name,
			"instagram": instagram,
			"linkedin":  linkedin,
			"facebook":  facebook,
			"bio":       bio,
			"updatedAt": now,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedUser models.User
	err = r.collection.FindOneAndUpdate(ctx, bson.M{"_id": oid}, update, opts).Decode(&updatedUser)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	updatedUser.SyncIDs()
	return &updatedUser, nil
}

// UpdateProfilePicture modifies the user's profile image URL in MongoDB, returning the modified document
func (r *MongoUserRepository) UpdateProfilePicture(ctx context.Context, id string, imageURL string) (*models.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectID: %w", err)
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"image":     imageURL,
			"updatedAt": now,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updatedUser models.User
	err = r.collection.FindOneAndUpdate(ctx, bson.M{"_id": oid}, update, opts).Decode(&updatedUser)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile picture: %w", err)
	}

	updatedUser.SyncIDs()
	return &updatedUser, nil
}
