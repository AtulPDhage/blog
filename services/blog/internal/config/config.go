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
	RedisURL           string
	UserService        string
	JWTSecret          string
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
		RedisURL:           os.Getenv("REDIS_URL"),
		UserService:        os.Getenv("USER_SERVICE"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
	}

	// Fail-close check: ensure essential variables are defined
	if cfg.Port == "" {
		cfg.Port = "5001" // default port for blog service
	}
	if cfg.DBURL == "" {
		return nil, fmt.Errorf("DB_URL environment variable is required")
	}
	if cfg.RabbitMQHost == "" {
		return nil, fmt.Errorf("Rabbitmq_Host environment variable is required")
	}
	if cfg.RedisURL == "" {
		return nil, fmt.Errorf("REDIS_URL environment variable is required")
	}
	if cfg.UserService == "" {
		return nil, fmt.Errorf("USER_SERVICE environment variable is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	return cfg, nil
}
