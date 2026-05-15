package env

import (
	"fmt"
	"time"
)

// PromoterWithAudit wraps a Promoter and records each promotion in an AuditLog.
type PromoterWithAudit struct {
	promoter *Promoter
	audit    *AuditLog
	actor    string
}

// NewPromoterWithAudit creates a PromoterWithAudit.
func NewPromoterWithAudit(v *Vault, audit *AuditLog, actor string) *PromoterWithAudit {
	return &PromoterWithAudit{
		promoter: NewPromoter(v),
		audit:    audit,
		actor:    actor,
	}
}

// Promote delegates to the underlying Promoter and records the outcome.
func (pa *PromoterWithAudit) Promote(src, dst, passphrase string, opts PromoteOptions) (*PromoteResult, error) {
	res, err := pa.promoter.Promote(src, dst, passphrase, opts)

	action := fmt.Sprintf("promote:%s->%s", src, dst)
	if opts.DryRun {
		action += ":dry-run"
	}

	auditErr := pa.audit.Record(AuditEntry{
		Timestamp:   time.Now().UTC(),
		Actor:       pa.actor,
		Action:      action,
		Environment: dst,
		Success:     err == nil,
	})
	if auditErr != nil && err == nil {
		return res, fmt.Errorf("promote audit: %w", auditErr)
	}

	return res, err
}
