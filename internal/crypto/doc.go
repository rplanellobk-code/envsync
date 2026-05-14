// Package crypto provides symmetric encryption and decryption utilities
// for envsync. It uses AES-256-GCM for authenticated encryption, ensuring
// both confidentiality and integrity of .env file contents.
//
// Keys are derived from user-supplied passphrases using SHA-256, producing
// a 32-byte key suitable for AES-256. Each encryption call generates a
// unique random nonce (12 bytes), which is prepended to the ciphertext
// for storage and transport. The resulting format is:
//
//	[ 12-byte nonce ][ ciphertext + 16-byte GCM auth tag ]
//
// Decryption expects input in the same format and will return an error if
// the authentication tag does not match (i.e. the data has been tampered
// with or the wrong passphrase was supplied).
package crypto
