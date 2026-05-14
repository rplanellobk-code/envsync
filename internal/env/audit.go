package env

import (
	"encoding/json"
	"fmt"
	"time"

	"envsync/internal/storage"
)

// AuditAction represents the type of operation performed.
type AuditAction string

const (
	ActionPush   AuditAction = "push"
	ActionPull   AuditAction = "pull"
	ActionDelete AuditAction = "delete"
)

// AuditEntry records a single operation against a remote environment.
type AuditEntry struct {
	Timestamp   time.Time   `json:"timestamp"`
	Environment string      `json:"environment"`
	Action      AuditAction `json:"action"`
	User        string      `json:"user,omitempty"`
	Note        string      `json:"note,omitempty"`
}

// AuditLog manages persisted audit entries for an environment.
type AuditLog struct {
	backend storage.Backend
}

// NewAuditLog creates an AuditLog backed by the given storage backend.
func NewAuditLog(b storage.Backend) *AuditLog {
	return &AuditLog{backend: b}
}

// Record appends a new entry to the audit log for the given environment.
func (a *AuditLog) Record(entry AuditEntry) error {
	if entry.Environment == "" {
		return fmt.Errorf("audit: environment must not be empty")
	}
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}

	existing, err := a.List(entry.Environment)
	if err != nil && !storage.IsNotFound(err) {
		return fmt.Errorf("audit: load existing log: %w", err)
	}
	existing = append(existing, entry)

	data, err := json.Marshal(existing)
	if err != nil {
		return fmt.Errorf("audit: marshal: %w", err)
	}
	return a.backend.Put(auditKey(entry.Environment), data)
}

// List returns all audit entries for the given environment.
func (a *AuditLog) List(environment string) ([]AuditEntry, error) {
	data, err := a.backend.Get(auditKey(environment))
	if err != nil {
		return nil, err
	}
	var entries []AuditEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("audit: unmarshal: %w", err)
	}
	return entries, nil
}

func auditKey(environment string) string {
	return fmt.Sprintf("audit/%s.json", environment)
}
