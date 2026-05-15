package env

import (
	"strings"
	"testing"
)

func TestParseSchemaEmptyKey(t *testing.T) {
	_, err := ParseSchema([]SchemaField{{Key: "", Required: true}})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestValidateAllPresent(t *testing.T) {
	schema, err := ParseSchema([]SchemaField{
		{Key: "APP_ENV", Required: true},
		{Key: "PORT", Required: true, Pattern: `^\d+$`},
	})
	if err != nil {
		t.Fatalf("ParseSchema: %v", err)
	}

	env := map[string]string{"APP_ENV": "production", "PORT": "8080"}
	if err := schema.Validate(env); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidateMissingRequired(t *testing.T) {
	schema, _ := ParseSchema([]SchemaField{
		{Key: "DATABASE_URL", Required: true},
	})

	err := schema.Validate(map[string]string{})
	if err == nil {
		t.Fatal("expected validation error")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(ve.Violations))
	}
}

func TestValidatePatternMismatch(t *testing.T) {
	schema, _ := ParseSchema([]SchemaField{
		{Key: "PORT", Required: true, Pattern: `^\d+$`},
	})

	err := schema.Validate(map[string]string{"PORT": "not-a-number"})
	if err == nil {
		t.Fatal("expected pattern violation")
	}
	ve := err.(*ValidationError)
	if !strings.Contains(ve.Violations[0], "PORT") {
		t.Errorf("violation should mention key PORT: %s", ve.Violations[0])
	}
}

func TestValidateOptionalMissing(t *testing.T) {
	schema, _ := ParseSchema([]SchemaField{
		{Key: "LOG_LEVEL", Required: false},
	})

	if err := schema.Validate(map[string]string{}); err != nil {
		t.Fatalf("optional missing key should not fail: %v", err)
	}
}

func TestValidateMultipleViolations(t *testing.T) {
	schema, _ := ParseSchema([]SchemaField{
		{Key: "A", Required: true},
		{Key: "B", Required: true},
		{Key: "C", Required: true, Pattern: `^yes|no$`},
	})

	err := schema.Validate(map[string]string{"C": "maybe"})
	if err == nil {
		t.Fatal("expected violations")
	}
	ve := err.(*ValidationError)
	if len(ve.Violations) != 3 {
		t.Fatalf("expected 3 violations, got %d: %v", len(ve.Violations), ve.Violations)
	}
}
