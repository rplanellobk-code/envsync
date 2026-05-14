package storage

import (
	"errors"
	"testing"
)

func TestMemoryPutGet(t *testing.T) {
	m := NewMemoryBackend()
	if err := m.Put("k", []byte("hello")); err != nil {
		t.Fatalf("put: %v", err)
	}
	got, err := m.Get("k")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if string(got) != "hello" {
		t.Errorf("want hello got %q", got)
	}
}

func TestMemoryGetNotFound(t *testing.T) {
	m := NewMemoryBackend()
	_, err := m.Get("missing")
	if !IsNotFound(err) {
		t.Fatalf("want not-found error, got %v", err)
	}
}

func TestMemoryDelete(t *testing.T) {
	m := NewMemoryBackend()
	_ = m.Put("x", []byte("v"))
	if err := m.Delete("x"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := m.Get("x"); !IsNotFound(err) {
		t.Fatal("expected not-found after delete")
	}
}

func TestMemoryDeleteNotFound(t *testing.T) {
	m := NewMemoryBackend()
	if err := m.Delete("nope"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

func TestMemoryList(t *testing.T) {
	m := NewMemoryBackend()
	for _, k := range []string{"a/1", "a/2", "b/1"} {
		_ = m.Put(k, []byte("v"))
	}
	keys, err := m.List("a/")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(keys) != 2 {
		t.Fatalf("want 2 keys got %d", len(keys))
	}
	if keys[0] != "a/1" || keys[1] != "a/2" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestMemoryIsolatesMutations(t *testing.T) {
	m := NewMemoryBackend()
	orig := []byte("original")
	_ = m.Put("k", orig)
	got, _ := m.Get("k")
	got[0] = 'X'
	again, _ := m.Get("k")
	if again[0] == 'X' {
		t.Error("backend should return a copy, not a reference")
	}
}
