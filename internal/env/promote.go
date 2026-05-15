package env

import (
	"fmt"
)

// PromoteOptions controls how a promotion between environments is performed.
type PromoteOptions struct {
	// OmitKeys lists keys to exclude from the promotion.
	OmitKeys []string
	// OverrideKeys are key=value pairs that override values during promotion.
	OverrideKeys map[string]string
	// DryRun when true returns the result without writing to the destination.
	DryRun bool
}

// PromoteResult holds the outcome of a promotion.
type PromoteResult struct {
	Source      string
	Destination string
	Applied     map[string]string
	Diffs       []DiffEntry
}

// Promoter copies env vars from one environment to another via a Vault.
type Promoter struct {
	vault *Vault
}

// NewPromoter creates a Promoter backed by the given Vault.
func NewPromoter(v *Vault) *Promoter {
	return &Promoter{vault: v}
}

// Promote reads the source environment, applies options, and writes the result
// to the destination environment. It returns a PromoteResult describing what
// changed.
func (p *Promoter) Promote(src, dst, passphrase string, opts PromoteOptions) (*PromoteResult, error) {
	srcVars, err := p.vault.Pull(src, passphrase)
	if err != nil {
		return nil, fmt.Errorf("promote: pull source %q: %w", src, err)
	}

	omit := make(map[string]bool, len(opts.OmitKeys))
	for _, k := range opts.OmitKeys {
		omit[k] = true
	}

	candidate := make(map[string]string, len(srcVars))
	for k, v := range srcVars {
		if omit[k] {
			continue
		}
		candidate[k] = v
	}
	for k, v := range opts.OverrideKeys {
		candidate[k] = v
	}

	dstVars, err := p.vault.Pull(dst, passphrase)
	if err != nil && !IsNotFound(err) {
		return nil, fmt.Errorf("promote: pull destination %q: %w", dst, err)
	}
	if dstVars == nil {
		dstVars = map[string]string{}
	}

	diffs := Diff(dstVars, candidate)

	result := &PromoteResult{
		Source:      src,
		Destination: dst,
		Applied:     candidate,
		Diffs:       diffs,
	}

	if opts.DryRun {
		return result, nil
	}

	if err := p.vault.Push(dst, passphrase, candidate); err != nil {
		return nil, fmt.Errorf("promote: push destination %q: %w", dst, err)
	}

	return result, nil
}
