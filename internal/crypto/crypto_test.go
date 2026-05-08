package crypto_test

import (
	"bytes"
	"testing"

	"github.com/yourorg/envsync/internal/crypto"
)

func TestDeriveKey(t *testing.T) {
	key := crypto.DeriveKey("my-secret-passphrase")
	if len(key) != 32 {
		t.Fatalf("expected key length 32, got %d", len(key))
	}

	key2 := crypto.DeriveKey("my-secret-passphrase")
	if !bytes.Equal(key, key2) {
		t.Fatal("expected same key for same passphrase")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key := crypto.DeriveKey("test-passphrase")
	plaintext := []byte("DB_PASSWORD=supersecret\nAPI_KEY=abc123")

	ciphertext, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("encrypt error: %v", err)
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Fatal("ciphertext should not equal plaintext")
	}

	decrypted, err := crypto.Decrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("decrypt error: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("expected %q, got %q", plaintext, decrypted)
	}
}

func TestDecryptWrongKey(t *testing.T) {
	key := crypto.DeriveKey("correct-passphrase")
	wrongKey := crypto.DeriveKey("wrong-passphrase")
	plaintext := []byte("SECRET=value")

	ciphertext, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("encrypt error: %v", err)
	}

	_, err = crypto.Decrypt(wrongKey, ciphertext)
	if err == nil {
		t.Fatal("expected error when decrypting with wrong key")
	}
}

func TestDecryptTooShort(t *testing.T) {
	key := crypto.DeriveKey("passphrase")
	_, err := crypto.Decrypt(key, []byte("short"))
	if err == nil {
		t.Fatal("expected error for too-short ciphertext")
	}
}
