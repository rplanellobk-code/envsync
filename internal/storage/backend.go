// Package storage defines the interface and types for remote backends.
package storage

import "context"

// Backend defines the interface for reading and writing encrypted env data
// to a remote or local storage system.
type Backend interface {
	// Put stores encrypted data under the given key.
	Put(ctx context.Context, key string, data []byte) error

	// Get retrieves encrypted data stored under the given key.
	Get(ctx context.Context, key string) ([]byte, error)

	// List returns all keys available in the backend.
	List(ctx context.Context) ([]string, error)

	// Delete removes the data stored under the given key.
	Delete(ctx context.Context, key string) error
}

// ErrNotFound is returned when a key does not exist in the backend.
type ErrNotFound struct {
	Key string
}

func (e *ErrNotFound) Error() string {
	return "storage: key not found: " + e.Key
}

// IsNotFound reports whether err is an ErrNotFound error.
func IsNotFound(err error) bool {
	_, ok := err.(*ErrNotFound)
	return ok
}
