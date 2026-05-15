package env

import (
	"testing"
	"time"

	"envsync/internal/storage"
)

func newPinStore() *PinStore {
	return NewPinStore(storage.NewMemoryBackend())
}

func TestPinSaveAndGet(t *testing.T) {
	ps := newPinStore()
	p := Pin{
		Name:        "v1.0",
		Environment: "production",
		SnapshotID:  "snap-abc123",
		CreatedAt:   time.Now().UTC().Truncate(time.Second),
		Note:        "initial release",
	}
	if err := ps.Save(p); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := ps.Get("production", "v1.0")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.SnapshotID != p.SnapshotID {
		t.Errorf("SnapshotID: got %q, want %q", got.SnapshotID, p.SnapshotID)
	}
	if got.Note != p.Note {
		t.Errorf("Note: got %q, want %q", got.Note, p.Note)
	}
}

func TestPinDuplicateReturnsError(t *testing.T) {
	ps := newPinStore()
	p := Pin{Name: "v1.0", Environment: "staging", SnapshotID: "snap-001", CreatedAt: time.Now()}
	if err := ps.Save(p); err != nil {
		t.Fatalf("first Save: %v", err)
	}
	if err := ps.Save(p); err == nil {
		t.Error("expected error on duplicate Save, got nil")
	}
}

func TestPinGetNotFound(t *testing.T) {
	ps := newPinStore()
	_, err := ps.Get("production", "missing")
	if !storage.IsNotFound(err) {
		t.Errorf("expected not-found error, got %v", err)
	}
}

func TestPinDelete(t *testing.T) {
	ps := newPinStore()
	p := Pin{Name: "v2.0", Environment: "dev", SnapshotID: "snap-xyz", CreatedAt: time.Now()}
	_ = ps.Save(p)
	if err := ps.Delete("dev", "v2.0"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := ps.Get("dev", "v2.0")
	if !storage.IsNotFound(err) {
		t.Errorf("expected not-found after delete, got %v", err)
	}
}

func TestPinList(t *testing.T) {
	ps := newPinStore()
	env := "production"
	names := []string{"v1.0", "v1.1", "v2.0"}
	for _, n := range names {
		_ = ps.Save(Pin{Name: n, Environment: env, SnapshotID: "snap-" + n, CreatedAt: time.Now()})
	}
	// add a pin for a different env — should not appear
	_ = ps.Save(Pin{Name: "v1.0", Environment: "staging", SnapshotID: "snap-s1", CreatedAt: time.Now()})

	pins, err := ps.List(env)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(pins) != len(names) {
		t.Errorf("List count: got %d, want %d", len(pins), len(names))
	}
}

func TestPinListEmpty(t *testing.T) {
	ps := newPinStore()
	pins, err := ps.List("nonexistent")
	if err != nil {
		t.Fatalf("List on empty: %v", err)
	}
	if len(pins) != 0 {
		t.Errorf("expected 0 pins, got %d", len(pins))
	}
}
