package sync

import (
	"testing"
)

func TestNewTransformer_ValidRules(t *testing.T) {
	tr, err := NewTransformer([]string{"DB_PASS=upper", "API_KEY=trimspace"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil transformer")
	}
}

func TestNewTransformer_InvalidRule_MissingEquals(t *testing.T) {
	_, err := NewTransformer([]string{"BADFORMAT"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestNewTransformer_InvalidRule_UnknownTransform(t *testing.T) {
	_, err := NewTransformer([]string{"KEY=rot13"})
	if err == nil {
		t.Fatal("expected error for unknown transform")
	}
}

func TestNewTransformer_EmptyKey_ReturnsError(t *testing.T) {
	_, err := NewTransformer([]string{"=upper"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestTransformer_Apply_Upper(t *testing.T) {
	tr, _ := NewTransformer([]string{"SECRET=upper"})
	out, err := tr.Apply(map[string]string{"SECRET": "hello", "OTHER": "world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["SECRET"] != "HELLO" {
		t.Errorf("expected HELLO, got %q", out["SECRET"])
	}
	if out["OTHER"] != "world" {
		t.Errorf("expected world unchanged, got %q", out["OTHER"])
	}
}

func TestTransformer_Apply_Lower(t *testing.T) {
	tr, _ := NewTransformer([]string{"TOKEN=lower"})
	out, err := tr.Apply(map[string]string{"TOKEN": "ABCXYZ"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["TOKEN"] != "abcxyz" {
		t.Errorf("expected abcxyz, got %q", out["TOKEN"])
	}
}

func TestTransformer_Apply_TrimSpace(t *testing.T) {
	tr, _ := NewTransformer([]string{"VAL=trimspace"})
	out, err := tr.Apply(map[string]string{"VAL": "  spaced  "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["VAL"] != "spaced" {
		t.Errorf("expected 'spaced', got %q", out["VAL"])
	}
}

func TestTransformer_Apply_MissingKey_Skips(t *testing.T) {
	tr, _ := NewTransformer([]string{"MISSING=upper"})
	out, err := tr.Apply(map[string]string{"PRESENT": "value"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["MISSING"]; ok {
		t.Error("expected MISSING key to not be added")
	}
}

func TestTransformer_Apply_NoRules_ReturnsOriginal(t *testing.T) {
	tr, _ := NewTransformer(nil)
	input := map[string]string{"A": "1", "B": "2"}
	out, err := tr.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "1" || out["B"] != "2" {
		t.Error("expected original values to be preserved")
	}
}
