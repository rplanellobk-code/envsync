package storage

import (
	"fmt"
	"sort"
	"sync"
)

// MemoryBackend is an in-process Backend useful for testing.
type MemoryBackend struct {
	mu   sync.RWMutex
	data map[string][]byte
}

// NewMemoryBackend returns an empty MemoryBackend.
func NewMemoryBackend() *MemoryBackend {
	return &MemoryBackend{data: make(map[string][]byte)}
}

// Put stores value under key, overwriting any previous value.
func (m *MemoryBackend) Put(key string, value []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	buf := make([]byte, len(value))
	copy(buf, value)
	m.data[key] = buf
	return nil
}

// Get retrieves the value stored under key.
func (m *MemoryBackend) Get(key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[key]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, key)
	}
	buf := make([]byte, len(v))
	copy(buf, v)
	return buf, nil
}

// Delete removes the entry for key. Returns ErrNotFound if absent.
func (m *MemoryBackend) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.data[key]; !ok {
		return fmt.Errorf("%w: %s", ErrNotFound, key)
	}
	delete(m.data, key)
	return nil
}

// List returns all keys with the given prefix, sorted lexicographically.
func (m *MemoryBackend) List(prefix string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var keys []string
	for k := range m.data {
		if len(prefix) == 0 || len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return keys, nil
}
