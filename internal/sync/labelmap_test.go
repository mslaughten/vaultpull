package sync

import (
	"testing"
)

func TestNewLabelMapper_Valid(t *testing.T) {
	lm, err := NewLabelMapper([]string{"OLD_KEY=NEW_KEY", "FOO=BAR"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lm.Len() != 2 {
		t.Fatalf("expected 2 rules, got %d", lm.Len())
	}
}

func TestNewLabelMapper_NoRules_IsNoop(t *testing.T) {
	lm, err := NewLabelMapper(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lm.Len() != 0 {
		t.Fatalf("expected 0 rules, got %d", lm.Len())
	}
}

func TestNewLabelMapper_MissingSeparator_ReturnsError(t *testing.T) {
	_, err := NewLabelMapper([]string{"BADFORMAT"})
	if err == nil {
		t.Fatal("expected error for missing separator")
	}
}

func TestNewLabelMapper_EmptyLabel_ReturnsError(t *testing.T) {
	_, err := NewLabelMapper([]string{"=NEW_KEY"})
	if err == nil {
		t.Fatal("expected error for empty label")
	}
}

func TestNewLabelMapper_EmptyTargetKey_ReturnsError(t *testing.T) {
	_, err := NewLabelMapper([]string{"OLD_KEY="})
	if err == nil {
		t.Fatal("expected error for empty target key")
	}
}

func TestLabelMapper_Apply_RenamesMatchingKeys(t *testing.T) {
	lm, _ := NewLabelMapper([]string{"DB_PASS=DATABASE_PASSWORD", "API=API_KEY"})
	input := map[string]string{
		"DB_PASS": "secret123",
		"API":     "tok_abc",
		"OTHER":   "unchanged",
	}
	out, err := lm.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DATABASE_PASSWORD"] != "secret123" {
		t.Errorf("expected DATABASE_PASSWORD=secret123, got %q", out["DATABASE_PASSWORD"])
	}
	if out["API_KEY"] != "tok_abc" {
		t.Errorf("expected API_KEY=tok_abc, got %q", out["API_KEY"])
	}
	if out["OTHER"] != "unchanged" {
		t.Errorf("expected OTHER=unchanged, got %q", out["OTHER"])
	}
	if _, exists := out["DB_PASS"]; exists {
		t.Error("old key DB_PASS should not exist in output")
	}
}

func TestLabelMapper_Apply_NoMatchPassesThrough(t *testing.T) {
	lm, _ := NewLabelMapper([]string{"X=Y"})
	input := map[string]string{"A": "1", "B": "2"}
	out, err := lm.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestLabelMapper_Apply_EmptyMap_ReturnsEmpty(t *testing.T) {
	lm, _ := NewLabelMapper([]string{"A=B"})
	out, err := lm.Apply(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
