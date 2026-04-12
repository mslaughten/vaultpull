package sync

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestNewSchemaValidator_Valid(t *testing.T) {
	src := `[{"key":"DB_URL","required":true,"pattern":"^postgres://"},{"key":"PORT","required":false}]`
	v, err := NewSchemaValidator(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(v.rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(v.rules))
	}
}

func TestNewSchemaValidator_InvalidJSON(t *testing.T) {
	_, err := NewSchemaValidator(`not json`)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestNewSchemaValidator_MissingKey(t *testing.T) {
	_, err := NewSchemaValidator(`[{"required":true}]`)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestNewSchemaValidator_InvalidPattern(t *testing.T) {
	_, err := NewSchemaValidator(`[{"key":"X","pattern":"["}]`)
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestValidate_RequiredMissing(t *testing.T) {
	v, _ := NewSchemaValidator(`[{"key":"DB_URL","required":true}]`)
	violations := v.Validate(map[string]string{})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidate_PatternMismatch(t *testing.T) {
	v, _ := NewSchemaValidator(`[{"key":"DB_URL","pattern":"^postgres://"}]`)
	violations := v.Validate(map[string]string{"DB_URL": "mysql://host"})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidate_AllPassing(t *testing.T) {
	v, _ := NewSchemaValidator(`[{"key":"DB_URL","required":true,"pattern":"^postgres://"}]`)
	violations := v.Validate(map[string]string{"DB_URL": "postgres://localhost/db"})
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestWriteReport_Pass(t *testing.T) {
	v, _ := NewSchemaValidator(`[]`)
	var buf bytes.Buffer
	ok := v.WriteReport(&buf, nil)
	if !ok {
		t.Fatal("expected pass")
	}
}

func TestWriteReport_Fail(t *testing.T) {
	v, _ := NewSchemaValidator(`[]`)
	var buf bytes.Buffer
	ok := v.WriteReport(&buf, []string{"required key \"X\" is missing"})
	if ok {
		t.Fatal("expected failure")
	}
	if buf.Len() == 0 {
		t.Fatal("expected output")
	}
}

func TestLoadSchemaFile_Valid(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "schema.json")
	_ = os.WriteFile(p, []byte(`[{"key":"TOKEN","required":true}]`), 0600)
	v, err := LoadSchemaFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(v.rules) != 1 {
		t.Fatalf("expected 1 rule")
	}
}

func TestLoadSchemaFile_Missing(t *testing.T) {
	_, err := LoadSchemaFile("/nonexistent/schema.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
