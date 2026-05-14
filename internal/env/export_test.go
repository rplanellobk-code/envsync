package env

import (
	"testing"
)

func TestExportDotenvFormat(t *testing.T) {
	vars := map[string]string{
		"APP_ENV": "production",
		"DB_URL":  "postgres://localhost/db",
	}
	out := Export(vars, ExportOptions{Format: FormatDotenv})
	want := "APP_ENV=production\nDB_URL=postgres://localhost/db"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestExportShellFormat(t *testing.T) {
	vars := map[string]string{
		"FOO": "bar",
	}
	out := Export(vars, ExportOptions{Format: FormatExport})
	want := "export FOO=bar"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestExportDockerFormat(t *testing.T) {
	vars := map[string]string{
		"PORT": "8080",
	}
	out := Export(vars, ExportOptions{Format: FormatDocker})
	want := "--env PORT=8080"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestExportOmitKeys(t *testing.T) {
	vars := map[string]string{
		"SECRET": "topsecret",
		"PUBLIC": "visible",
	}
	out := Export(vars, ExportOptions{
		Format:   FormatDotenv,
		OmitKeys: []string{"SECRET"},
	})
	want := "PUBLIC=visible"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestExportQuotesSpaces(t *testing.T) {
	vars := map[string]string{
		"GREETING": "hello world",
	}
	out := Export(vars, ExportOptions{Format: FormatDotenv})
	want := `GREETING="hello world"`
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestExportEmptyValue(t *testing.T) {
	vars := map[string]string{
		"EMPTY": "",
	}
	out := Export(vars, ExportOptions{Format: FormatDotenv})
	want := `EMPTY=""`
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestExportDeterministicOrder(t *testing.T) {
	vars := map[string]string{
		"Z_KEY": "z",
		"A_KEY": "a",
		"M_KEY": "m",
	}
	out := Export(vars, ExportOptions{Format: FormatDotenv})
	want := "A_KEY=a\nM_KEY=m\nZ_KEY=z"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}
