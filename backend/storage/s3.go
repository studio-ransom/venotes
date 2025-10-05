package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Storage implements Storage interface for AWS S3
type S3Storage struct {
	client   *s3.Client
	bucket   string
	basePath string
}

// S3Config holds S3 configuration
type S3Config struct {
	AccessKey  string
	SecretKey  string
	Endpoint   string
	Region     string
	Bucket     string
	BasePath   string
}

// NewS3Storage creates a new S3 storage instance
func NewS3Storage(cfg S3Config) (*S3Storage, error) {
	// Create AWS config
	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKey,
			cfg.SecretKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Override endpoint if provided (for S3-compatible services like MinIO)
	if cfg.Endpoint != "" {
		awsConfig.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL: cfg.Endpoint,
				}, nil
			},
		)
	}

	client := s3.NewFromConfig(awsConfig)

	return &S3Storage{
		client:   client,
		bucket:   cfg.Bucket,
		basePath: strings.Trim(cfg.BasePath, "/"),
	}, nil
}

// StoreFile stores a file in S3
func (s *S3Storage) StoreFile(ctx context.Context, filename string, content io.Reader) (string, error) {
	key := s.getKey(filename)
	
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   content,
	})
	
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}
	
	return key, nil
}

// GetFile retrieves a file from S3
func (s *S3Storage) GetFile(ctx context.Context, key string) (io.ReadCloser, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to get file from S3: %w", err)
	}
	
	return result.Body, nil
}

// DeleteFile deletes a file from S3
func (s *S3Storage) DeleteFile(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}
	
	return nil
}

// FileExists checks if a file exists in S3
func (s *S3Storage) FileExists(ctx context.Context, key string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	
	if err != nil {
		// Check if it's a "not found" error
		if strings.Contains(err.Error(), "NoSuchKey") || strings.Contains(err.Error(), "NotFound") {
			return false, nil
		}
		return false, fmt.Errorf("failed to check file existence in S3: %w", err)
	}
	
	return true, nil
}

// getKey constructs the full S3 key for a file
func (s *S3Storage) getKey(filename string) string {
	if s.basePath == "" {
		return filename
	}
	return s.basePath + "/" + filename
}
