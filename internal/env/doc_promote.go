// Package env provides the Promoter type for copying environment variables
// from one environment to another via an encrypted Vault backend.
//
// # Overview
//
// A Promoter reads the source environment, applies optional transformations
// (key omission, value overrides), diffs the result against the destination,
// and optionally writes the merged result back.
//
// # Usage
//
//	v := env.NewVault(backend)
//	p := env.NewPromoter(v)
//	res, err := p.Promote("staging", "production", passphrase, env.PromoteOptions{
//		OmitKeys:     []string{"DEBUG"},
//		OverrideKeys: map[string]string{"LOG_LEVEL": "warn"},
//	})
//
// Set DryRun: true to inspect diffs without persisting changes.
package env
