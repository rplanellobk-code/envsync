// Package env provides environment variable management primitives for envsync.
//
// # Advisory Locking
//
// LockStore provides cooperative, TTL-based advisory locks for environments.
// Locks prevent concurrent push operations from overwriting each other when
// multiple operators work against the same remote backend simultaneously.
//
// Locks are stored as JSON documents in the configured storage backend under
// the "locks/<environment>.json" key.
//
// Usage:
//
//	store := env.NewLockStore(backend)
//
//	lock, err := store.Acquire("production", "alice", 5*time.Minute)
//	if err != nil {
//	    // environment is locked by someone else
//	}
//	defer store.Release("production", "alice")
//
// Locks are advisory: they rely on all callers checking before writing.
// Expired locks are automatically superseded on the next Acquire call.
package env
