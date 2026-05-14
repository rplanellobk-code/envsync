package env

import (
	"testing"
	"time"

	"envsync/internal/storage"
)

func newTestAuditLog(t *testing.T) *AuditLog {
	t.Helper()
	dir := t.TempDir()
	b, err := storage.NewFileBackend(dir)
	if err != nil {
		t.Fatalf("NewFileBackend: %v", err)
	}
	return NewAuditLog(b)
}

func TestAuditRecordAndList(t *testing.T) {
	log := newTestAuditLog(t)

	entry := AuditEntry{
		Environment: "production",
		Action:      ActionPush,
		User:        "alice",
	}
	if err := log.Record(entry); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := log.List("production")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Action != ActionPush {
		t.Errorf("expected action %q, got %q", ActionPush, entries[0].Action)
	}
	if entries[0].User != "alice" {
		t.Errorf("expected user alice, got %q", entries[0].User)
	}
	if entries[0].Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestAuditMultipleEntries(t *testing.T) {
	log := newTestAuditLog(t)

	actions := []AuditAction{ActionPush, ActionPull, ActionDelete}
	for _, a := range actions {
		if err := log.Record(AuditEntry{Environment: "staging", Action: a}); err != nil {
			t.Fatalf("Record %q: %v", a, err)
		}
	}

	entries, err := log.List("staging")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestAuditListNotFound(t *testing.T) {
	log := newTestAuditLog(t)
	_, err := log.List("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing environment, got nil")
	}
	if !storage.IsNotFound(err) {
		t.Errorf("expected not-found error, got %v", err)
	}
}

func TestAuditRecordEmptyEnvironment(t *testing.T) {
	log := newTestAuditLog(t)
	err := log.Record(AuditEntry{Action: ActionPush})
	if err == nil {
		t.Fatal("expected error for empty environment")
	}
}

func TestAuditTimestampPreserved(t *testing.T) {
	log := newTestAuditLog(t)
	ts := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	entry := AuditEntry{Environment: "dev", Action: ActionPull, Timestamp: ts}
	if err := log.Record(entry); err != nil {
		t.Fatalf("Record: %v", err)
	}
	entries, _ := log.List("dev")
	if !entries[0].Timestamp.Equal(ts) {
		t.Errorf("timestamp not preserved: got %v", entries[0].Timestamp)
	}
}
