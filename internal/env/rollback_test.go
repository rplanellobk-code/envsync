package env

import (
	"testing"
)

func TestRollbackToSpecificSnapshot(t *testing.T) {
	v, store, backend := newRollbackFixture(t)
	_ = backend

	data1 := map[string]string{"KEY": "v1"}
	data2 := map[string]string{"KEY": "v2"}

	if err := v.Push("prod", data1, "alice"); err != nil {
		t.Fatalf("push v1: %v", err)
	}
	if err := store.Save(NewSnapshot("prod", data1)); err != nil {
		t.Fatalf("save snap1: %v", err)
	}

	if err := v.Push("prod", data2, "alice"); err != nil {
		t.Fatalf("push v2: %v", err)
	}
	snap2 := NewSnapshot("prod", data2)
	if err := store.Save(snap2); err != nil {
		t.Fatalf("save snap2: %v", err)
	}

	snaps, _ := store.List("prod")
	if err := Rollback(v, store, "prod", snaps[1].ID, "alice"); err != nil {
		t.Fatalf("rollback: %v", err)
	}

	got, err := v.Pull("prod")
	if err != nil {
		t.Fatalf("pull after rollback: %v", err)
	}
	if got["KEY"] != "v1" {
		t.Errorf("expected KEY=v1, got %q", got["KEY"])
	}
}

func TestRollbackSnapshotNotFound(t *testing.T) {
	v, store, _ := newRollbackFixture(t)

	if err := v.Push("prod", map[string]string{"K": "1"}, "alice"); err != nil {
		t.Fatalf("push: %v", err)
	}
	if err := store.Save(NewSnapshot("prod", map[string]string{"K": "1"})); err != nil {
		t.Fatalf("save: %v", err)
	}

	err := Rollback(v, store, "prod", "nonexistent-id", "alice")
	if err == nil {
		t.Fatal("expected error for missing snapshot ID")
	}
}

func TestRollbackToLatest(t *testing.T) {
	v, store, _ := newRollbackFixture(t)

	for i, kv := range []map[string]string{{"N": "1"}, {"N": "2"}, {"N": "3"}} {
		if err := v.Push("staging", kv, "bob"); err != nil {
			t.Fatalf("push %d: %v", i, err)
		}
		if err := store.Save(NewSnapshot("staging", kv)); err != nil {
			t.Fatalf("save snap %d: %v", i, err)
		}
	}

	if err := RollbackToLatest(v, store, "staging", "bob"); err != nil {
		t.Fatalf("rollback to latest: %v", err)
	}

	got, err := v.Pull("staging")
	if err != nil {
		t.Fatalf("pull: %v", err)
	}
	if got["N"] != "2" {
		t.Errorf("expected N=2, got %q", got["N"])
	}
}

func TestRollbackToLatestInsufficientSnapshots(t *testing.T) {
	v, store, _ := newRollbackFixture(t)

	if err := v.Push("dev", map[string]string{"X": "1"}, "carol"); err != nil {
		t.Fatalf("push: %v", err)
	}
	if err := store.Save(NewSnapshot("dev", map[string]string{"X": "1"})); err != nil {
		t.Fatalf("save: %v", err)
	}

	if err := RollbackToLatest(v, store, "dev", "carol"); err == nil {
		t.Fatal("expected error with only one snapshot")
	}
}

// newRollbackFixture wires up a VaultWithAudit and SnapshotStore backed by
// in-memory file backends suitable for rollback tests.
func newRollbackFixture(t *testing.T) (*VaultWithAudit, *SnapshotStore, interface{}) {
	t.Helper()
	vault := newTestVault(t)
	snapshotBackend := newTestBackend(t)
	store := NewSnapshotStore(snapshotBackend)
	return vault, store, snapshotBackend
}
