package env

import (
	"encoding/json"
	"fmt"

	"envsync/internal/storage"
)

const accessPrefix = "access/"

// AccessStore persists AccessPolicy objects via a storage backend.
type AccessStore struct {
	backend storage.Backend
}

// NewAccessStore creates a new AccessStore backed by the given backend.
func NewAccessStore(b storage.Backend) *AccessStore {
	return &AccessStore{backend: b}
}

// Save serialises and stores the policy for its environment.
func (s *AccessStore) Save(policy *AccessPolicy) error {
	data, err := json.Marshal(policy)
	if err != nil {
		return fmt.Errorf("access: marshal policy: %w", err)
	}
	return s.backend.Put(accessKey(policy.Environment), data)
}

// Load retrieves the policy for the named environment.
func (s *AccessStore) Load(env string) (*AccessPolicy, error) {
	data, err := s.backend.Get(accessKey(env))
	if err != nil {
		return nil, fmt.Errorf("access: load policy for %q: %w", env, err)
	}
	var policy AccessPolicy
	if err := json.Unmarshal(data, &policy); err != nil {
		return nil, fmt.Errorf("access: unmarshal policy: %w", err)
	}
	return &policy, nil
}

// Delete removes the stored policy for the named environment.
func (s *AccessStore) Delete(env string) error {
	return s.backend.Delete(accessKey(env))
}

// ListEnvironments returns all environments that have a stored access policy.
func (s *AccessStore) ListEnvironments() ([]string, error) {
	keys, err := s.backend.List(accessPrefix)
	if err != nil {
		return nil, fmt.Errorf("access: list environments: %w", err)
	}
	envs := make([]string, 0, len(keys))
	for _, k := range keys {
		envs = append(envs, k[len(accessPrefix):])
	}
	return envs, nil
}

func accessKey(env string) string {
	return accessPrefix + env
}
