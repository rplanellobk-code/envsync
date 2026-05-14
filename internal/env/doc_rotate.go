// Package env provides environment variable management utilities for envsync.
//
// # Passphrase Rotation
//
// The Rotator type enables safe re-encryption of a vault environment under a
// new passphrase without data loss:
//
//  1. The current ciphertext is decrypted with the old passphrase.
//  2. An optional snapshot is saved so the state can be recovered via
//     RollbackToLatest if the push step fails.
//  3. The plaintext is re-encrypted and stored under the new passphrase.
//  4. An optional audit entry is appended to record the rotation event.
//
// Use NewRotator to construct a Rotator, then call Rotate for a single
// environment or RotateAll to iterate every environment known to the vault.
//
// Example:
//
//	rot := env.NewRotator(vault, auditLog, snapshotStore)
//	_, err := rot.Rotate("production", oldPass, newPass)
package env
