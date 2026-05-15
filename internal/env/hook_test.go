package env

import (
	"errors"
	"testing"
)

func newTestRegistry() *HookRegistry {
	return NewHookRegistry()
}

func TestHookRegisterAndRun(t *testing.T) {
	r := newTestRegistry()
	called := false
	err := r.Register(HookEventPostPush, "logger", func(event HookEvent, env string, vars map[string]string) error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := r.Run(HookEventPostPush, "production", nil); err != nil {
		t.Fatalf("run error: %v", err)
	}
	if !called {
		t.Fatal("expected hook to be called")
	}
}

func TestHookRunPassesArgs(t *testing.T) {
	r := newTestRegistry()
	vars := map[string]string{"KEY": "value"}
	var gotEvent HookEvent
	var gotEnv string
	var gotVars map[string]string
	_ = r.Register(HookEventPrePush, "capture", func(event HookEvent, env string, v map[string]string) error {
		gotEvent = event
		gotEnv = env
		gotVars = v
		return nil
	})
	_ = r.Run(HookEventPrePush, "staging", vars)
	if gotEvent != HookEventPrePush {
		t.Errorf("expected event %q, got %q", HookEventPrePush, gotEvent)
	}
	if gotEnv != "staging" {
		t.Errorf("expected env staging, got %q", gotEnv)
	}
	if gotVars["KEY"] != "value" {
		t.Errorf("expected vars to be passed")
	}
}

func TestHookDuplicateNameReturnsError(t *testing.T) {
	r := newTestRegistry()
	noop := func(HookEvent, string, map[string]string) error { return nil }
	_ = r.Register(HookEventPostPull, "dup", noop)
	err := r.Register(HookEventPostPull, "dup", noop)
	if err == nil {
		t.Fatal("expected error for duplicate hook name")
	}
}

func TestHookRunStopsOnError(t *testing.T) {
	r := newTestRegistry()
	sentinel := errors.New("hook failure")
	secondCalled := false
	_ = r.Register(HookEventPrePull, "fail", func(HookEvent, string, map[string]string) error {
		return sentinel
	})
	_ = r.Register(HookEventPrePull, "second", func(HookEvent, string, map[string]string) error {
		secondCalled = true
		return nil
	})
	err := r.Run(HookEventPrePull, "dev", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
	if secondCalled {
		t.Fatal("second hook should not have been called")
	}
}

func TestHookUnregister(t *testing.T) {
	r := newTestRegistry()
	noop := func(HookEvent, string, map[string]string) error { return nil }
	_ = r.Register(HookEventPostPush, "remove-me", noop)
	if err := r.Unregister(HookEventPostPush, "remove-me"); err != nil {
		t.Fatalf("unregister error: %v", err)
	}
	if names := r.List(HookEventPostPush); len(names) != 0 {
		t.Errorf("expected no hooks, got %v", names)
	}
}

func TestHookUnregisterNotFound(t *testing.T) {
	r := newTestRegistry()
	if err := r.Unregister(HookEventPostPush, "ghost"); err == nil {
		t.Fatal("expected error for missing hook")
	}
}

func TestHookList(t *testing.T) {
	r := newTestRegistry()
	noop := func(HookEvent, string, map[string]string) error { return nil }
	_ = r.Register(HookEventPostPull, "alpha", noop)
	_ = r.Register(HookEventPostPull, "beta", noop)
	names := r.List(HookEventPostPull)
	if len(names) != 2 || names[0] != "alpha" || names[1] != "beta" {
		t.Errorf("unexpected list: %v", names)
	}
}

func TestHookRegisterEmptyNameError(t *testing.T) {
	r := newTestRegistry()
	err := r.Register(HookEventPrePush, "", func(HookEvent, string, map[string]string) error { return nil })
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}
