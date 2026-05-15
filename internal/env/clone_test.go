package env

import (
	"testing"

	"github.com/your/envsync/internal/storage"
)

func newCloneFixture(t *testing.T) (*Cloner, *Vault) {
	t.Helper()
	b := storage.NewMemoryBackend()
	v := NewVault(b)
	return NewCloner(v), v
}

func TestCloneBasic(t *testing.T) {
	cloner, vault := newCloneFixture(t)
	pass := "passphrase"

	src := map[string]string{"A": "1", "B": "2"}
	if err := vault.Push(pass, "prod", src); err != nil {
		t.Fatalf("push: %v", err)
	}

	if err := cloner.Clone(pass, "prod", "staging", CloneOptions{}); err != nil {
		t.Fatalf("clone: %v", err)
	}

	got, err := vault.Pull(pass, "staging")
	if err != nil {
		t.Fatalf("pull: %v", err)
	}
	if got["A"] != "1" || got["B"] != "2" {
		t.Errorf("unexpected values: %v", got)
	}
}

func TestCloneOmitKeys(t *testing.T) {
	cloner, vault := newCloneFixture(t)
	pass := "passphrase"

	vault.Push(pass, "prod", map[string]string{"A": "1", "SECRET": "s3cr3t"})

	cloner.Clone(pass, "prod", "staging", CloneOptions{OmitKeys: []string{"SECRET"}})

	got, _ := vault.Pull(pass, "staging")
	if _, ok := got["SECRET"]; ok {
		t.Error("SECRET should have been omitted")
	}
	if got["A"] != "1" {
		t.Errorf("expected A=1, got %v", got["A"])
	}
}

func TestCloneOverrideKeys(t *testing.T) {
	cloner, vault := newCloneFixture(t)
	pass := "passphrase"

	vault.Push(pass, "prod", map[string]string{"DB_HOST": "prod-db", "PORT": "5432"})

	cloner.Clone(pass, "prod", "staging", CloneOptions{
		OverrideKeys: map[string]string{"DB_HOST": "staging-db"},
	})

	got, _ := vault.Pull(pass, "staging")
	if got["DB_HOST"] != "staging-db" {
		t.Errorf("expected staging-db, got %v", got["DB_HOST"])
	}
	if got["PORT"] != "5432" {
		t.Errorf("expected 5432, got %v", got["PORT"])
	}
}

func TestCloneNoOverwrite(t *testing.T) {
	cloner, vault := newCloneFixture(t)
	pass := "passphrase"

	vault.Push(pass, "prod", map[string]string{"A": "new"})
	vault.Push(pass, "staging", map[string]string{"A": "existing"})

	cloner.Clone(pass, "prod", "staging", CloneOptions{Overwrite: false})

	got, _ := vault.Pull(pass, "staging")
	if got["A"] != "existing" {
		t.Errorf("expected existing value preserved, got %v", got["A"])
	}
}

func TestCloneOverwrite(t *testing.T) {
	cloner, vault := newCloneFixture(t)
	pass := "passphrase"

	vault.Push(pass, "prod", map[string]string{"A": "new"})
	vault.Push(pass, "staging", map[string]string{"A": "existing"})

	cloner.Clone(pass, "prod", "staging", CloneOptions{Overwrite: true})

	got, _ := vault.Pull(pass, "staging")
	if got["A"] != "new" {
		t.Errorf("expected overwritten value, got %v", got["A"])
	}
}

func TestCloneEmptySource(t *testing.T) {
	cloner, _ := newCloneFixture(t)
	err := cloner.Clone("pass", "", "dst", CloneOptions{})
	if err == nil {
		t.Error("expected error for empty source")
	}
}

func TestCloneEmptyDestination(t *testing.T) {
	cloner, _ := newCloneFixture(t)
	err := cloner.Clone("pass", "src", "", CloneOptions{})
	if err == nil {
		t.Error("expected error for empty destination")
	}
}
