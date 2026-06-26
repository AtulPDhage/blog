package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DBURL              string
	DBMaxConns         int32
	DBMinConns         int32
	DBMaxConnLifetime  time.Duration
	DBMaxConnIdleTime  time.Duration
	RabbitMQHost       string
	RabbitMQUsername   string
	RabbitMQPassword   string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string
	AWSS3Bucket        string
	JWTSecret          string
	GeminiAPIKey       string
}

func LoadConfig() (*Config, error) {
	// Load .env if it exists (for local development)
	_ = godotenv.Load()

	// Parse database pool settings with default fallbacks
	maxConns := int32(25)
	if val := os.Getenv("DB_MAX_CONNS"); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			maxConns = int32(i)
		}
	}

	minConns := int32(5)
	if val := os.Getenv("DB_MIN_CONNS"); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			minConns = int32(i)
		}
	}

	maxLifetime := 30 * time.Minute
	if val := os.Getenv("DB_MAX_CONN_LIFETIME"); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			maxLifetime = d
		}
	}

	maxIdleTime := 15 * time.Minute
	if val := os.Getenv("DB_MAX_CONN_IDLE_TIME"); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			maxIdleTime = d
		}
	}

	cfg := &Config{
		Port:               os.Getenv("PORT"),
		DBURL:              os.Getenv("DB_URL"),
		DBMaxConns:         maxConns,
		DBMinConns:         minConns,
		DBMaxConnLifetime:  maxLifetime,
		DBMaxConnIdleTime:  maxIdleTime,
		RabbitMQHost:       os.Getenv("Rabbitmq_Host"),
		RabbitMQUsername:   os.Getenv("Rabbitmq_Username"),
		RabbitMQPassword:   os.Getenv("Rabbitmq_Password"),
		AWSAccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWSRegion:          os.Getenv("AWS_REGION"),
		AWSS3Bucket:        os.Getenv("AWS_S3_BUCKET"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		GeminiAPIKey:       os.Getenv("GEMINI_API_KEY"),
	}

	// Fail-close check: ensure essential variables are defined
	if cfg.Port == "" {
		cfg.Port = "5000" // default port
	}
	if cfg.DBURL == "" {
		return nil, fmt.Errorf("DB_URL environment variable is required")
	}
	if cfg.RabbitMQHost == "" {
		return nil, fmt.Errorf("rabbitmq_Host environment variable is required")
	}
	if cfg.AWSAccessKeyID == "" || cfg.AWSSecretAccessKey == "" || cfg.AWSRegion == "" || cfg.AWSS3Bucket == "" {
		return nil, fmt.Errorf("AWS/S3 configuration environment variables (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_REGION, AWS_S3_BUCKET) are required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}
	if cfg.GeminiAPIKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable is required")
	}

	return cfg, nil
}
