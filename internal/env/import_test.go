package env

import (
	"testing"
)

func TestImportDotenvBasic(t *testing.T) {
	dst := map[string]string{"EXISTING": "yes"}
	result, err := Import(dst, "KEY=value\nFOO=bar\n", ImportOptions{Format: FormatDotenv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["KEY"] != "value" || result["FOO"] != "bar" || result["EXISTING"] != "yes" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestImportNoOverwrite(t *testing.T) {
	dst := map[string]string{"KEY": "original"}
	result, err := Import(dst, "KEY=new\n", ImportOptions{Format: FormatDotenv, Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["KEY"] != "original" {
		t.Errorf("expected original, got %q", result["KEY"])
	}
}

func TestImportOverwrite(t *testing.T) {
	dst := map[string]string{"KEY": "original"}
	result, err := Import(dst, "KEY=new\n", ImportOptions{Format: FormatDotenv, Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["KEY"] != "new" {
		t.Errorf("expected new, got %q", result["KEY"])
	}
}

func TestImportIgnoreKeys(t *testing.T) {
	dst := map[string]string{}
	result, err := Import(dst, "KEY=value\nSECRET=hidden\n", ImportOptions{
		Format:     FormatDotenv,
		IgnoreKeys: []string{"SECRET"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["SECRET"]; ok {
		t.Error("SECRET should have been ignored")
	}
	if result["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", result["KEY"])
	}
}

func TestImportShellFormat(t *testing.T) {
	dst := map[string]string{}
	input := "export DB_HOST=localhost\nexport DB_PORT=5432\n"
	result, err := Import(dst, input, ImportOptions{Format: FormatShell})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["DB_HOST"] != "localhost" || result["DB_PORT"] != "5432" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestImportShellFormatQuoted(t *testing.T) {
	dst := map[string]string{}
	input := `export APP_NAME="my app"` + "\n"
	result, err := Import(dst, input, ImportOptions{Format: FormatShell})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["APP_NAME"] != "my app" {
		t.Errorf("expected 'my app', got %q", result["APP_NAME"])
	}
}

func TestImportUnsupportedFormat(t *testing.T) {
	_, err := Import(map[string]string{}, "KEY=val", ImportOptions{Format: "yaml"})
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestImportDockerFormat(t *testing.T) {
	dst := map[string]string{}
	result, err := Import(dst, "PORT=8080\nHOST=0.0.0.0\n", ImportOptions{Format: FormatDocker})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["PORT"] != "8080" || result["HOST"] != "0.0.0.0" {
		t.Errorf("unexpected result: %v", result)
	}
}
