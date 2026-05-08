package env

import (
	"testing"
)

func TestDiffNoChanges(t *testing.T) {
	base := map[string]string{"FOO": "bar", "BAZ": "qux"}
	target := map[string]string{"FOO": "bar", "BAZ": "qux"}
	changes := Diff(base, target)
	if len(changes) != 0 {
		t.Fatalf("expected no changes, got %d", len(changes))
	}
}

func TestDiffAdded(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	target := map[string]string{"FOO": "bar", "NEW_KEY": "newval"}
	changes := Diff(base, target)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != ChangeAdded || changes[0].Key != "NEW_KEY" || changes[0].NewValue != "newval" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestDiffRemoved(t *testing.T) {
	base := map[string]string{"FOO": "bar", "OLD_KEY": "oldval"}
	target := map[string]string{"FOO": "bar"}
	changes := Diff(base, target)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != ChangeRemoved || changes[0].Key != "OLD_KEY" || changes[0].OldValue != "oldval" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestDiffUpdated(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	target := map[string]string{"FOO": "baz"}
	changes := Diff(base, target)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	c := changes[0]
	if c.Type != ChangeUpdated || c.Key != "FOO" || c.OldValue != "bar" || c.NewValue != "baz" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestDiffDeterministicOrder(t *testing.T) {
	base := map[string]string{}
	target := map[string]string{"Z_KEY": "1", "A_KEY": "2", "M_KEY": "3"}
	changes := Diff(base, target)
	if len(changes) != 3 {
		t.Fatalf("expected 3 changes, got %d", len(changes))
	}
	if changes[0].Key != "A_KEY" || changes[1].Key != "M_KEY" || changes[2].Key != "Z_KEY" {
		t.Errorf("changes not sorted: %v", changes)
	}
}

func TestApply(t *testing.T) {
	base := map[string]string{"FOO": "bar", "OLD": "val"}
	changes := []Change{
		{Key: "FOO", Type: ChangeUpdated, OldValue: "bar", NewValue: "newbar"},
		{Key: "OLD", Type: ChangeRemoved, OldValue: "val"},
		{Key: "FRESH", Type: ChangeAdded, NewValue: "fresh"},
	}
	result := Apply(base, changes)
	if result["FOO"] != "newbar" {
		t.Errorf("expected FOO=newbar, got %s", result["FOO"])
	}
	if _, ok := result["OLD"]; ok {
		t.Error("expected OLD to be removed")
	}
	if result["FRESH"] != "fresh" {
		t.Errorf("expected FRESH=fresh, got %s", result["FRESH"])
	}
	// base should not be mutated
	if base["FOO"] != "bar" {
		t.Error("base map was mutated")
	}
}
