// Package env provides environment variable management utilities for envsync.
//
// The audit sub-feature records every push, pull, and delete operation
// performed against a remote environment. Entries are stored as JSON arrays
// in the configured storage backend under the "audit/<environment>.json" key.
//
// Example usage:
//
//	log := env.NewAuditLog(backend)
//	_ = log.Record(env.AuditEntry{
//		Environment: "production",
//		Action:      env.ActionPush,
//		User:        "alice",
//	})
//	entries, _ := log.List("production")
package env
