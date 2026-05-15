package env

// Clone copies all key-value pairs from one environment to another,
// optionally overwriting existing keys in the destination.
//
// It uses the Vault to pull the source environment, applies optional
// key filtering and overrides, then pushes the result to the destination.

import (
	"fmt"
)

// CloneOptions controls how the clone operation behaves.
type CloneOptions struct {
	// Overwrite replaces existing keys in the destination.
	Overwrite bool
	// OmitKeys lists keys to exclude from the clone.
	OmitKeys []string
	// OverrideKeys replaces specific key values before writing to destination.
	OverrideKeys map[string]string
}

// Cloner copies an environment from one name to another via a Vault.
type Cloner struct {
	vault *Vault
}

// NewCloner returns a Cloner backed by the given Vault.
func NewCloner(v *Vault) *Cloner {
	return &Cloner{vault: v}
}

// Clone copies srcEnv into dstEnv according to opts.
// If opts.Overwrite is false, existing keys in dstEnv are preserved.
func (c *Cloner) Clone(passphrase, srcEnv, dstEnv string, opts CloneOptions) error {
	if srcEnv == "" {
		return fmt.Errorf("clone: source environment must not be empty")
	}
	if dstEnv == "" {
		return fmt.Errorf("clone: destination environment must not be empty")
	}

	src, err := c.vault.Pull(passphrase, srcEnv)
	if err != nil {
		return fmt.Errorf("clone: pull source %q: %w", srcEnv, err)
	}

	omit := make(map[string]bool, len(opts.OmitKeys))
	for _, k := range opts.OmitKeys {
		omit[k] = true
	}

	filtered := make(map[string]string, len(src))
	for k, v := range src {
		if omit[k] {
			continue
		}
		if override, ok := opts.OverrideKeys[k]; ok {
			v = override
		}
		filtered[k] = v
	}

	if !opts.Overwrite {
		existing, err := c.vault.Pull(passphrase, dstEnv)
		if err != nil && !IsNotFound(err) {
			return fmt.Errorf("clone: pull destination %q: %w", dstEnv, err)
		}
		for k, v := range existing {
			filtered[k] = v
		}
	}

	if err := c.vault.Push(passphrase, dstEnv, filtered); err != nil {
		return fmt.Errorf("clone: push destination %q: %w", dstEnv, err)
	}
	return nil
}
