package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"

	"blog/internal/db"
	"blog/internal/logger"
	"blog/internal/redis"
)

// StartCacheConsumer connects to the queue and listens for invalidation events
func StartCacheConsumer() {
	if Channel == nil {
		logger.Logger.Error("RabbitMQ channel is not established, cannot start consumer")
		return
	}

	queueName := "cache-invalidation"
	_, err := Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		logger.Logger.Error("Failed to declare cache-invalidation queue", zap.Error(err))
		return
	}

	msgs, err := Channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack (explicit ack)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		logger.Logger.Error("Failed to register a consumer", zap.Error(err))
		return
	}

	logger.Logger.Info("Blog Cache Consumer connected to RabbitMQ, waiting for messages...")

	go func() {
		for d := range msgs {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			var msg CacheInvalidationMessage
			err := json.Unmarshal(d.Body, &msg)
			if err != nil {
				logger.Logger.Error("Error parsing cache invalidation message", zap.Error(err))
				// Reject message without requeue if invalid format
				_ = d.Nack(false, false)
				cancel()
				continue
			}

			logger.Logger.Info("Blog Service Received cache invalidation message", zap.Any("message", msg))

			if msg.Action == "invalidateCache" && len(msg.Keys) > 0 {
				repo := db.NewPostgresBlogRepository()
				for _, pattern := range msg.Keys {
					deleted, err := redis.InvalidateKeysByPattern(ctx, pattern)
					if err != nil {
						logger.Logger.Error("Failed to invalidate keys by pattern", zap.String("pattern", pattern), zap.Error(err))
						continue
					}

					if deleted > 0 {
						// Refill default empty search category list key: blogs:::12:0
						cacheKey := "blogs:::12:0"
						blogs, err := repo.GetAllBlogs(ctx, "", "", 12, 0)
						if err != nil {
							logger.Logger.Error("Failed to fetch blogs from database for cache refill", zap.Error(err))
							continue
						}

						blogsJSON, err := json.Marshal(blogs)
						if err != nil {
							logger.Logger.Error("Failed to marshal blogs for cache refill", zap.Error(err))
							continue
						}

						err = redis.Set(ctx, cacheKey, string(blogsJSON), 3600*time.Second)
						if err != nil {
							logger.Logger.Error("Failed to update cache key after invalidation", zap.String("key", cacheKey), zap.Error(err))
							continue
						}

						logger.Logger.Info("Blog Service Refreshed cache for key", zap.String("key", cacheKey))
					}
				}
			}

			_ = d.Ack(false)
			cancel()
		}
	}()
}
