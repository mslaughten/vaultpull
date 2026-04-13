package sync

import (
	"testing"
)

func TestNewTypeCaster_ValidRules(t *testing.T) {
	tc, err := NewTypeCaster([]string{"PORT=int", "RATIO=float", "ENABLED=bool", "NAME=string"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tc.rules) != 4 {
		t.Fatalf("expected 4 rules, got %d", len(tc.rules))
	}
}

func TestNewTypeCaster_InvalidRule_MissingEquals(t *testing.T) {
	_, err := NewTypeCaster([]string{"PORT"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestNewTypeCaster_UnknownType(t *testing.T) {
	_, err := NewTypeCaster([]string{"PORT=integer"})
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestNewTypeCaster_EmptyKey_ReturnsError(t *testing.T) {
	_, err := NewTypeCaster([]string{"=int"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestTypeCaster_Apply_Int_Valid(t *testing.T) {
	tc, _ := NewTypeCaster([]string{"PORT=int"})
	out, err := tc.Apply(map[string]string{"PORT": "8080"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected '8080', got %q", out["PORT"])
	}
}

func TestTypeCaster_Apply_Int_Invalid(t *testing.T) {
	tc, _ := NewTypeCaster([]string{"PORT=int"})
	_, err := tc.Apply(map[string]string{"PORT": "not-a-number"})
	if err == nil {
		t.Fatal("expected error for invalid int")
	}
}

func TestTypeCaster_Apply_Bool_Normalises(t *testing.T) {
	tc, _ := NewTypeCaster([]string{"ENABLED=bool"})
	out, err := tc.Apply(map[string]string{"ENABLED": "1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["ENABLED"] != "true" {
		t.Errorf("expected 'true', got %q", out["ENABLED"])
	}
}

func TestTypeCaster_Apply_Float_Normalises(t *testing.T) {
	tc, _ := NewTypeCaster([]string{"RATIO=float"})
	out, err := tc.Apply(map[string]string{"RATIO": "3.14000"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["RATIO"] != "3.14" {
		t.Errorf("expected '3.14', got %q", out["RATIO"])
	}
}

func TestTypeCaster_Apply_MissingKey_IsSkipped(t *testing.T) {
	tc, _ := NewTypeCaster([]string{"MISSING=int"})
	out, err := tc.Apply(map[string]string{"OTHER": "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["OTHER"] != "hello" {
		t.Errorf("expected 'hello', got %q", out["OTHER"])
	}
}

func TestTypeCaster_Apply_DoesNotMutateInput(t *testing.T) {
	tc, _ := NewTypeCaster([]string{"FLAG=bool"})
	input := map[string]string{"FLAG": "0"}
	_, _ = tc.Apply(input)
	if input["FLAG"] != "0" {
		t.Error("Apply must not mutate the input map")
	}
}
