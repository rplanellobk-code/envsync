package env

import (
	"context"
	"testing"
	"time"

	"envsync/internal/storage"
)

func newWatchFixture(t *testing.T) (*Vault, string) {
	t.Helper()
	back := storage.NewMemoryBackend()
	v := NewVault(back)
	return v, "watch-passphrase"
}

func TestWatchNoChange(t *testing.T) {
	v, pass := newWatchFixture(t)
	env := "staging"

	if err := v.Push(env, map[string]string{"K": "1"}, pass); err != nil {
		t.Fatalf("push: %v", err)
	}

	w := NewWatcher(v, pass, 20*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	ch, err := w.Watch(ctx, env)
	if err != nil {
		t.Fatalf("Watch: %v", err)
	}

	var events []WatchEvent
	for e := range ch {
		events = append(events, e)
	}
	if len(events) != 0 {
		t.Errorf("expected 0 events, got %d", len(events))
	}
}

func TestWatchDetectsChange(t *testing.T) {
	v, pass := newWatchFixture(t)
	env := "production"

	if err := v.Push(env, map[string]string{"K": "v1"}, pass); err != nil {
		t.Fatalf("push: %v", err)
	}

	w := NewWatcher(v, pass, 20*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	ch, err := w.Watch(ctx, env)
	if err != nil {
		t.Fatalf("Watch: %v", err)
	}

	// Mutate the remote after a short delay.
	go func() {
		time.Sleep(40 * time.Millisecond)
		_ = v.Push(env, map[string]string{"K": "v2"}, pass)
	}()

	select {
	case ev := <-ch:
		if ev.Environment != env {
			t.Errorf("expected env %q, got %q", env, ev.Environment)
		}
		if ev.PreviousSum == ev.CurrentSum {
			t.Error("expected sums to differ")
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for change event")
	}
}

func TestWatchInitialPullError(t *testing.T) {
	v, pass := newWatchFixture(t)
	// No data pushed — Pull will return NotFound.
	w := NewWatcher(v, pass, 20*time.Millisecond)
	_, err := w.Watch(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error for missing environment, got nil")
	}
}

// TestWatchDetectsMultipleChanges verifies that the watcher emits an event for
// each successive mutation while the context remains active.
func TestWatchDetectsMultipleChanges(t *testing.T) {
	v, pass := newWatchFixture(t)
	env := "multi"

	if err := v.Push(env, map[string]string{"K": "v0"}, pass); err != nil {
		t.Fatalf("push: %v", err)
	}

	w := NewWatcher(v, pass, 20*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	ch, err := w.Watch(ctx, env)
	if err != nil {
		t.Fatalf("Watch: %v", err)
	}

	// Push two successive changes with enough spacing for the poller to catch each.
	go func() {
		time.Sleep(40 * time.Millisecond)
		_ = v.Push(env, map[string]string{"K": "v1"}, pass)
		time.Sleep(60 * time.Millisecond)
		_ = v.Push(env, map[string]string{"K": "v2"}, pass)
	}()

	var received int
	for received < 2 {
		select {
		case <-ch:
			received++
		case <-ctx.Done():
			t.Fatalf("timed out after %d/2 change events", received)
		}
	}
}
