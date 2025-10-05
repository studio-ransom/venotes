package storage

import (
	"context"
	"io"
)

// Storage interface defines methods for file storage operations
type Storage interface {
	// StoreFile stores a file and returns the path/key
	StoreFile(ctx context.Context, filename string, content io.Reader) (string, error)
	
	// GetFile retrieves a file by path/key
	GetFile(ctx context.Context, path string) (io.ReadCloser, error)
	
	// DeleteFile deletes a file by path/key
	DeleteFile(ctx context.Context, path string) error
	
	// FileExists checks if a file exists
	FileExists(ctx context.Context, path string) (bool, error)
}

// StorageType represents the type of storage backend
type StorageType string

const (
	StorageTypeLocal StorageType = "local"
	StorageTypeS3    StorageType = "s3"
)
