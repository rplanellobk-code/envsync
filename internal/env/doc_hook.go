// Package env provides the hook subsystem for envsync lifecycle events.
//
// # Hook Registry
//
// A HookRegistry allows callers to register named callback functions
// (HookFunc) that are invoked at specific points in the push/pull
// lifecycle:
//
//	- pre-push  — before vars are encrypted and written to the backend
//	- post-push — after a successful push
//	- pre-pull  — before vars are fetched from the backend
//	- post-pull — after vars have been decrypted and returned
//
// Hooks are executed in registration order. If any hook returns a
// non-nil error the chain is aborted and the error is propagated to
// the caller with the hook name and event embedded for easy debugging.
//
// Example usage:
//
//	reg := env.NewHookRegistry()
//	reg.Register(env.HookEventPostPush, "notify", func(event env.HookEvent, environment string, vars map[string]string) error {
//		log.Printf("pushed %d vars to %s", len(vars), environment)
//		return nil
//	})
package env
