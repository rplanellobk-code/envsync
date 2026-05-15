// Package env provides utilities for managing environment variable maps.
//
// # Redaction
//
// The redact module offers helpers for masking sensitive values before they
// are displayed, logged, or exported.
//
// IsSensitiveKey uses a heuristic regular expression to detect keys that
// commonly hold secrets (passwords, tokens, API keys, etc.).
//
// Redact returns a shallow copy of an env map with all sensitive-looking
// values replaced by a fixed-width mask string.  The original map is never
// modified.
//
// RedactKeys provides explicit control: only the keys listed by the caller
// are masked, regardless of their names.
//
// RedactOptions lets callers customise the mask character and optionally
// preserve a visible suffix of the original value (e.g. to show the last
// four characters of a token for identification purposes).
package env
