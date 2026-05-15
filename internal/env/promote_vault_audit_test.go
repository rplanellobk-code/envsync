package env

import (
	"testing"

	"github.com/user/envsync/internal/storage"
)

func newAuditedPromoterFixture(actor string) (*PromoterWithAudit, *Vault, *AuditLog) {
	backend := storage.NewMemoryBackend()
	v := NewVault(backend)
	auditLog := NewAuditLog(backend)
	pa := NewPromoterWithAudit(v, auditLog, actor)
	return pa, v, auditLog
}

func TestAuditedPromoteRecordsEntry(t *testing.T) {
	pa, v, auditLog := newAuditedPromoterFixture("ci-bot")
	const pass = "passphrase-32-bytes-long-exactly!"

	_ = v.Push("staging", pass, map[string]string{"KEY": "val"})

	_, err := pa.Promote("staging", "production", pass, PromoteOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := auditLog.List("production")
	if err != nil {
		t.Fatalf("list audit: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 audit entry, got %d", len(entries))
	}
	if entries[0].Actor != "ci-bot" {
		t.Errorf("expected actor ci-bot, got %q", entries[0].Actor)
	}
	if !entries[0].Success {
		t.Error("expected success=true")
	}
}

func TestAuditedPromoteDryRunRecordsEntry(t *testing.T) {
	pa, v, auditLog := newAuditedPromoterFixture("human")
	const pass = "passphrase-32-bytes-long-exactly!"

	_ = v.Push("staging", pass, map[string]string{"A": "1"})

	_, err := pa.Promote("staging", "production", pass, PromoteOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := auditLog.List("production")
	if err != nil {
		t.Fatalf("list audit: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 audit entry, got %d", len(entries))
	}
	if entries[0].Action != "promote:staging->production:dry-run" {
		t.Errorf("unexpected action: %q", entries[0].Action)
	}
}
