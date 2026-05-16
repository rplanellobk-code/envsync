// Package env provides the Sign and Verify helpers for tamper-evident
// .env maps.
//
// # Overview
//
// Sign computes an HMAC-SHA256 signature over all key=value pairs in a map
// (sorted lexicographically) and stores the hex digest under the reserved key
// _ENVSYNC_SIG. The original map is never mutated; a new map is returned.
//
// Verify recomputes the signature and compares it in constant time against the
// stored value. It returns ErrMissingSignature when the key is absent and
// ErrSignatureMismatch when the digest does not match, indicating that the
// contents were modified after signing.
//
// # Reserved keys
//
// The key _ENVSYNC_SIG is reserved for internal use. Callers should not set
// this key manually; doing so will cause Verify to treat the map as unsigned
// or produce a spurious ErrSignatureMismatch if the value does not match the
// computed digest.
//
// # Typical usage
//
//	signed := env.Sign(vars, passphrase)
//	// … store or transmit signed …
//	if err := env.Verify(received, passphrase); err != nil {
//		log.Fatal(err)
//	}
package env
