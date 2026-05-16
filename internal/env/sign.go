package env

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
)

// ErrSignatureMismatch is returned when a signature does not match the env map.
var ErrSignatureMismatch = errors.New("signature mismatch: env contents may have been tampered with")

// ErrMissingSignature is returned when no signature is found during verification.
var ErrMissingSignature = errors.New("no signature found in env map")

const signatureKey = "_ENVSYNC_SIG"

// Sign computes an HMAC-SHA256 signature over the sorted key=value pairs in
// env (excluding any existing signature key) using the provided secret, and
// inserts the hex-encoded signature under _ENVSYNC_SIG.
func Sign(env map[string]string, secret string) map[string]string {
	sig := computeSignature(env, secret)
	out := make(map[string]string, len(env)+1)
	for k, v := range env {
		if k != signatureKey {
			out[k] = v
		}
	}
	out[signatureKey] = sig
	return out
}

// Verify checks that the _ENVSYNC_SIG value in env matches a freshly computed
// HMAC over the remaining keys. Returns ErrMissingSignature or
// ErrSignatureMismatch on failure.
func Verify(env map[string]string, secret string) error {
	stored, ok := env[signatureKey]
	if !ok {
		return ErrMissingSignature
	}
	expected := computeSignature(env, secret)
	if !hmac.Equal([]byte(stored), []byte(expected)) {
		return ErrSignatureMismatch
	}
	return nil
}

func computeSignature(env map[string]string, secret string) string {
	keys := make([]string, 0, len(env))
	for k := range env {
		if k != signatureKey {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, env[k]))
	}
	payload := strings.Join(parts, "\n")

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}
