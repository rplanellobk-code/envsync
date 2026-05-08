// Package crypto provides symmetric encryption and decryption utilities
// for envsync. It uses AES-256-GCM for authenticated encryption, ensuring
// both confidentiality and integrity of .env file contents.
//
// Keys are derived from user-supplied passphrases using SHA-256.
// Each encryption call generates a unique random nonce, which is prepended
// to the ciphertext for storage and transport.
package crypto
