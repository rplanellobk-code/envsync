package env

import (
	"testing"
	"time"

	"github.com/user/envsync/internal/storage"
)

func newSnapshotStore(t *testing.T) *SnapshotStore {
	t.Helper()
	tmpDir := t.TempDir()
	b, err := storage.NewFileBackend(tmpDir)
	if err != nil {
		t.Fatalf("NewFileBackend: %v", err)
	}
	return NewSnapshotStore(b)
}

func TestNewSnapshot(t *testing.T) {
	values := map[string]string{"KEY": "val", "FOO": "bar"}
	snap := NewSnapshot("production", values)
	if snap.Environment != "production" {
		t.Errorf("expected environment 'production', got %q", snap.Environment)
	}
	if len(snap.Values) != 2 {
		t.Errorf("expected 2 values, got %d", len(snap.Values))
	}
	// Ensure it's a copy
	values["EXTRA"] = "leaked"
	if _, ok := snap.Values["EXTRA"]; ok {
		t.Error("snapshot values should be a copy, not a reference")
	}
}

func TestSnapshotSaveAndLatest(t *testing.T) {
	store := newSnapshotStore(t)

	snap1 := &Snapshot{Environment: "staging", CreatedAt: time.Now().UTC().Add(-time.Second), Values: map[string]string{"A": "1"}}
	snap2 := &Snapshot{Environment: "staging", CreatedAt: time.Now().UTC(), Values: map[string]string{"A": "2", "B": "3"}}

	if err := store.Save(snap1); err != nil {
		t.Fatalf("Save snap1: %v", err)
	}
	if err := store.Save(snap2); err != nil {
		t.Fatalf("Save snap2: %v", err)
	}

	latest, err := store.Latest("staging")
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}
	if latest.Values["A"] != "2" {
		t.Errorf("expected latest A=2, got %q", latest.Values["A"])
	}
	if len(latest.Values) != 2 {
		t.Errorf("expected 2 keys in latest snapshot, got %d", len(latest.Values))
	}
}

func TestSnapshotLatestNotFound(t *testing.T) {
	store := newSnapshotStore(t)
	_, err := store.Latest("nonexistent")
	if err == nil {
		t.Error("expected error for missing environment, got nil")
	}
}

func TestSnapshotListEnvironments(t *testing.T) {
	store := newSnapshotStore(t)

	for _, env := range []string{"dev", "staging", "prod"} {
		snap := NewSnapshot(env, map[string]string{"X": "1"})
		if err := store.Save(snap); err != nil {
			t.Fatalf("Save %s: %v", env, err)
		}
	}

	envs, err := store.ListEnvironments()
	if err != nil {
		t.Fatalf("ListEnvironments: %v", err)
	}
	if len(envs) != 3 {
		t.Errorf("expected 3 environments, got %d: %v", len(envs), envs)
	}
}

func TestSnapshotDiffFrom(t *testing.T) {
	old := NewSnapshot("dev", map[string]string{"A": "1", "B": "2"})
	new := NewSnapshot("dev", map[string]string{"A": "updated", "C": "3"})

	entries := old.DiffFrom(new)
	if len(entries) != 3 {
		t.Errorf("expected 3 diff entries, got %d", len(entries))
	}
}
