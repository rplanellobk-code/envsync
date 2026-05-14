package env

import (
	"context"
	"fmt"
	"strings"

	"envsync/internal/storage"
)

const templatePrefix = "templates/"

// TemplateStore persists and retrieves env templates via a storage backend.
type TemplateStore struct {
	backend storage.Backend
}

// NewTemplateStore creates a TemplateStore backed by the given storage.Backend.
func NewTemplateStore(b storage.Backend) *TemplateStore {
	return &TemplateStore{backend: b}
}

// Save serialises and stores a template under the given name.
func (ts *TemplateStore) Save(ctx context.Context, name string, tmpl map[string]string) error {
	if name == "" {
		return fmt.Errorf("template name must not be empty")
	}
	var sb strings.Builder
	for _, k := range sortedKeys(tmpl) {
		if tmpl[k] == "" {
			fmt.Fprintf(&sb, "%s\n", k)
		} else {
			fmt.Fprintf(&sb, "%s=%s\n", k, tmpl[k])
		}
	}
	return ts.backend.Put(ctx, templatePrefix+name, []byte(sb.String()))
}

// Load retrieves and parses a stored template by name.
func (ts *TemplateStore) Load(ctx context.Context, name string) (map[string]string, error) {
	data, err := ts.backend.Get(ctx, templatePrefix+name)
	if err != nil {
		return nil, err
	}
	return ParseTemplate(string(data))
}

// List returns all stored template names.
func (ts *TemplateStore) List(ctx context.Context) ([]string, error) {
	keys, err := ts.backend.List(ctx, templatePrefix)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(keys))
	for _, k := range keys {
		names = append(names, strings.TrimPrefix(k, templatePrefix))
	}
	return names, nil
}

// Delete removes a stored template by name.
func (ts *TemplateStore) Delete(ctx context.Context, name string) error {
	return ts.backend.Delete(ctx, templatePrefix+name)
}
