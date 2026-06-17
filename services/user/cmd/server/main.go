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

	"user/internal/config"
	"user/internal/s3"
	"user/internal/db"
	"user/internal/google"
	"user/internal/handlers"
	"user/internal/logger"
	"user/internal/middleware"
	"user/internal/service"
)

func main() {
	// 1. Initialize Zap Logger
	logger.InitLogger()
	defer func() {
		_ = logger.Logger.Sync()
	}()

	logger.Logger.Info("🚀 Starting User Service in Go with Chi Router...")

	// 2. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Logger.Fatal("❌ Configuration error", zap.Error(err))
	}

	// 3. Connect MongoDB
	err = db.ConnectDB(cfg)
	if err != nil {
		logger.Logger.Fatal("❌ MongoDB connection error", zap.Error(err))
	}
	defer db.CloseDB()

	// 4. Initialize AWS S3 client
	err = s3.InitS3(cfg.AWSAccessKeyID, cfg.AWSSecretAccessKey, cfg.AWSRegion, cfg.AWSS3Bucket)
	if err != nil {
		logger.Logger.Fatal("❌ AWS S3 client initialization error", zap.Error(err))
	}

	// 5. Initialize Google Client
	googleClient := google.NewGoogleClient(cfg.GoogleClientID, cfg.GoogleClientSecret)

	// 6. Instantiate Service-Repository pattern layers
	repo := db.NewMongoUserRepository()
	userService := service.NewUserService(repo, googleClient, cfg.JWTSecret)
	handler := handlers.NewUserHandler(userService)

	// 7. Setup Chi Router
	r := chi.NewRouter()

	// Apply global middlewares
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RateLimitMiddleware())

	// Register routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/login", handler.LoginUser)
		r.Get("/user/{id}", handler.GetUserProfile)

		// Authenticated routes group
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(cfg.JWTSecret))

			r.Get("/me", handler.MyProfile)
			r.Post("/user/update", handler.UpdateUser)

			// Size limited file uploads
			r.With(middleware.MaxBodySizeMiddleware(10 << 20)).Post("/user/update/pic", handler.UpdateProfilePicture)
		})
	})

	// Local security binding check (must bind to localhost during local/test executions)
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

	// 8. Start Server with Graceful Shutdown support
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
