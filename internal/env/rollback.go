package env

import (
	"errors"
	"fmt"
)

// ErrNoSnapshots is returned when no snapshots exist for the given environment.
var ErrNoSnapshots = errors.New("no snapshots found for environment")

// Rollback restores the environment in the vault to the state captured in the
// given snapshot. It records the rollback action in the audit log when one is
// provided (audit may be nil).
func Rollback(
	v *VaultWithAudit,
	store *SnapshotStore,
	environment string,
	snapshotID string,
	author string,
) error {
	if environment == "" {
		return errors.New("environment must not be empty")
	}
	if snapshotID == "" {
		return errors.New("snapshotID must not be empty")
	}

	snaps, err := store.List(environment)
	if err != nil {
		return fmt.Errorf("rollback: list snapshots: %w", err)
	}
	if len(snaps) == 0 {
		return ErrNoSnapshots
	}

	var target *Snapshot
	for _, s := range snaps {
		if s.ID == snapshotID {
			copy := s
			target = &copy
			break
		}
	}
	if target == nil {
		return fmt.Errorf("rollback: snapshot %q not found for environment %q", snapshotID, environment)
	}

	if err := v.Push(environment, target.Data, author); err != nil {
		return fmt.Errorf("rollback: push snapshot data: %w", err)
	}

	return nil
}

// RollbackToLatest is a convenience wrapper that rolls back to the most recent
// snapshot preceding the current HEAD (i.e. the second-latest snapshot).
func RollbackToLatest(
	v *VaultWithAudit,
	store *SnapshotStore,
	environment string,
	author string,
) error {
	snaps, err := store.List(environment)
	if err != nil {
		return fmt.Errorf("rollback latest: %w", err)
	}
	if len(snaps) < 2 {
		return fmt.Errorf("rollback latest: need at least 2 snapshots, have %d", len(snaps))
	}
	// List returns entries newest-first; index 1 is the previous version.
	previous := snaps[1]
	return Rollback(v, store, environment, previous.ID, author)
}
