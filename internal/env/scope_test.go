package env

import (
	"context"
	"testing"

	"envsync/internal/storage"
)

func newScopeStore() *ScopeStore {
	return NewScopeStore(storage.NewMemoryBackend())
}

func TestScopeMatchPrefix(t *testing.T) {
	s := &Scope{Name: "db", Prefix: "DB_"}
	if !s.Match("DB_HOST") {
		t.Error("expected DB_HOST to match")
	}
	if s.Match("APP_KEY") {
		t.Error("expected APP_KEY not to match")
	}
}

func TestScopeMatchPattern(t *testing.T) {
	s := &Scope{Name: "secrets", Pattern: `^SECRET_`}
	_ = s.compile()
	if !s.Match("SECRET_KEY") {
		t.Error("expected SECRET_KEY to match")
	}
	if s.Match("DB_PASSWORD") {
		t.Error("expected DB_PASSWORD not to match")
	}
}

func TestScopeFilter(t *testing.T) {
	s := &Scope{Name: "app", Prefix: "APP_"}
	env := map[string]string{
		"APP_NAME": "myapp",
		"APP_PORT": "8080",
		"DB_HOST":  "localhost",
	}
	got := s.Filter(env)
	if len(got) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(got))
	}
	if got["APP_NAME"] != "myapp" || got["APP_PORT"] != "8080" {
		t.Error("unexpected filter result")
	}
}

func TestScopeSaveAndLoad(t *testing.T) {
	ctx := context.Background()
	store := newScopeStore()
	s := &Scope{Name: "db", Prefix: "DB_", Pattern: `^DB_`}
	if err := store.Save(ctx, s); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := store.Load(ctx, "db")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Name != "db" || loaded.Prefix != "DB_" || loaded.Pattern != `^DB_` {
		t.Errorf("unexpected loaded scope: %+v", loaded)
	}
}

func TestScopeLoadNotFound(t *testing.T) {
	ctx := context.Background()
	store := newScopeStore()
	_, err := store.Load(ctx, "missing")
	if err == nil {
		t.Fatal("expected error for missing scope")
	}
}

func TestScopeList(t *testing.T) {
	ctx := context.Background()
	store := newScopeStore()
	for _, name := range []string{"alpha", "beta", "gamma"} {
		_ = store.Save(ctx, &Scope{Name: name, Prefix: name + "_"})
	}
	names, err := store.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 3 {
		t.Fatalf("expected 3 scopes, got %d", len(names))
	}
}

func TestScopeDelete(t *testing.T) {
	ctx := context.Background()
	store := newScopeStore()
	_ = store.Save(ctx, &Scope{Name: "temp", Prefix: "TEMP_"})
	if err := store.Delete(ctx, "temp"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := store.Load(ctx, "temp")
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestScopeEmptyNameError(t *testing.T) {
	ctx := context.Background()
	store := newScopeStore()
	err := store.Save(ctx, &Scope{Name: "", Prefix: "X_"})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}
