package env

import (
	"testing"
)

func TestLintNoViolations(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"APP_PORT":     "8080",
	}
	violations := Lint(env, nil)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d: %v", len(violations), violations)
	}
}

func TestLintUppercaseKey(t *testing.T) {
	env := map[string]string{
		"database_url": "postgres://localhost/db",
		"VALID_KEY":    "value",
	}
	violations := Lint(env, nil)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Rule != "uppercase-key" {
		t.Errorf("expected rule uppercase-key, got %q", violations[0].Rule)
	}
	if violations[0].Key != "database_url" {
		t.Errorf("expected key database_url, got %q", violations[0].Key)
	}
}

func TestLintTrailingSpace(t *testing.T) {
	env := map[string]string{
		"APP_HOST": "localhost  ",
	}
	violations := Lint(env, nil)
	found := false
	for _, v := range violations {
		if v.Rule == "no-trailing-space" {
			found = true
		}
	}
	if !found {
		t.Error("expected no-trailing-space violation")
	}
}

func TestLintEmptyValue(t *testing.T) {
	env := map[string]string{
		"EMPTY_KEY": "",
	}
	violations := Lint(env, nil)
	found := false
	for _, v := range violations {
		if v.Rule == "no-empty-value" {
			found = true
		}
	}
	if !found {
		t.Error("expected no-empty-value violation")
	}
}

func TestLintCustomRule(t *testing.T) {
	customRules := []LintRule{
		{
			Name:    "no-localhost",
			Message: "value must not reference localhost in production",
			Check: func(_, value string) bool {
				return value == "localhost"
			},
		},
	}
	env := map[string]string{
		"DB_HOST": "localhost",
		"API_HOST": "api.example.com",
	}
	violations := Lint(env, customRules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Rule != "no-localhost" {
		t.Errorf("expected rule no-localhost, got %q", violations[0].Rule)
	}
}

func TestLintViolationError(t *testing.T) {
	v := LintViolation{Key: "foo", Rule: "bar", Message: "baz"}
	got := v.Error()
	if got == "" {
		t.Error("expected non-empty error string")
	}
}
