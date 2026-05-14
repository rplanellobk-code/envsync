package env

import (
	"testing"
	"time"

	"envsync/internal/storage"
)

func newLockStore() *LockStore {
	return NewLockStore(storage.NewMemoryBackend())
}

func TestLockAcquire(t *testing.T) {
	s := newLockStore()
	lock, err := s.Acquire("production", "alice", time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lock.Owner != "alice" {
		t.Errorf("expected owner alice, got %q", lock.Owner)
	}
	if lock.Environment != "production" {
		t.Errorf("expected env production, got %q", lock.Environment)
	}
	if lock.IsExpired() {
		t.Error("lock should not be expired immediately after acquire")
	}
}

func TestLockAcquireConflict(t *testing.T) {
	s := newLockStore()
	if _, err := s.Acquire("staging", "alice", time.Minute); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	_, err := s.Acquire("staging", "bob", time.Minute)
	if err == nil {
		t.Fatal("expected error acquiring lock held by alice")
	}
}

func TestLockAcquireExpired(t *testing.T) {
	s := newLockStore()
	// Acquire with a TTL already in the past.
	if _, err := s.Acquire("dev", "alice", -time.Second); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	// Second acquire should succeed because the first lock is expired.
	lock, err := s.Acquire("dev", "bob", time.Minute)
	if err != nil {
		t.Fatalf("expected to acquire expired lock, got: %v", err)
	}
	if lock.Owner != "bob" {
		t.Errorf("expected owner bob, got %q", lock.Owner)
	}
}

func TestLockRelease(t *testing.T) {
	s := newLockStore()
	if _, err := s.Acquire("production", "alice", time.Minute); err != nil {
		t.Fatalf("acquire failed: %v", err)
	}
	if err := s.Release("production", "alice"); err != nil {
		t.Fatalf("release failed: %v", err)
	}
	// After release, another owner should be able to acquire.
	if _, err := s.Acquire("production", "bob", time.Minute); err != nil {
		t.Fatalf("acquire after release failed: %v", err)
	}
}

func TestLockReleaseWrongOwner(t *testing.T) {
	s := newLockStore()
	if _, err := s.Acquire("production", "alice", time.Minute); err != nil {
		t.Fatalf("acquire failed: %v", err)
	}
	if err := s.Release("production", "bob"); err == nil {
		t.Fatal("expected error releasing lock owned by alice as bob")
	}
}

func TestLockCurrent(t *testing.T) {
	s := newLockStore()
	_, err := s.Current("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing lock")
	}
	if _, err := s.Acquire("staging", "carol", time.Minute); err != nil {
		t.Fatalf("acquire failed: %v", err)
	}
	lock, err := s.Current("staging")
	if err != nil {
		t.Fatalf("current failed: %v", err)
	}
	if lock.Owner != "carol" {
		t.Errorf("expected owner carol, got %q", lock.Owner)
	}
}
