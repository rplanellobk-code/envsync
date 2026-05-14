package env

import (
	"encoding/json"
	"fmt"
	"time"

	"envsync/internal/storage"
)

// Lock represents an advisory lock on an environment.
type Lock struct {
	Environment string    `json:"environment"`
	Owner       string    `json:"owner"`
	AcquiredAt  time.Time `json:"acquired_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// IsExpired reports whether the lock has passed its expiry time.
func (l *Lock) IsExpired() bool {
	return time.Now().After(l.ExpiresAt)
}

// LockStore manages advisory locks stored in a backend.
type LockStore struct {
	backend storage.Backend
}

// NewLockStore creates a LockStore backed by the given storage backend.
func NewLockStore(b storage.Backend) *LockStore {
	return &LockStore{backend: b}
}

func lockKey(env string) string {
	return fmt.Sprintf("locks/%s.json", env)
}

// Acquire attempts to acquire a lock for the given environment.
// It returns an error if a non-expired lock already exists.
func (s *LockStore) Acquire(env, owner string, ttl time.Duration) (*Lock, error) {
	existing, err := s.Current(env)
	if err == nil && !existing.IsExpired() {
		return nil, fmt.Errorf("environment %q is locked by %q until %s",
			env, existing.Owner, existing.ExpiresAt.Format(time.RFC3339))
	}

	now := time.Now().UTC()
	lock := &Lock{
		Environment: env,
		Owner:       owner,
		AcquiredAt:  now,
		ExpiresAt:   now.Add(ttl),
	}

	data, err := json.Marshal(lock)
	if err != nil {
		return nil, fmt.Errorf("lock: marshal: %w", err)
	}
	if err := s.backend.Put(lockKey(env), data); err != nil {
		return nil, fmt.Errorf("lock: put: %w", err)
	}
	return lock, nil
}

// Release removes the lock for the given environment.
// Only the owner that holds the lock may release it.
func (s *LockStore) Release(env, owner string) error {
	existing, err := s.Current(env)
	if err != nil {
		return fmt.Errorf("lock: release: %w", err)
	}
	if existing.Owner != owner {
		return fmt.Errorf("lock: release: %q does not own lock for %q", owner, env)
	}
	if err := s.backend.Delete(lockKey(env)); err != nil {
		return fmt.Errorf("lock: delete: %w", err)
	}
	return nil
}

// Current returns the active lock for the environment, if any.
func (s *LockStore) Current(env string) (*Lock, error) {
	data, err := s.backend.Get(lockKey(env))
	if err != nil {
		return nil, err
	}
	var lock Lock
	if err := json.Unmarshal(data, &lock); err != nil {
		return nil, fmt.Errorf("lock: unmarshal: %w", err)
	}
	return &lock, nil
}
