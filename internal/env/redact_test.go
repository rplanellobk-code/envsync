package env

import (
	"strings"
	"testing"
)

func TestIsSensitiveKey(t *testing.T) {
	sensitive := []string{
		"PASSWORD", "DB_PASSWORD", "API_KEY", "API_SECRET",
		"AUTH_TOKEN", "PRIVATE_KEY", "ACCESS_KEY", "SIGNING_KEY",
		"aws_secret", "github_token", "stripe_api_key",
	}
	for _, k := range sensitive {
		if !IsSensitiveKey(k) {
			t.Errorf("expected %q to be sensitive", k)
		}
	}

	insensitive := []string{"PORT", "HOST", "LOG_LEVEL", "DEBUG", "APP_NAME"}
	for _, k := range insensitive {
		if IsSensitiveKey(k) {
			t.Errorf("expected %q to NOT be sensitive", k)
		}
	}
}

func TestRedactValueDefault(t *testing.T) {
	result := RedactValue("supersecret", nil)
	if result != "********" {
		t.Errorf("expected 8 asterisks, got %q", result)
	}
}

func TestRedactValueVisibleSuffix(t *testing.T) {
	opts := &RedactOptions{MaskChar: "*", VisibleSuffix: 4}
	result := RedactValue("supersecret", opts)
	if !strings.HasSuffix(result, "cret") {
		t.Errorf("expected suffix 'cret', got %q", result)
	}
	if !strings.HasPrefix(result, "********") {
		t.Errorf("expected 8 asterisk prefix, got %q", result)
	}
}

func TestRedactValueSuffixLargerThanValue(t *testing.T) {
	opts := &RedactOptions{MaskChar: "*", VisibleSuffix: 100}
	result := RedactValue("short", opts)
	if result != "********" {
		t.Errorf("expected full mask when suffix >= len, got %q", result)
	}
}

func TestRedact(t *testing.T) {
	env := map[string]string{
		"PORT":     "8080",
		"DB_PASSWORD": "s3cr3t",
		"API_KEY":  "abc123",
		"APP_NAME": "envsync",
	}
	redacted := Redact(env, nil)

	if redacted["PORT"] != "8080" {
		t.Errorf("PORT should be unchanged, got %q", redacted["PORT"])
	}
	if redacted["APP_NAME"] != "envsync" {
		t.Errorf("APP_NAME should be unchanged, got %q", redacted["APP_NAME"])
	}
	if redacted["DB_PASSWORD"] == "s3cr3t" {
		t.Error("DB_PASSWORD should be redacted")
	}
	if redacted["API_KEY"] == "abc123" {
		t.Error("API_KEY should be redacted")
	}
}

func TestRedactKeys(t *testing.T) {
	env := map[string]string{
		"PORT":    "8080",
		"HOST":    "localhost",
		"SPECIAL": "topsecret",
	}
	redacted := RedactKeys(env, []string{"SPECIAL", "PORT"}, nil)

	if redacted["HOST"] != "localhost" {
		t.Errorf("HOST should be unchanged, got %q", redacted["HOST"])
	}
	if redacted["PORT"] == "8080" {
		t.Error("PORT should be redacted")
	}
	if redacted["SPECIAL"] == "topsecret" {
		t.Error("SPECIAL should be redacted")
	}
}

func TestRedactDoesNotMutateInput(t *testing.T) {
	env := map[string]string{"DB_PASSWORD": "original"}
	_ = Redact(env, nil)
	if env["DB_PASSWORD"] != "original" {
		t.Error("Redact must not mutate the input map")
	}
}
