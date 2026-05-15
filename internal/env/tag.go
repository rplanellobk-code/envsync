package env

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"envsync/internal/storage"
)

var validTagName = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)

// Tag associates a human-readable label with a specific snapshot index
// for a given environment.
type Tag struct {
	Name        string    `json:"name"`
	Environment string    `json:"environment"`
	SnapshotIdx int       `json:"snapshot_idx"`
	CreatedAt   time.Time `json:"created_at"`
	Message     string    `json:"message,omitempty"`
}

// TagStore persists and retrieves tags via a storage backend.
type TagStore struct {
	backend storage.Backend
}

// NewTagStore creates a TagStore backed by the given storage backend.
func NewTagStore(b storage.Backend) *TagStore {
	return &TagStore{backend: b}
}

func tagKey(env, name string) string {
	return fmt.Sprintf("tags/%s/%s", env, name)
}

// Save persists a tag. Returns an error if the tag name is invalid.
func (ts *TagStore) Save(t Tag) error {
	if !validTagName.MatchString(t.Name) {
		return fmt.Errorf("tag: invalid name %q (alphanumeric, -, _, . only)", t.Name)
	}
	if t.Environment == "" {
		return fmt.Errorf("tag: environment must not be empty")
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now().UTC()
	}
	data, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("tag: marshal: %w", err)
	}
	return ts.backend.Put(tagKey(t.Environment, t.Name), data)
}

// Get retrieves a tag by environment and name.
func (ts *TagStore) Get(env, name string) (Tag, error) {
	data, err := ts.backend.Get(tagKey(env, name))
	if err != nil {
		return Tag{}, fmt.Errorf("tag: get %q/%q: %w", env, name, err)
	}
	var t Tag
	if err := json.Unmarshal(data, &t); err != nil {
		return Tag{}, fmt.Errorf("tag: unmarshal: %w", err)
	}
	return t, nil
}

// Delete removes a tag by environment and name.
func (ts *TagStore) Delete(env, name string) error {
	return ts.backend.Delete(tagKey(env, name))
}

// List returns all tags for the given environment.
func (ts *TagStore) List(env string) ([]Tag, error) {
	prefix := fmt.Sprintf("tags/%s/", env)
	keys, err := ts.backend.List(prefix)
	if err != nil {
		return nil, fmt.Errorf("tag: list: %w", err)
	}
	tags := make([]Tag, 0, len(keys))
	for _, k := range keys {
		data, err := ts.backend.Get(k)
		if err != nil {
			return nil, fmt.Errorf("tag: list get %q: %w", k, err)
		}
		var t Tag
		if err := json.Unmarshal(data, &t); err != nil {
			return nil, fmt.Errorf("tag: list unmarshal: %w", err)
		}
		tags = append(tags, t)
	}
	return tags, nil
}
