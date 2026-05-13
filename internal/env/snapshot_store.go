package env

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/user/envsync/internal/storage"
)

const snapshotPrefix = "snapshots/"

// SnapshotStore persists and retrieves Snapshots via a storage backend.
type SnapshotStore struct {
	backend storage.Backend
}

// NewSnapshotStore creates a SnapshotStore backed by the given storage.Backend.
func NewSnapshotStore(b storage.Backend) *SnapshotStore {
	return &SnapshotStore{backend: b}
}

// Save persists a Snapshot. The key is derived from environment + timestamp.
func (ss *SnapshotStore) Save(snap *Snapshot) error {
	data, err := json.Marshal(snap)
	if err != nil {
		return fmt.Errorf("snapshot marshal: %w", err)
	}
	key := fmt.Sprintf("%s%s/%d", snapshotPrefix, snap.Environment, snap.CreatedAt.UnixNano())
	return ss.backend.Put(key, data)
}

// Latest retrieves the most recently saved Snapshot for the given environment.
// Returns storage.ErrNotFound (via IsNotFound) if none exist.
func (ss *SnapshotStore) Latest(environment string) (*Snapshot, error) {
	keys, err := ss.backend.List(snapshotPrefix + environment + "/")
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return nil, fmt.Errorf("no snapshots for environment %q", environment)
	}
	// Keys are stored with UnixNano suffix; lexicographic last == most recent.
	latest := keys[0]
	for _, k := range keys[1:] {
		if k > latest {
			latest = k
		}
	}
	data, err := ss.backend.Get(latest)
	if err != nil {
		return nil, err
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot unmarshal: %w", err)
	}
	return &snap, nil
}

// ListEnvironments returns the unique environment names that have snapshots.
func (ss *SnapshotStore) ListEnvironments() ([]string, error) {
	keys, err := ss.backend.List(snapshotPrefix)
	if err != nil {
		return nil, err
	}
	seen := make(map[string]struct{})
	for _, k := range keys {
		trimmed := strings.TrimPrefix(k, snapshotPrefix)
		parts := strings.SplitN(trimmed, "/", 2)
		if len(parts) > 0 && parts[0] != "" {
			seen[parts[0]] = struct{}{}
		}
	}
	envs := make([]string, 0, len(seen))
	for e := range seen {
		envs = append(envs, e)
	}
	return envs, nil
}
