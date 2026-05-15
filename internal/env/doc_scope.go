// Package env provides the core environment variable management primitives
// used by envsync.
//
// # Scopes
//
// A Scope is a named filter over environment variable keys. It can match keys
// by a string prefix, a regular expression pattern, or both. Scopes allow
// operations such as diff, export, and import to be narrowed to a logical
// subset of variables (e.g. only database-related keys).
//
// Use ScopeStore to persist and retrieve Scope definitions in any storage
// backend:
//
//	store := env.NewScopeStore(backend)
//	store.Save(ctx, &env.Scope{Name: "db", Prefix: "DB_"})
//	scope, _ := store.Load(ctx, "db")
//	filtered := scope.Filter(myEnvMap)
package env
