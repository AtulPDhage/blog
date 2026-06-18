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

	"author/internal/config"
	"author/internal/db"
	"author/internal/handlers"
	"author/internal/logger"
	"author/internal/middleware"
	"author/internal/rabbitmq"
	"author/internal/s3"
	"author/internal/service"
	"author/internal/swagger"
)

// InjectGeminiAPIKey injects the Gemini API Key into the request context
func InjectGeminiAPIKey(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "GeminiAPIKey", apiKey)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func main() {
	// 1. Initialize Zap Logger
	logger.InitLogger()
	defer func() {
		_ = logger.Logger.Sync()
	}()

	logger.Logger.Info("🚀 Starting Author Service in Go with Chi Router...")

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

	// 5. Connect RabbitMQ
	err = rabbitmq.ConnectRabbitMQ(cfg.RabbitMQHost, cfg.RabbitMQUsername, cfg.RabbitMQPassword)
	if err != nil {
		logger.Logger.Warn("⚠️ RabbitMQ connection warning (running without RabbitMQ)", zap.Error(err))
	} else {
		defer rabbitmq.CloseRabbitMQ()
	}

	// 6. Connect S3
	err = s3.InitS3(cfg.AWSAccessKeyID, cfg.AWSSecretAccessKey, cfg.AWSRegion, cfg.AWSS3Bucket)
	if err != nil {
		logger.Logger.Fatal("❌ S3 initialization error", zap.Error(err))
	}

	// 7. Instantiate Service-Repository pattern layers
	repo := db.NewPostgresBlogRepository(db.Pool)
	blogService := service.NewBlogService(repo)
	handler := handlers.NewBlogHandler(blogService)

	// 8. Setup Chi Router
	r := chi.NewRouter()

	// Apply global middlewares
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RateLimitMiddleware())
	r.Use(InjectGeminiAPIKey(cfg.GeminiAPIKey))

	// Register Swagger routes
	swagger.RegisterRoutes(r)

	// Register routes with route groups and specific middleware chains
	r.Route("/api/v1", func(r chi.Router) {
		r.With(middleware.MaxBodySizeMiddleware(10<<20), middleware.AuthMiddleware(cfg.JWTSecret)).Post("/blog/new", handler.CreateBlog)
		r.With(middleware.MaxBodySizeMiddleware(10<<20), middleware.AuthMiddleware(cfg.JWTSecret)).Post("/blog/{id}", handler.UpdateBlog)
		r.With(middleware.AuthMiddleware(cfg.JWTSecret)).Delete("/blog/{id}", handler.DeleteBlog)

		r.Post("/ai/title", handler.AITitle)
		r.Post("/ai/description", handler.AIDescription)
		r.Post("/ai/blog", handler.AIBlog)
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
