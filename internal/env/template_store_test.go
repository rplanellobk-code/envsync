package env

import (
	"context"
	"testing"

	"envsync/internal/storage"
)

func newTemplateStore(t *testing.T) *TemplateStore {
	t.Helper()
	dir := t.TempDir()
	b, err := storage.NewFileBackend(dir)
	if err != nil {
		t.Fatalf("NewFileBackend: %v", err)
	}
	return NewTemplateStore(b)
}

func TestTemplateSaveAndLoad(t *testing.T) {
	ts := newTemplateStore(t)
	ctx := context.Background()

	tmpl := map[string]string{"HOST": "", "PORT": "3000", "ENV": "production"}
	if err := ts.Save(ctx, "web", tmpl); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := ts.Load(ctx, "web")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded["HOST"] != "" {
		t.Errorf("expected empty default for HOST")
	}
	if loaded["PORT"] != "3000" {
		t.Errorf("expected 3000 for PORT, got %q", loaded["PORT"])
	}
	if loaded["ENV"] != "production" {
		t.Errorf("expected production for ENV")
	}
}

func TestTemplateLoadNotFound(t *testing.T) {
	ts := newTemplateStore(t)
	_, err := ts.Load(context.Background(), "ghost")
	if err == nil || !storage.IsNotFound(err) {
		t.Fatalf("expected not-found error, got %v", err)
	}
}

func TestTemplateList(t *testing.T) {
	ts := newTemplateStore(t)
	ctx := context.Background()

	for _, name := range []string{"alpha", "beta", "gamma"} {
		if err := ts.Save(ctx, name, map[string]string{"K": ""}); err != nil {
			t.Fatalf("Save %s: %v", name, err)
		}
	}

	names, err := ts.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 3 {
		t.Errorf("expected 3 templates, got %d", len(names))
	}
}

func TestTemplateDelete(t *testing.T) {
	ts := newTemplateStore(t)
	ctx := context.Background()

	if err := ts.Save(ctx, "tmp", map[string]string{"X": "1"}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if err := ts.Delete(ctx, "tmp"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := ts.Load(ctx, "tmp")
	if err == nil || !storage.IsNotFound(err) {
		t.Fatalf("expected not-found after delete, got %v", err)
	}
}

func TestTemplateSaveEmptyName(t *testing.T) {
	ts := newTemplateStore(t)
	err := ts.Save(context.Background(), "", map[string]string{"K": ""})
	if err == nil {
		t.Fatal("expected error for empty template name")
	}
}
