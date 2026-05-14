package env

import (
	"testing"
)

func TestParseTemplateBasic(t *testing.T) {
	input := `
# required keys
DB_HOST
DB_PORT=5432
APP_ENV=development
`
	tmpl, err := ParseTemplate(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tmpl["DB_HOST"] != "" {
		t.Errorf("expected empty default for DB_HOST")
	}
	if tmpl["DB_PORT"] != "5432" {
		t.Errorf("expected default 5432 for DB_PORT, got %q", tmpl["DB_PORT"])
	}
	if tmpl["APP_ENV"] != "development" {
		t.Errorf("expected default development for APP_ENV")
	}
}

func TestParseTemplateEmptyKey(t *testing.T) {
	_, err := ParseTemplate("=value")
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestRenderTemplateAllPresent(t *testing.T) {
	tmpl := map[string]string{"HOST": "", "PORT": "3000"}
	env := map[string]string{"HOST": "localhost", "PORT": "8080"}

	res := RenderTemplate(tmpl, env)
	if len(res.Missing) != 0 {
		t.Errorf("expected no missing keys, got %v", res.Missing)
	}
	if res.Rendered["HOST"] != "localhost" {
		t.Errorf("expected localhost")
	}
	if res.Rendered["PORT"] != "8080" {
		t.Errorf("expected 8080")
	}
}

func TestRenderTemplateMissingRequired(t *testing.T) {
	tmpl := map[string]string{"SECRET": "", "PORT": "3000"}
	env := map[string]string{"PORT": "9000"}

	res := RenderTemplate(tmpl, env)
	if len(res.Missing) != 1 || res.Missing[0] != "SECRET" {
		t.Errorf("expected SECRET in missing, got %v", res.Missing)
	}
}

func TestRenderTemplateUsesDefault(t *testing.T) {
	tmpl := map[string]string{"PORT": "3000"}
	env := map[string]string{}

	res := RenderTemplate(tmpl, env)
	if len(res.Missing) != 0 {
		t.Errorf("expected no missing keys")
	}
	if res.Rendered["PORT"] != "3000" {
		t.Errorf("expected default 3000, got %q", res.Rendered["PORT"])
	}
}

func TestRenderTemplateUnusedKeys(t *testing.T) {
	tmpl := map[string]string{"HOST": ""}
	env := map[string]string{"HOST": "localhost", "EXTRA": "value", "ANOTHER": "x"}

	res := RenderTemplate(tmpl, env)
	if len(res.Unused) != 2 {
		t.Errorf("expected 2 unused keys, got %v", res.Unused)
	}
	if res.Unused[0] != "ANOTHER" || res.Unused[1] != "EXTRA" {
		t.Errorf("unexpected unused order: %v", res.Unused)
	}
}

func TestRenderTemplateEmpty(t *testing.T) {
	res := RenderTemplate(map[string]string{}, map[string]string{})
	if len(res.Missing) != 0 || len(res.Unused) != 0 || len(res.Rendered) != 0 {
		t.Error("expected all empty results")
	}
}
