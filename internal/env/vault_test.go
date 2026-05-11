package env_test

import (
	"os"
	"testing"

	"github.com/user/envsync/internal/env"
	"github.com/user/envsync/internal/storage"
)

func newTestVault(t *testing.T) *env.Vault {
	t.Helper()
	dir, err := os.MkdirTemp("", "vault-test-*")
	if err != nil {
		t.Fatalf("create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })

	backend, err := storage.NewFileBackend(dir)
	if err != nil {
		t.Fatalf("new file backend: %v", err)
	}

	v, err := env.NewVault(backend, "correct-horse-battery-staple")
	if err != nil {
		t.Fatalf("new vault: %v", err)
	}
	return v
}

func TestVaultPushPull(t *testing.T) {
	v := newTestVault(t)

	orig := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := v.Push("prod", orig); err != nil {
		t.Fatalf("push: %v", err)
	}

	got, err := v.Pull("prod")
	if err != nil {
		t.Fatalf("pull: %v", err)
	}

	for k, want := range orig {
		if got[k] != want {
			t.Errorf("key %q: got %q, want %q", k, got[k], want)
		}
	}
}

func TestVaultPullNotFound(t *testing.T) {
	v := newTestVault(t)

	_, err := v.Pull("missing")
	if err == nil {
		t.Fatal("expected error pulling missing env, got nil")
	}
}

func TestVaultList(t *testing.T) {
	v := newTestVault(t)

	envs := []string{"prod", "staging", "dev"}
	for _, name := range envs {
		if err := v.Push(name, map[string]string{"KEY": name}); err != nil {
			t.Fatalf("push %s: %v", name, err)
		}
	}

	list, err := v.List()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != len(envs) {
		t.Errorf("list length: got %d, want %d", len(list), len(envs))
	}
}

func TestVaultWrongPassphrase(t *testing.T) {
	dir, _ := os.MkdirTemp("", "vault-wp-*")
	t.Cleanup(func() { os.RemoveAll(dir) })

	backend, _ := storage.NewFileBackend(dir)

	v1, _ := env.NewVault(backend, "correct-horse-battery-staple")
	v1.Push("prod", map[string]string{"SECRET": "value"})

	v2, _ := env.NewVault(backend, "wrong-passphrase-here-xx")
	_, err := v2.Pull("prod")
	if err == nil {
		t.Fatal("expected error with wrong passphrase, got nil")
	}
}
