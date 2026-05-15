package env

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"envsync/internal/storage"
)

const scopeKeyPrefix = "scopes/"

// ScopeStore persists and retrieves Scope definitions in a storage backend.
type ScopeStore struct {
	backend storage.Backend
}

// NewScopeStore creates a ScopeStore backed by the given storage.Backend.
func NewScopeStore(b storage.Backend) *ScopeStore {
	return &ScopeStore{backend: b}
}

func scopeKey(name string) string {
	return scopeKeyPrefix + name
}

// Save persists a Scope. Returns an error if the name is empty or the
// pattern is invalid.
func (ss *ScopeStore) Save(ctx context.Context, s *Scope) error {
	if s.Name == "" {
		return fmt.Errorf("scope name must not be empty")
	}
	// ensure pattern compiles before persisting
	if err := s.compile(); err != nil {
		return err
	}
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("scope: marshal: %w", err)
	}
	return ss.backend.Put(ctx, scopeKey(s.Name), data)
}

// Load retrieves a Scope by name.
func (ss *ScopeStore) Load(ctx context.Context, name string) (*Scope, error) {
	data, err := ss.backend.Get(ctx, scopeKey(name))
	if err != nil {
		return nil, err
	}
	var s Scope
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("scope: unmarshal: %w", err)
	}
	return &s, nil
}

// Delete removes a Scope by name.
func (ss *ScopeStore) Delete(ctx context.Context, name string) error {
	return ss.backend.Delete(ctx, scopeKey(name))
}

// List returns all stored Scope names.
func (ss *ScopeStore) List(ctx context.Context) ([]string, error) {
	keys, err := ss.backend.List(ctx, scopeKeyPrefix)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(keys))
	for _, k := range keys {
		names = append(names, strings.TrimPrefix(k, scopeKeyPrefix))
	}
	return names, nil
}
