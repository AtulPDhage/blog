package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"blog/internal/config"
	"blog/internal/db"
	"blog/internal/handlers"
	"blog/internal/logger"
	"blog/internal/middleware"
	"blog/internal/rabbitmq"
	"blog/internal/redis"
	"blog/internal/service"
)

func main() {
	// 1. Initialize Zap Logger
	logger.InitLogger()
	defer func() {
		_ = logger.Logger.Sync()
	}()

	logger.Logger.Info("🚀 Starting Blog Service in Go with Chi Router...")

	// 2. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Logger.Fatal("❌ Configuration error", zap.Error(err))
	}

	// 3. Connect Database using Config pooling settings
	err = db.ConnectDB(cfg)
	if err != nil {
		logger.Logger.Fatal("❌ Database connection error", zap.Error(err))
	}
	defer db.Pool.Close()

	// 4. Initialize Tables
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = db.InitTables(ctx)
	if err != nil {
		logger.Logger.Fatal("❌ Table initialization error", zap.Error(err))
	}

	// 5. Connect Redis
	err = redis.ConnectRedis(cfg.RedisURL)
	if err != nil {
		logger.Logger.Fatal("❌ Redis connection error", zap.Error(err))
	}
	defer func() {
		if redis.Client != nil {
			_ = redis.Client.Close()
		}
	}()

	// 6. Connect RabbitMQ and start Cache Consumer
	err = rabbitmq.ConnectRabbitMQ(cfg.RabbitMQHost, cfg.RabbitMQUsername, cfg.RabbitMQPassword)
	if err != nil {
		logger.Logger.Warn("⚠️ RabbitMQ connection warning (running without RabbitMQ)", zap.Error(err))
	} else {
		defer rabbitmq.CloseRabbitMQ()
		// Start consuming invalidation events in the background
		rabbitmq.StartCacheConsumer()
	}

	// 7. Instantiate Service-Repository pattern layers
	repo := db.NewPostgresBlogRepository()
	blogService := service.NewBlogService(repo, cfg.UserService)
	handler := handlers.NewBlogHandler(blogService)

	// 8. Setup Chi Router
	r := chi.NewRouter()

	// Apply global middlewares
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RateLimitMiddleware())

	// Register routes with route groups and specific middleware chains
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/blog/all", handler.GetAllBlogs)
		r.Get("/blog/{id}", handler.GetSingleBlog)
		r.Get("/comment/{id}", handler.GetAllComments)

		// Authenticated & size-limited routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.MaxBodySizeMiddleware(10 << 20)) // 10MB request size limit
			r.Use(middleware.AuthMiddleware(cfg.JWTSecret))

			r.Post("/comment/{id}", handler.AddComment)
			r.Delete("/comment/{commentid}", handler.DeleteComment)
			r.Post("/save/{blogid}", handler.SaveBlog)
			r.Get("/blogs/saved/all", handler.GetSavedBlogs)
		})
	})

	// Local security binding check
	// Listen on 127.0.0.1 by default for secure local testing (MUST NOT listen on 0.0.0.0 during test/local)
	host := "127.0.0.1"
	if os.Getenv("BIND_ADDR") != "" {
		host = os.Getenv("BIND_ADDR")
	}
	addr := fmt.Sprintf("%s:%s", host, cfg.Port)

	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 9. Start Server with Graceful Shutdown support
	go func() {
		logger.Logger.Info("📡 Server is running", zap.String("url", fmt.Sprintf("http://%s", addr)))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Logger.Fatal("❌ HTTP server listen error", zap.Error(err))
		}
	}()

	// Wait for termination signals
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)
	<-shutdownChan

	logger.Logger.Info("🛑 Shutting down server gracefully...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Logger.Error("⚠️ Graceful shutdown failed", zap.Error(err))
	} else {
		logger.Logger.Info("✅ Server stopped gracefully")
	}
}
