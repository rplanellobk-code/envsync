package env

import (
	"testing"

	"envsync/internal/storage"
)

func newTagStore() *TagStore {
	return NewTagStore(storage.NewMemoryBackend())
}

func TestTagSaveAndGet(t *testing.T) {
	ts := newTagStore()
	tag := Tag{
		Name:        "v1.0",
		Environment: "production",
		SnapshotIdx: 3,
		Message:     "initial release",
	}
	if err := ts.Save(tag); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := ts.Get("production", "v1.0")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.SnapshotIdx != 3 {
		t.Errorf("SnapshotIdx: want 3, got %d", got.SnapshotIdx)
	}
	if got.Message != "initial release" {
		t.Errorf("Message: want 'initial release', got %q", got.Message)
	}
}

func TestTagGetNotFound(t *testing.T) {
	ts := newTagStore()
	_, err := ts.Get("staging", "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing tag, got nil")
	}
}

func TestTagInvalidName(t *testing.T) {
	ts := newTagStore()
	tag := Tag{Name: "bad name!", Environment: "dev", SnapshotIdx: 0}
	if err := ts.Save(tag); err == nil {
		t.Fatal("expected error for invalid tag name")
	}
}

func TestTagEmptyEnvironment(t *testing.T) {
	ts := newTagStore()
	tag := Tag{Name: "v1", Environment: "", SnapshotIdx: 1}
	if err := ts.Save(tag); err == nil {
		t.Fatal("expected error for empty environment")
	}
}

func TestTagDelete(t *testing.T) {
	ts := newTagStore()
	tag := Tag{Name: "v2.0", Environment: "staging", SnapshotIdx: 5}
	if err := ts.Save(tag); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if err := ts.Delete("staging", "v2.0"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := ts.Get("staging", "v2.0"); err == nil {
		t.Fatal("expected error after deletion")
	}
}

func TestTagList(t *testing.T) {
	ts := newTagStore()
	for _, name := range []string{"v1.0", "v1.1", "v2.0"} {
		if err := ts.Save(Tag{Name: name, Environment: "prod", SnapshotIdx: 1}); err != nil {
			t.Fatalf("Save %q: %v", name, err)
		}
	}
	// save a tag for a different env — should not appear
	if err := ts.Save(Tag{Name: "v1.0", Environment: "dev", SnapshotIdx: 0}); err != nil {
		t.Fatalf("Save dev tag: %v", err)
	}
	tags, err := ts.List("prod")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(tags) != 3 {
		t.Errorf("want 3 tags, got %d", len(tags))
	}
}
