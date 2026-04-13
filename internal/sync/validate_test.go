package sync

import (
	"strings"
	"testing"
)

func TestNewValidator_ValidRules(t *testing.T) {
	v, err := NewValidator([]string{"DB_HOST=.+", "PORT=^[0-9]+$"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(v.rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(v.rules))
	}
}

func TestNewValidator_MissingSeparator(t *testing.T) {
	_, err := NewValidator([]string{"NOEQUALS"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestNewValidator_EmptyKey(t *testing.T) {
	_, err := NewValidator([]string{"=pattern"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestNewValidator_InvalidPattern(t *testing.T) {
	_, err := NewValidator([]string{"KEY=[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestValidator_Valid_PassesAllRules(t *testing.T) {
	v, _ := NewValidator([]string{"HOST=.+", "PORT=^[0-9]+$"})
	err := v.Validate(map[string]string{"HOST": "localhost", "PORT": "5432"})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidator_MissingRequiredKey(t *testing.T) {
	v, _ := NewValidator([]string{"TOKEN=.+"})
	err := v.Validate(map[string]string{})
	if err == nil {
		t.Fatal("expected validation error for missing key")
	}
	if !strings.Contains(err.Error(), "TOKEN") {
		t.Errorf("expected error to mention TOKEN, got: %v", err)
	}
}

func TestValidator_PatternMismatch(t *testing.T) {
	v, _ := NewValidator([]string{"PORT=^[0-9]+$"})
	err := v.Validate(map[string]string{"PORT": "not-a-number"})
	if err == nil {
		t.Fatal("expected validation error for pattern mismatch")
	}
	if !strings.Contains(err.Error(), "PORT") {
		t.Errorf("expected error to mention PORT, got: %v", err)
	}
}

func TestValidator_EmptyPatternRequiresPresence(t *testing.T) {
	v, _ := NewValidator([]string{"SECRET="})
	if err := v.Validate(map[string]string{"SECRET": "value"}); err != nil {
		t.Errorf("expected no error when key present, got: %v", err)
	}
	if err := v.Validate(map[string]string{}); err == nil {
		t.Error("expected error when key missing")
	}
}

func TestValidationError_CollectsMultiple(t *testing.T) {
	v, _ := NewValidator([]string{"A=.+", "B=.+"})
	err := v.Validate(map[string]string{})
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Violations) != 2 {
		t.Errorf("expected 2 violations, got %d", len(ve.Violations))
	}
}
