package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"blog/internal/logger"
)

var Client *redis.Client

// ConnectRedis establishes a connection pool to Redis
func ConnectRedis(redisURL string) error {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	Client = redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	logger.Logger.Info("Connected to Redis successfully")
	return nil
}

// Get retrieves a string key value from Redis
func Get(ctx context.Context, key string) (string, error) {
	if Client == nil {
		return "", fmt.Errorf("Redis client is not initialized")
	}
	return Client.Get(ctx, key).Result()
}

// Set saves a key value pair to Redis with a specific expiration TTL
func Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	if Client == nil {
		return fmt.Errorf("Redis client is not initialized")
	}
	return Client.Set(ctx, key, value, ttl).Err()
}

// InvalidateKeysByPattern deletes all keys matching a specific pattern (e.g. blogs:*)
func InvalidateKeysByPattern(ctx context.Context, pattern string) (int, error) {
	if Client == nil {
		return 0, fmt.Errorf("Redis client is not initialized")
	}

	keys, err := Client.Keys(ctx, pattern).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to scan keys for pattern: %w", err)
	}

	deletedCount := len(keys)
	if deletedCount > 0 {
		err = Client.Del(ctx, keys...).Err()
		if err != nil {
			return 0, fmt.Errorf("failed to delete invalidation keys: %w", err)
		}
		logger.Logger.Info("Invalidated keys in Redis", zap.String("pattern", pattern), zap.Int("deleted_count", deletedCount))
	}

	return deletedCount, nil
}

