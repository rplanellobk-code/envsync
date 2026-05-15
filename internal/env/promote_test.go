package env

import (
	"testing"

	"github.com/user/envsync/internal/storage"
)

func newPromoterFixture() (*Promoter, *Vault) {
	backend := storage.NewMemoryBackend()
	v := NewVault(backend)
	return NewPromoter(v), v
}

func TestPromoteBasic(t *testing.T) {
	p, v := newPromoterFixture()
	const pass = "passphrase-32-bytes-long-exactly!"

	_ = v.Push("staging", pass, map[string]string{"FOO": "bar", "BAZ": "qux"})

	res, err := p.Promote("staging", "production", pass, PromoteOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Applied["FOO"] != "bar" || res.Applied["BAZ"] != "qux" {
		t.Errorf("unexpected applied vars: %v", res.Applied)
	}
	if len(res.Diffs) == 0 {
		t.Error("expected diffs to be non-empty on first promotion")
	}
}

func TestPromoteOmitKeys(t *testing.T) {
	p, v := newPromoterFixture()
	const pass = "passphrase-32-bytes-long-exactly!"

	_ = v.Push("staging", pass, map[string]string{"FOO": "bar", "SECRET": "hidden"})

	res, err := p.Promote("staging", "production", pass, PromoteOptions{OmitKeys: []string{"SECRET"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Applied["SECRET"]; ok {
		t.Error("SECRET should have been omitted")
	}
	if res.Applied["FOO"] != "bar" {
		t.Errorf("FOO should be present, got %v", res.Applied)
	}
}

func TestPromoteOverrideKeys(t *testing.T) {
	p, v := newPromoterFixture()
	const pass = "passphrase-32-bytes-long-exactly!"

	_ = v.Push("staging", pass, map[string]string{"FOO": "staging-value"})

	res, err := p.Promote("staging", "production", pass, PromoteOptions{
		OverrideKeys: map[string]string{"FOO": "prod-value"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Applied["FOO"] != "prod-value" {
		t.Errorf("expected FOO=prod-value, got %q", res.Applied["FOO"])
	}
}

func TestPromoteDryRun(t *testing.T) {
	p, v := newPromoterFixture()
	const pass = "passphrase-32-bytes-long-exactly!"

	_ = v.Push("staging", pass, map[string]string{"FOO": "bar"})

	_, err := p.Promote("staging", "production", pass, PromoteOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = v.Pull("production", pass)
	if !IsNotFound(err) {
		t.Error("dry run should not have written to destination")
	}
}

func TestPromoteNoDiff(t *testing.T) {
	p, v := newPromoterFixture()
	const pass = "passphrase-32-bytes-long-exactly!"

	vars := map[string]string{"FOO": "bar"}
	_ = v.Push("staging", pass, vars)
	_ = v.Push("production", pass, vars)

	res, err := p.Promote("staging", "production", pass, PromoteOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Diffs) != 0 {
		t.Errorf("expected no diffs when envs are identical, got %v", res.Diffs)
	}
}
