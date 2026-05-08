package env

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseBasic(t *testing.T) {
	input := `
# comment
DB_HOST=localhost
DB_PORT=5432
APP_SECRET="mysecret"
`
	m, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expect := Map{
		"DB_HOST":    "localhost",
		"DB_PORT":    "5432",
		"APP_SECRET": "mysecret",
	}
	for k, v := range expect {
		if m[k] != v {
			t.Errorf("key %q: got %q, want %q", k, m[k], v)
		}
	}
}

func TestParseSingleQuotes(t *testing.T) {
	m, err := Parse(strings.NewReader("TOKEN='abc123'\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["TOKEN"] != "abc123" {
		t.Errorf("got %q, want %q", m["TOKEN"], "abc123")
	}
}

func TestParseInvalidLine(t *testing.T) {
	_, err := Parse(strings.NewReader("NOEQUALSIGN\n"))
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestParseEmptyKey(t *testing.T) {
	_, err := Parse(strings.NewReader("=value\n"))
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestParseEmpty(t *testing.T) {
	m, err := Parse(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m) != 0 {
		t.Errorf("expected empty map, got %v", m)
	}
}

func TestSerializeRoundtrip(t *testing.T) {
	original := Map{
		"Z_KEY": "last",
		"A_KEY": "first",
		"M_KEY": "middle",
	}
	var buf bytes.Buffer
	if err := Serialize(original, &buf); err != nil {
		t.Fatalf("serialize error: %v", err)
	}
	parsed, err := Parse(&buf)
	if err != nil {
		t.Fatalf("parse error after serialize: %v", err)
	}
	for k, v := range original {
		if parsed[k] != v {
			t.Errorf("key %q: got %q, want %q", k, parsed[k], v)
		}
	}
}
