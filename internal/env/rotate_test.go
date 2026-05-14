package env

import (
	"testing"

	"github.com/user/envsync/internal/storage"
)

func newRotateFixture(t *testing.T) (*Rotator, *Vault, *AuditLog, *SnapshotStore) {
	t.Helper()
	b := storage.NewMemoryBackend()
	vault := NewVault(b)
	audit := NewAuditLog(b)
	snaps := NewSnapshotStore(b)
	return NewRotator(vault, audit, snaps), vault, audit, snaps
}

func TestRotateSuccess(t *testing.T) {
	rot, vault, _, snaps := newRotateFixture(t)
	pairs := map[string]string{"KEY": "value", "OTHER": "data"}
	if err := vault.Push("prod", pairs, "old-pass"); err != nil {
		t.Fatalf("push: %v", err)
	}
	if _, err := rot.Rotate("prod", "old-pass", "new-pass"); err != nil {
		t.Fatalf("rotate: %v", err)
	}
	got, err := vault.Pull("prod", "new-pass")
	if err != nil {
		t.Fatalf("pull after rotate: %v", err)
	}
	for k, v := range pairs {
		if got[k] != v {
			t.Errorf("key %q: want %q got %q", k, v, got[k])
		}
	}
	snap, err := snaps.Latest("prod")
	if err != nil {
		t.Fatalf("snapshot after rotate: %v", err)
	}
	if snap.Environment != "prod" {
		t.Errorf("snapshot env: want prod got %q", snap.Environment)
	}
}

func TestRotateSamePassphrase(t *testing.T) {
	rot, vault, _, _ := newRotateFixture(t)
	_ = vault.Push("dev", map[string]string{"A": "1"}, "pass")
	if _, err := rot.Rotate("dev", "pass", "pass"); err == nil {
		t.Fatal("expected error when new passphrase equals old")
	}
}

func TestRotateEmptyEnvironment(t *testing.T) {
	rot, _, _, _ := newRotateFixture(t)
	if _, err := rot.Rotate("", "old", "new"); err == nil {
		t.Fatal("expected error for empty environment")
	}
}

func TestRotateWrongOldPassphrase(t *testing.T) {
	rot, vault, _, _ := newRotateFixture(t)
	_ = vault.Push("staging", map[string]string{"X": "y"}, "correct")
	if _, err := rot.Rotate("staging", "wrong", "new-pass"); err == nil {
		t.Fatal("expected error for wrong old passphrase")
	}
}

func TestRotateAll(t *testing.T) {
	rot, vault, _, _ := newRotateFixture(t)
	for _, env := range []string{"dev", "staging", "prod"} {
		_ = vault.Push(env, map[string]string{"K": "v"}, "old")
	}
	results, err := rot.RotateAll("old", "new")
	if err != nil {
		t.Fatalf("rotate-all: %v", err)
	}
	if len(results) != 3 {
		t.Errorf("want 3 results got %d", len(results))
	}
}

func TestRotateAuditRecorded(t *testing.T) {
	rot, vault, audit, _ := newRotateFixture(t)
	_ = vault.Push("prod", map[string]string{"A": "b"}, "old")
	if _, err := rot.Rotate("prod", "old", "new"); err != nil {
		t.Fatalf("rotate: %v", err)
	}
	entries, err := audit.List("prod")
	if err != nil {
		t.Fatalf("audit list: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected at least one audit entry after rotate")
	}
	if entries[len(entries)-1].Action != "rotate" {
		t.Errorf("last action: want rotate got %q", entries[len(entries)-1].Action)
	}
}
