package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	return p
}

func TestRead_ParsesKeyValues(t *testing.T) {
	p := writeTemp(t, "FOO=bar\nBAZ=qux\n")
	r := NewReader(p)
	got, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["FOO"] != "bar" || got["BAZ"] != "qux" {
		t.Errorf("unexpected map: %v", got)
	}
}

func TestRead_SkipsComments(t *testing.T) {
	p := writeTemp(t, "# comment\nKEY=value\n")
	r := NewReader(p)
	got, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got["KEY"] != "value" {
		t.Errorf("unexpected map: %v", got)
	}
}

func TestRead_StripsQuotes(t *testing.T) {
	p := writeTemp(t, `SINGLE='hello world'
DOUBLE="goodbye world"
`)
	r := NewReader(p)
	got, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["SINGLE"] != "hello world" {
		t.Errorf("SINGLE: got %q", got["SINGLE"])
	}
	if got["DOUBLE"] != "goodbye world" {
		t.Errorf("DOUBLE: got %q", got["DOUBLE"])
	}
}

func TestRead_MissingFile_ReturnsEmpty(t *testing.T) {
	r := NewReader("/nonexistent/path/.env")
	got, err := r.Read()
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got: %v", got)
	}
}

func TestRead_SkipsLinesWithoutEquals(t *testing.T) {
	p := writeTemp(t, "NOEQUALS\nVALID=yes\n")
	r := NewReader(p)
	got, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := got["NOEQUALS"]; ok {
		t.Error("expected NOEQUALS to be skipped")
	}
	if got["VALID"] != "yes" {
		t.Errorf("VALID: got %q", got["VALID"])
	}
}
