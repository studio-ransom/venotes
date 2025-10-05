package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

// LocalStorage implements Storage interface for local file system
type LocalStorage struct {
	basePath string
}

// NewLocalStorage creates a new local storage instance
func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{
		basePath: basePath,
	}
}

// StoreFile stores a file locally
func (l *LocalStorage) StoreFile(ctx context.Context, filename string, content io.Reader) (string, error) {
	// Ensure directory exists
	if err := os.MkdirAll(l.basePath, 0755); err != nil {
		return "", err
	}
	
	filePath := filepath.Join(l.basePath, filename)
	
	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	
	// Copy content
	_, err = io.Copy(file, content)
	if err != nil {
		os.Remove(filePath) // Clean up on error
		return "", err
	}
	
	return filePath, nil
}

// GetFile retrieves a file from local storage
func (l *LocalStorage) GetFile(ctx context.Context, path string) (io.ReadCloser, error) {
	return os.Open(path)
}

// DeleteFile deletes a file from local storage
func (l *LocalStorage) DeleteFile(ctx context.Context, path string) error {
	return os.Remove(path)
}

// FileExists checks if a file exists in local storage
func (l *LocalStorage) FileExists(ctx context.Context, path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}
