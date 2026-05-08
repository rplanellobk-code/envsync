package storage

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

// FileBackend is a Backend implementation that stores data in a local directory.
// Each key maps to a file inside the base directory.
type FileBackend struct {
	BaseDir string
}

// NewFileBackend creates a FileBackend rooted at baseDir.
// The directory is created if it does not exist.
func NewFileBackend(baseDir string) (*FileBackend, error) {
	if err := os.MkdirAll(baseDir, 0700); err != nil {
		return nil, err
	}
	return &FileBackend{BaseDir: baseDir}, nil
}

func (f *FileBackend) path(key string) string {
	return filepath.Join(f.BaseDir, key+".enc")
}

// Put writes data to a file named <key>.enc inside BaseDir.
func (f *FileBackend) Put(_ context.Context, key string, data []byte) error {
	return os.WriteFile(f.path(key), data, 0600)
}

// Get reads the file for the given key.
func (f *FileBackend) Get(_ context.Context, key string) ([]byte, error) {
	data, err := os.ReadFile(f.path(key))
	if os.IsNotExist(err) {
		return nil, &ErrNotFound{Key: key}
	}
	return data, err
}

// List returns all keys found in BaseDir (files ending in .enc).
func (f *FileBackend) List(_ context.Context) ([]string, error) {
	entries, err := os.ReadDir(f.BaseDir)
	if err != nil {
		return nil, err
	}
	var keys []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".enc") {
			keys = append(keys, strings.TrimSuffix(e.Name(), ".enc"))
		}
	}
	return keys, nil
}

// Delete removes the file for the given key.
func (f *FileBackend) Delete(_ context.Context, key string) error {
	err := os.Remove(f.path(key))
	if os.IsNotExist(err) {
		return &ErrNotFound{Key: key}
	}
	return err
}
