package s3

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"user/internal/logger"
)

var s3Client *s3.Client
var s3BucketName string
var s3Region string

// UploadImageFn points to the default S3 uploader, editable for testing mocks
var UploadImageFn = UploadToS3

// InitS3 initializes the S3 client using credentials
func InitS3(accessKey, secretKey, region, bucket string) error {
	s3BucketName = bucket
	s3Region = region

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	)
	if err != nil {
		return fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	s3Client = s3.NewFromConfig(cfg)
	logger.Logger.Info("AWS S3 client initialized successfully")
	return nil
}

// UploadToS3 uploads a file reader stream to AWS S3 and returns the public secure URL
func UploadToS3(ctx context.Context, file io.Reader, originalFilename, contentType string) (string, error) {
	if s3Client == nil {
		return "", fmt.Errorf("S3 client not initialized")
	}

	// Generate a unique file key to prevent namespace collisions and directory traversal
	ext := filepath.Ext(originalFilename)
	uniqueID := uuid.New().String()
	key := fmt.Sprintf("profile_pictures/%s%s", uniqueID, ext)

	// Perform standard PutObject upload
	_, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s3BucketName),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object to S3: %w", err)
	}

	// Generate public S3 URL
	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s3BucketName, s3Region, key)
	logger.Logger.Info("File uploaded successfully to S3", zap.String("url", fileURL))
	return fileURL, nil
}
