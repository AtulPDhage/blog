package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	MongoURI           string
	DBName             string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string
	AWSS3Bucket        string
	GoogleClientID     string
	GoogleClientSecret string
	JWTSecret          string
}

// LoadConfig reads configuration from environment variables and checks required fields
func LoadConfig() (*Config, error) {
	// Optional load for local development
	_ = godotenv.Load()

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "MasterDB"
	}

	cfg := &Config{
		Port:               os.Getenv("PORT"),
		MongoURI:           os.Getenv("MONGO_URI"),
		DBName:             dbName,
		AWSAccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWSRegion:          os.Getenv("AWS_REGION"),
		AWSS3Bucket:        os.Getenv("AWS_S3_BUCKET"),
		GoogleClientID:     os.Getenv("Google_Client_Id"),
		GoogleClientSecret: os.Getenv("Google_client_Secret"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
	}

	if cfg.Port == "" {
		cfg.Port = "5002"
	}

	// Fail-close check: ensure essential variables are defined
	if cfg.MongoURI == "" {
		return nil, fmt.Errorf("MONGO_URI environment variable is required")
	}
	if cfg.AWSAccessKeyID == "" || cfg.AWSSecretAccessKey == "" || cfg.AWSRegion == "" || cfg.AWSS3Bucket == "" {
		return nil, fmt.Errorf("AWS S3 settings (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_REGION, AWS_S3_BUCKET) are required")
	}
	if cfg.GoogleClientID == "" || cfg.GoogleClientSecret == "" {
		return nil, fmt.Errorf("google OAuth credentials (Google_Client_Id, Google_client_Secret) are required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	return cfg, nil
}
