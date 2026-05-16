package env

import (
	"testing"
)

func TestSignProducesSignatureKey(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	signed := Sign(env, "mysecret")
	if _, ok := signed[signatureKey]; !ok {
		t.Fatal("expected signature key to be present after Sign")
	}
}

func TestVerifyValidSignature(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	signed := Sign(env, "mysecret")
	if err := Verify(signed, "mysecret"); err != nil {
		t.Fatalf("expected valid signature, got: %v", err)
	}
}

func TestVerifyWrongSecret(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	signed := Sign(env, "mysecret")
	err := Verify(signed, "wrongsecret")
	if err != ErrSignatureMismatch {
		t.Fatalf("expected ErrSignatureMismatch, got: %v", err)
	}
}

func TestVerifyMissingSignature(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	err := Verify(env, "mysecret")
	if err != ErrMissingSignature {
		t.Fatalf("expected ErrMissingSignature, got: %v", err)
	}
}

func TestVerifyTamperedValue(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	signed := Sign(env, "mysecret")
	signed["FOO"] = "tampered"
	err := Verify(signed, "mysecret")
	if err != ErrSignatureMismatch {
		t.Fatalf("expected ErrSignatureMismatch after tampering, got: %v", err)
	}
}

func TestVerifyTamperedKey(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	signed := Sign(env, "mysecret")
	signed["EXTRA"] = "injected"
	err := Verify(signed, "mysecret")
	if err != ErrSignatureMismatch {
		t.Fatalf("expected ErrSignatureMismatch after key injection, got: %v", err)
	}
}

func TestSignIsDeterministic(t *testing.T) {
	env := map[string]string{"Z": "last", "A": "first", "M": "middle"}
	s1 := Sign(env, "secret")
	s2 := Sign(env, "secret")
	if s1[signatureKey] != s2[signatureKey] {
		t.Fatal("Sign should be deterministic regardless of map iteration order")
	}
}

func TestSignDoesNotMutateOriginal(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	_ = Sign(env, "secret")
	if _, ok := env[signatureKey]; ok {
		t.Fatal("Sign must not mutate the original map")
	}
}
