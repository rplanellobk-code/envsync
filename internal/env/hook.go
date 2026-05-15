package env

import (
	"errors"
	"fmt"
)

// HookEvent represents the lifecycle event that triggered a hook.
type HookEvent string

const (
	HookEventPrePush  HookEvent = "pre-push"
	HookEventPostPush HookEvent = "post-push"
	HookEventPrePull  HookEvent = "pre-pull"
	HookEventPostPull HookEvent = "post-pull"
)

// HookFunc is a function invoked during a lifecycle event.
// It receives the environment name and the current env map (may be nil for pre-events).
type HookFunc func(event HookEvent, environment string, vars map[string]string) error

// HookRegistry holds named hooks registered for specific lifecycle events.
type HookRegistry struct {
	hooks map[HookEvent][]namedHook
}

type namedHook struct {
	name string
	fn   HookFunc
}

// NewHookRegistry creates an empty HookRegistry.
func NewHookRegistry() *HookRegistry {
	return &HookRegistry{
		hooks: make(map[HookEvent][]namedHook),
	}
}

// Register adds a named hook for the given event.
// Returns an error if a hook with the same name is already registered for that event.
func (r *HookRegistry) Register(event HookEvent, name string, fn HookFunc) error {
	if name == "" {
		return errors.New("hook name must not be empty")
	}
	if fn == nil {
		return errors.New("hook function must not be nil")
	}
	for _, h := range r.hooks[event] {
		if h.name == name {
			return fmt.Errorf("hook %q already registered for event %q", name, event)
		}
	}
	r.hooks[event] = append(r.hooks[event], namedHook{name: name, fn: fn})
	return nil
}

// Unregister removes a named hook from the given event.
// Returns an error if the hook is not found.
func (r *HookRegistry) Unregister(event HookEvent, name string) error {
	list := r.hooks[event]
	for i, h := range list {
		if h.name == name {
			r.hooks[event] = append(list[:i], list[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("hook %q not found for event %q", name, event)
}

// Run executes all hooks registered for the given event in registration order.
// If any hook returns an error, execution stops and the error is returned.
func (r *HookRegistry) Run(event HookEvent, environment string, vars map[string]string) error {
	for _, h := range r.hooks[event] {
		if err := h.fn(event, environment, vars); err != nil {
			return fmt.Errorf("hook %q (%s): %w", h.name, event, err)
		}
	}
	return nil
}

// List returns the names of all hooks registered for the given event.
func (r *HookRegistry) List(event HookEvent) []string {
	list := r.hooks[event]
	names := make([]string, len(list))
	for i, h := range list {
		names[i] = h.name
	}
	return names
}
