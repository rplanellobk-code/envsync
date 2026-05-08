package storage_test

import (
	"context"
	"testing"

	"envsync/internal/storage"
)

func newTestBackend(t *testing.T) *storage.FileBackend {
	t.Helper()
	b, err := storage.NewFileBackend(t.TempDir())
	if err != nil {
		t.Fatalf("NewFileBackend: %v", err)
	}
	return b
}

func TestFilePutGet(t *testing.T) {
	ctx := context.Background()
	b := newTestBackend(t)

	data := []byte("encrypted-payload")
	if err := b.Put(ctx, "production", data); err != nil {
		t.Fatalf("Put: %v", err)
	}

	got, err := b.Get(ctx, "production")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if string(got) != string(data) {
		t.Errorf("got %q, want %q", got, data)
	}
}

func TestFileGetNotFound(t *testing.T) {
	ctx := context.Background()
	b := newTestBackend(t)

	_, err := b.Get(ctx, "missing")
	if !storage.IsNotFound(err) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestFileList(t *testing.T) {
	ctx := context.Background()
	b := newTestBackend(t)

	for _, k := range []string{"dev", "staging", "prod"} {
		if err := b.Put(ctx, k, []byte(k)); err != nil {
			t.Fatalf("Put %s: %v", k, err)
		}
	}

	keys, err := b.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d: %v", len(keys), keys)
	}
}

func TestFileDelete(t *testing.T) {
	ctx := context.Background()
	b := newTestBackend(t)

	if err := b.Put(ctx, "tmp", []byte("x")); err != nil {
		t.Fatalf("Put: %v", err)
	}
	if err := b.Delete(ctx, "tmp"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := b.Get(ctx, "tmp"); !storage.IsNotFound(err) {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestFileDeleteNotFound(t *testing.T) {
	ctx := context.Background()
	b := newTestBackend(t)

	if err := b.Delete(ctx, "ghost"); !storage.IsNotFound(err) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
