package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"user/internal/config"
	"user/internal/logger"
)

var Client *mongo.Client
var Database *mongo.Database

// ConnectDB establishes a connection pool to MongoDB
func ConnectDB(cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoURI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping database to verify connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	Client = client
	Database = client.Database(cfg.DBName)

	logger.Logger.Info("Connected to MongoDB successfully", zap.String("db_name", cfg.DBName))
	return nil
}

// CloseDB disconnects the MongoDB client
func CloseDB() {
	if Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = Client.Disconnect(ctx)
		logger.Logger.Info("Disconnected from MongoDB successfully")
	}
}

