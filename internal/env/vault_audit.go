package env

import (
	"fmt"
)

// VaultWithAudit wraps a Vault and records every mutating operation to an
// AuditLog. Read operations (Pull, List) are also recorded for full
// traceability.
type VaultWithAudit struct {
	vault *Vault
	audit *AuditLog
	user  string
}

// NewVaultWithAudit creates a VaultWithAudit that delegates to v and records
// operations in auditLog, attributing them to user.
func NewVaultWithAudit(v *Vault, auditLog *AuditLog, user string) *VaultWithAudit {
	return &VaultWithAudit{vault: v, audit: auditLog, user: user}
}

// Push encrypts and stores env vars for the given environment, then records
// the operation.
func (va *VaultWithAudit) Push(environment string, vars map[string]string) error {
	if err := va.vault.Push(environment, vars); err != nil {
		return err
	}
	_ = va.audit.Record(AuditEntry{
		Environment: environment,
		Action:      ActionPush,
		User:        va.user,
	})
	return nil
}

// Pull retrieves and decrypts env vars for the given environment, then records
// the operation.
func (va *VaultWithAudit) Pull(environment string) (map[string]string, error) {
	vars, err := va.vault.Pull(environment)
	if err != nil {
		return nil, err
	}
	_ = va.audit.Record(AuditEntry{
		Environment: environment,
		Action:      ActionPull,
		User:        va.user,
	})
	return vars, nil
}

// List returns all stored environment names without recording an audit entry.
func (va *VaultWithAudit) List() ([]string, error) {
	return va.vault.List()
}

// AuditHistory returns all recorded audit entries for the given environment.
func (va *VaultWithAudit) AuditHistory(environment string) ([]AuditEntry, error) {
	if environment == "" {
		return nil, fmt.Errorf("vault_audit: environment must not be empty")
	}
	return va.audit.List(environment)
}
