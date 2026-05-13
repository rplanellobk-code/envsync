package env

import (
	"fmt"
	"time"
)

// Snapshot represents a point-in-time capture of an environment's key-value pairs.
type Snapshot struct {
	Environment string            `json:"environment"`
	CreatedAt   time.Time         `json:"created_at"`
	Values      map[string]string `json:"values"`
}

// NewSnapshot creates a Snapshot from a parsed env map.
func NewSnapshot(environment string, values map[string]string) *Snapshot {
	copy := make(map[string]string, len(values))
	for k, v := range values {
		copy[k] = v
	}
	return &Snapshot{
		Environment: environment,
		CreatedAt:   time.Now().UTC(),
		Values:      copy,
	}
}

// DiffFrom returns the diff between this snapshot and a newer snapshot.
// The receiver is treated as the "before" state.
func (s *Snapshot) DiffFrom(newer *Snapshot) []DiffEntry {
	return Diff(s.Values, newer.Values)
}

// String returns a human-readable summary of the snapshot.
func (s *Snapshot) String() string {
	return fmt.Sprintf("Snapshot[env=%s keys=%d created=%s]",
		s.Environment, len(s.Values), s.CreatedAt.Format(time.RFC3339))
}
