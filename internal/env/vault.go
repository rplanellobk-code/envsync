package env

import (
	"fmt"

	"github.com/user/envsync/internal/crypto"
	"github.com/user/envsync/internal/storage"
)

// Vault provides encrypted storage and retrieval of env files using a
// storage backend and a passphrase-derived key.
type Vault struct {
	backend    storage.Backend
	passphrase string
}

// NewVault creates a Vault backed by the given storage.Backend.
func NewVault(backend storage.Backend, passphrase string) (*Vault, error) {
	if err := crypto.ValidatePassphrase(passphrase); err != nil {
		return nil, err
	}
	return &Vault{backend: backend, passphrase: passphrase}, nil
}

// Push serializes, encrypts, and stores an env map under the given name.
func (v *Vault) Push(name string, env map[string]string) error {
	plain := Serialize(env)

	key, err := crypto.DeriveKey(v.passphrase, nil)
	if err != nil {
		return fmt.Errorf("vault push: derive key: %w", err)
	}

	ciphertext, err := crypto.Encrypt(key, []byte(plain))
	if err != nil {
		return fmt.Errorf("vault push: encrypt: %w", err)
	}

	if err := v.backend.Put(name, ciphertext); err != nil {
		return fmt.Errorf("vault push: store: %w", err)
	}
	return nil
}

// Pull retrieves, decrypts, and parses an env map stored under the given name.
func (v *Vault) Pull(name string) (map[string]string, error) {
	ciphertext, err := v.backend.Get(name)
	if err != nil {
		return nil, fmt.Errorf("vault pull: fetch: %w", err)
	}

	key, err := crypto.DeriveKey(v.passphrase, nil)
	if err != nil {
		return nil, fmt.Errorf("vault pull: derive key: %w", err)
	}

	plain, err := crypto.Decrypt(key, ciphertext)
	if err != nil {
		return nil, fmt.Errorf("vault pull: decrypt: %w", err)
	}

	env, err := Parse(string(plain))
	if err != nil {
		return nil, fmt.Errorf("vault pull: parse: %w", err)
	}
	return env, nil
}

// List returns the names of all env files stored in the backend.
func (v *Vault) List() ([]string, error) {
	names, err := v.backend.List()
	if err != nil {
		return nil, fmt.Errorf("vault list: %w", err)
	}
	return names, nil
}
