package env_test

import (
	"testing"

	"envsync/internal/env"
	"envsync/internal/storage"
)

func newAccessStore(t *testing.T) *env.AccessStore {
	t.Helper()
	b, err := storage.NewFileBackend(t.TempDir())
	if err != nil {
		t.Fatalf("create backend: %v", err)
	}
	return env.NewAccessStore(b)
}

func TestParseAccessLevel(t *testing.T) {
	cases := []struct{ in string; want env.AccessLevel }{
		{"read", env.AccessRead},
		{"write", env.AccessWrite},
		{"admin", env.AccessAdmin},
		{"READ", env.AccessRead},
	}
	for _, c := range cases {
		got, err := env.ParseAccessLevel(c.in)
		if err != nil || got != c.want {
			t.Errorf("ParseAccessLevel(%q) = %v, %v; want %v", c.in, got, err, c.want)
		}
	}
	if _, err := env.ParseAccessLevel("superuser"); err == nil {
		t.Error("expected error for unknown level")
	}
}

func TestAccessPolicyGrantCheck(t *testing.T) {
	p := env.NewAccessPolicy("production")
	_ = p.Grant("alice", env.AccessAdmin)
	_ = p.Grant("bob", env.AccessRead)

	if got := p.Check("alice"); got != env.AccessAdmin {
		t.Errorf("alice: got %v, want admin", got)
	}
	if got := p.Check("bob"); got != env.AccessRead {
		t.Errorf("bob: got %v, want read", got)
	}
	if got := p.Check("charlie"); got != env.AccessNone {
		t.Errorf("charlie: got %v, want none", got)
	}
}

func TestAccessPolicyRevoke(t *testing.T) {
	p := env.NewAccessPolicy("staging")
	_ = p.Grant("alice", env.AccessWrite)
	p.Revoke("alice")
	if got := p.Check("alice"); got != env.AccessNone {
		t.Errorf("after revoke: got %v, want none", got)
	}
}

func TestAccessStoreRoundTrip(t *testing.T) {
	store := newAccessStore(t)
	p := env.NewAccessPolicy("dev")
	_ = p.Grant("alice", env.AccessAdmin)
	_ = p.Grant("bob", env.AccessRead)

	if err := store.Save(p); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := store.Load("dev")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Check("alice") != env.AccessAdmin {
		t.Error("alice should be admin after round-trip")
	}
	if loaded.Check("bob") != env.AccessRead {
		t.Error("bob should be read after round-trip")
	}
}

func TestAccessStoreListEnvironments(t *testing.T) {
	store := newAccessStore(t)
	for _, name := range []string{"dev", "staging", "prod"} {
		p := env.NewAccessPolicy(name)
		_ = p.Grant("alice", env.AccessRead)
		if err := store.Save(p); err != nil {
			t.Fatalf("Save %s: %v", name, err)
		}
	}
	envs, err := store.ListEnvironments()
	if err != nil {
		t.Fatalf("ListEnvironments: %v", err)
	}
	if len(envs) != 3 {
		t.Errorf("expected 3 environments, got %d: %v", len(envs), envs)
	}
}
