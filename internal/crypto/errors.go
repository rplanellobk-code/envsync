package crypto

import "errors"

// Sentinel errors for the crypto package.
var (
	// ErrInvalidKeySize is returned when a key of incorrect length is provided.
	ErrInvalidKeySize = errors.New("crypto: invalid key size, expected 32 bytes")

	// ErrDecryptionFailed is returned when authenticated decryption fails,
	// typically indicating a wrong key or tampered ciphertext.
	ErrDecryptionFailed = errors.New("crypto: decryption failed, wrong key or corrupted data")

	// ErrEmptyPassphrase is returned when an empty passphrase is provided.
	ErrEmptyPassphrase = errors.New("crypto: passphrase must not be empty")
)

// ValidateKey returns an error if the key is not exactly 32 bytes.
func ValidateKey(key []byte) error {
	if len(key) != 32 {
		return ErrInvalidKeySize
	}
	return nil
}

// ValidatePassphrase returns an error if the passphrase is empty.
func ValidatePassphrase(passphrase string) error {
	if passphrase == "" {
		return ErrEmptyPassphrase
	}
	return nil
}
