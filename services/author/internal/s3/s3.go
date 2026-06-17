package s3

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"author/internal/logger"
)

var s3Client *s3.Client
var s3BucketName string
var s3Region string
var s3Endpoint string

var UploadImageFn = UploadToS3

// InitS3 initializes the S3 client
func InitS3(accessKey, secretKey, region, bucket string) error {
	s3BucketName = bucket
	s3Region = region
	s3Endpoint = os.Getenv("AWS_ENDPOINT")

	var cfg aws.Config
	var err error

	if s3Endpoint != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:               s3Endpoint,
				HostnameImmutable: true,
			}, nil
		})
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
			config.WithEndpointResolverWithOptions(customResolver),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		)
	}

	if err != nil {
		return fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		if s3Endpoint != "" || os.Getenv("AWS_S3_FORCE_PATH_STYLE") == "true" {
			o.UsePathStyle = true
		}
	})

	// Auto-create bucket for local development / MinIO setup
	_, _ = s3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})

	// Set public read policy on the bucket for local development / MinIO setup
	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Sid": "PublicReadGetObject",
				"Effect": "Allow",
				"Principal": "*",
				"Action": "s3:GetObject",
				"Resource": "arn:aws:s3:::%s/*"
			}
		]
	}`, bucket)
	_, _ = s3Client.PutBucketPolicy(context.TODO(), &s3.PutBucketPolicyInput{
		Bucket: aws.String(bucket),
		Policy: aws.String(policy),
	})

	logger.Logger.Info("AWS S3 client initialized successfully")
	return nil
}

// UploadToS3 uploads a file reader stream to AWS S3 and returns the public secure URL of the file
func UploadToS3(ctx context.Context, file io.Reader, originalFilename, contentType string) (string, error) {
	if s3Client == nil {
		return "", fmt.Errorf("S3 client not initialized")
	}

	// Generate a unique file key to prevent namespace collisions and directory traversal
	ext := filepath.Ext(originalFilename)
	uniqueID := uuid.New().String()
	key := fmt.Sprintf("blogs/%s%s", uniqueID, ext)

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
	var fileURL string
	if s3Endpoint != "" {
		fileURL = fmt.Sprintf("%s/%s/%s", s3Endpoint, s3BucketName, key)
	} else {
		fileURL = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s3BucketName, s3Region, key)
	}
	logger.Logger.Info("File uploaded successfully to S3/MinIO", zap.String("url", fileURL))
	return fileURL, nil
}
