// Package env — access control
//
// This file documents the access-control sub-feature of the env package.
//
// AccessPolicy
//
// An AccessPolicy records which principals (users, service accounts, CI tokens)
// may interact with a named environment and at what level:
//
//   - AccessRead  — fetch / pull the encrypted env file
//   - AccessWrite — push updates to the env file
//   - AccessAdmin — full control, including managing the policy itself
//
// AccessStore
//
// AccessStore persists policies through the standard storage.Backend interface,
// so policies can be stored locally (FileBackend) or on any future remote
// backend without changes to the access-control logic.
//
// Typical usage:
//
//	store := env.NewAccessStore(backend)
//	policy := env.NewAccessPolicy("production")
//	policy.Grant("alice@example.com", env.AccessAdmin)
//	store.Save(policy)
package env
