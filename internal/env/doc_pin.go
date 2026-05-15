// Package env — Pin module
//
// A Pin is a named, immutable reference that points to a specific snapshot
// of an environment's variables. Pins are analogous to release tags in
// version-control systems: once created they cannot be overwritten, ensuring
// a stable anchor that can be used for rollback, auditing, or promotion.
//
// Usage:
//
//	store := env.NewPinStore(backend)
//
//	// Create a pin after a successful deployment
//	err := store.Save(env.Pin{
//		Name:        "v1.2.0",
//		Environment: "production",
//		SnapshotID:  snap.ID,
//		CreatedAt:   time.Now().UTC(),
//		Note:        "post-deploy pin",
//	})
//
//	// Retrieve a pin later
//	pin, err := store.Get("production", "v1.2.0")
//
//	// List all pins for an environment
//	pins, err := store.List("production")
package env
