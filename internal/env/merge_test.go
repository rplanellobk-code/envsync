package env

import (
	"errors"
	"testing"
)

func TestMergeNoConflict(t *testing.T) {
	local := map[string]string{"A": "1", "B": "2"}
	remote := map[string]string{"C": "3"}

	result, err := Merge(local, remote, PreferLocal)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["A"] != "1" || result["B"] != "2" || result["C"] != "3" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestMergePreferLocal(t *testing.T) {
	local := map[string]string{"KEY": "local_value"}
	remote := map[string]string{"KEY": "remote_value"}

	result, err := Merge(local, remote, PreferLocal)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["KEY"] != "local_value" {
		t.Errorf("expected local_value, got %q", result["KEY"])
	}
}

func TestMergePreferRemote(t *testing.T) {
	local := map[string]string{"KEY": "local_value"}
	remote := map[string]string{"KEY": "remote_value"}

	result, err := Merge(local, remote, PreferRemote)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["KEY"] != "remote_value" {
		t.Errorf("expected remote_value, got %q", result["KEY"])
	}
}

func TestMergeErrorOnConflict(t *testing.T) {
	local := map[string]string{"KEY": "local_value"}
	remote := map[string]string{"KEY": "remote_value"}

	_, err := Merge(local, remote, ErrorOnConflict)
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}

	var ce *ConflictError
	if !errors.As(err, &ce) {
		t.Fatalf("expected *ConflictError, got %T", err)
	}
	if ce.Key != "KEY" {
		t.Errorf("expected conflict key KEY, got %q", ce.Key)
	}
}

func TestMergeNoConflictSameValue(t *testing.T) {
	local := map[string]string{"KEY": "same"}
	remote := map[string]string{"KEY": "same"}

	result, err := Merge(local, remote, ErrorOnConflict)
	if err != nil {
		t.Fatalf("unexpected error for identical values: %v", err)
	}
	if result["KEY"] != "same" {
		t.Errorf("expected same, got %q", result["KEY"])
	}
}

func TestMergeRemoteOnlyKeys(t *testing.T) {
	local := map[string]string{}
	remote := map[string]string{"NEW": "value"}

	result, err := Merge(local, remote, ErrorOnConflict)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["NEW"] != "value" {
		t.Errorf("expected value, got %q", result["NEW"])
	}
}
