package storage

import (
	"fmt"
	"os"
)

// StorageFactory creates storage instances based on configuration
type StorageFactory struct{}

// NewStorageFactory creates a new storage factory
func NewStorageFactory() *StorageFactory {
	return &StorageFactory{}
}

// CreateStorage creates a storage instance based on environment variables
func (f *StorageFactory) CreateStorage() (Storage, error) {
	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		storageType = string(StorageTypeLocal) // Default to local
	}

	switch StorageType(storageType) {
	case StorageTypeLocal:
		basePath := os.Getenv("STORAGE_LOCAL_PATH")
		if basePath == "" {
			basePath = "data/uploads"
		}
		return NewLocalStorage(basePath), nil

	case StorageTypeS3:
		cfg := S3Config{
			AccessKey: os.Getenv("S3_ACCESS_KEY"),
			SecretKey: os.Getenv("S3_SECRET_KEY"),
			Endpoint:  os.Getenv("S3_ENDPOINT"),
			Region:    os.Getenv("S3_REGION"),
			Bucket:    os.Getenv("S3_BUCKET"),
			BasePath:  os.Getenv("S3_BASE_PATH"),
		}

		// Validate required S3 configuration
		if cfg.AccessKey == "" {
			return nil, fmt.Errorf("S3_ACCESS_KEY environment variable is required")
		}
		if cfg.SecretKey == "" {
			return nil, fmt.Errorf("S3_SECRET_KEY environment variable is required")
		}
		if cfg.Bucket == "" {
			return nil, fmt.Errorf("S3_BUCKET environment variable is required")
		}
		if cfg.Region == "" {
			cfg.Region = "us-east-1" // Default region
		}

		return NewS3Storage(cfg)

	default:
		return nil, fmt.Errorf("unsupported storage type: %s", storageType)
	}
}
