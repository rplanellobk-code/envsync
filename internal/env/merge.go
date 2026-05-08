// Package env provides utilities for parsing, serializing, diffing, and merging .env files.
package env

import "fmt"

// ConflictStrategy determines how conflicts are resolved during a merge.
type ConflictStrategy int

const (
	// PreferLocal keeps the local value on conflict.
	PreferLocal ConflictStrategy = iota
	// PreferRemote keeps the remote value on conflict.
	PreferRemote
	// ErrorOnConflict returns an error if any conflict is detected.
	ErrorOnConflict
)

// ConflictError is returned when a merge conflict is detected and the strategy
// is ErrorOnConflict.
type ConflictError struct {
	Key   string
	Local string
	Remote string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("merge conflict on key %q: local=%q remote=%q", e.Key, e.Local, e.Remote)
}

// Merge combines local and remote env maps using the provided ConflictStrategy.
// Keys present only in one side are always included. Keys present in both sides
// with different values are resolved according to strategy.
func Merge(local, remote map[string]string, strategy ConflictStrategy) (map[string]string, error) {
	result := make(map[string]string, len(local))

	// Copy all local entries into result.
	for k, v := range local {
		result[k] = v
	}

	for k, remoteVal := range remote {
		localVal, exists := result[k]
		if !exists {
			// Key only in remote — always add it.
			result[k] = remoteVal
			continue
		}

		if localVal == remoteVal {
			// No conflict.
			continue
		}

		// Conflict detected.
		switch strategy {
		case PreferLocal:
			// Keep local value already in result.
		case PreferRemote:
			result[k] = remoteVal
		case ErrorOnConflict:
			return nil, &ConflictError{Key: k, Local: localVal, Remote: remoteVal}
		}
	}

	return result, nil
}
