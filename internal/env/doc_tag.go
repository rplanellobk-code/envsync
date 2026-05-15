// Package env provides the Tag and TagStore types for labelling snapshots
// with human-readable names.
//
// # Overview
//
// Tags allow operators to pin a meaningful name (e.g. "v1.2.3" or
// "pre-migration") to a specific snapshot index within an environment.
// This makes it easy to roll back to a known-good state by name rather
// than by numeric index.
//
// # Usage
//
//	store := env.NewTagStore(backend)
//
//	// Create a tag pointing at snapshot 7 of the production environment.
//	err := store.Save(env.Tag{
//		Name:        "v2.0.0",
//		Environment: "production",
//		SnapshotIdx: 7,
//		Message:     "stable release",
//	})
//
//	// Retrieve the tag later.
//	tag, err := store.Get("production", "v2.0.0")
package env
