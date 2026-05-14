package env

import (
	"fmt"
	"time"
)

// RotateResult holds the outcome of a passphrase rotation for one environment.
type RotateResult struct {
	Environment string
	RotatedAt   time.Time
}

// Rotator handles re-encrypting vault entries under a new passphrase.
type Rotator struct {
	vault    *Vault
	audit    *AuditLog
	snapshot *SnapshotStore
}

// NewRotator creates a Rotator backed by the given vault, audit log, and
// snapshot store. audit and snapshot may be nil.
func NewRotator(v *Vault, a *AuditLog, s *SnapshotStore) *Rotator {
	return &Rotator{vault: v, audit: a, snapshot: s}
}

// Rotate re-encrypts the named environment: it pulls with oldPass, then pushes
// with newPass. A snapshot is saved before the operation when a store is
// available. An audit entry is appended on success.
func (r *Rotator) Rotate(env, oldPass, newPass string) (RotateResult, error) {
	if env == "" {
		return RotateResult{}, fmt.Errorf("rotate: environment name must not be empty")
	}
	if oldPass == newPass {
		return RotateResult{}, fmt.Errorf("rotate: new passphrase must differ from old")
	}

	pairs, err := r.vault.Pull(env, oldPass)
	if err != nil {
		return RotateResult{}, fmt.Errorf("rotate: pull with old passphrase: %w", err)
	}

	if r.snapshot != nil {
		snap := NewSnapshot(env, pairs)
		if saveErr := r.snapshot.Save(snap); saveErr != nil {
			return RotateResult{}, fmt.Errorf("rotate: save snapshot: %w", saveErr)
		}
	}

	if err := r.vault.Push(env, pairs, newPass); err != nil {
		return RotateResult{}, fmt.Errorf("rotate: push with new passphrase: %w", err)
	}

	result := RotateResult{Environment: env, RotatedAt: time.Now().UTC()}

	if r.audit != nil {
		_ = r.audit.Record(AuditEntry{
			Environment: env,
			Action:      "rotate",
			Timestamp:   result.RotatedAt,
		})
	}

	return result, nil
}

// RotateAll rotates every environment returned by the vault listing.
func (r *Rotator) RotateAll(oldPass, newPass string) ([]RotateResult, error) {
	envs, err := r.vault.List()
	if err != nil {
		return nil, fmt.Errorf("rotate-all: list environments: %w", err)
	}
	var results []RotateResult
	for _, env := range envs {
		res, err := r.Rotate(env, oldPass, newPass)
		if err != nil {
			return results, fmt.Errorf("rotate-all: %w", err)
		}
		results = append(results, res)
	}
	return results, nil
}
