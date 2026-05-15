package env

import (
	"encoding/json"
	"fmt"
	"time"

	"envsync/internal/storage"
)

// Pin records a named, immutable reference to a specific snapshot version
// for a given environment, similar to a git tag pointing at a commit.
type Pin struct {
	Name        string    `json:"name"`
	Environment string    `json:"environment"`
	SnapshotID  string    `json:"snapshot_id"`
	CreatedAt   time.Time `json:"created_at"`
	Note        string    `json:"note,omitempty"`
}

// PinStore persists and retrieves pins using a storage backend.
type PinStore struct {
	backend storage.Backend
}

// NewPinStore returns a PinStore backed by the given storage.Backend.
func NewPinStore(b storage.Backend) *PinStore {
	return &PinStore{backend: b}
}

func pinKey(env, name string) string {
	return fmt.Sprintf("pins/%s/%s", env, name)
}

// Save persists a pin. Returns an error if a pin with the same name already exists.
func (ps *PinStore) Save(p Pin) error {
	key := pinKey(p.Environment, p.Name)
	_, err := ps.backend.Get(key)
	if err == nil {
		return fmt.Errorf("pin %q already exists for environment %q", p.Name, p.Environment)
	}
	if !storage.IsNotFound(err) {
		return err
	}
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return ps.backend.Put(key, data)
}

// Get retrieves a pin by environment and name.
func (ps *PinStore) Get(env, name string) (Pin, error) {
	data, err := ps.backend.Get(pinKey(env, name))
	if err != nil {
		return Pin{}, err
	}
	var p Pin
	if err := json.Unmarshal(data, &p); err != nil {
		return Pin{}, err
	}
	return p, nil
}

// Delete removes a pin by environment and name.
func (ps *PinStore) Delete(env, name string) error {
	return ps.backend.Delete(pinKey(env, name))
}

// List returns all pins for the given environment.
func (ps *PinStore) List(env string) ([]Pin, error) {
	prefix := fmt.Sprintf("pins/%s/", env)
	keys, err := ps.backend.List(prefix)
	if err != nil {
		return nil, err
	}
	pins := make([]Pin, 0, len(keys))
	for _, k := range keys {
		data, err := ps.backend.Get(k)
		if err != nil {
			return nil, err
		}
		var p Pin
		if err := json.Unmarshal(data, &p); err != nil {
			return nil, err
		}
		pins = append(pins, p)
	}
	return pins, nil
}
