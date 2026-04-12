package sync

import (
	"testing"
)

func TestNewRedactor_ValidPatterns(t *testing.T) {
	r, err := NewRedactor([]string{`password`, `token`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil Redactor")
	}
}

func TestNewRedactor_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := NewRedactor([]string{`[invalid`})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestRedactor_Apply_MasksMatchingValues(t *testing.T) {
	r, _ := NewRedactor([]string{`(?i)secret`})
	secrets := map[string]string{
		"DB_PASSWORD": "mysecretpass",
		"APP_NAME":    "vaultpull",
	}
	out := r.Apply(secrets)
	if out["DB_PASSWORD"] != "***REDACTED***" {
		t.Errorf("expected redacted, got %q", out["DB_PASSWORD"])
	}
	if out["APP_NAME"] != "vaultpull" {
		t.Errorf("expected unchanged, got %q", out["APP_NAME"])
	}
}

func TestRedactor_Apply_NoRules_ReturnsUnchanged(t *testing.T) {
	r, _ := NewRedactor(nil)
	secrets := map[string]string{"KEY": "value"}
	out := r.Apply(secrets)
	if out["KEY"] != "value" {
		t.Errorf("expected value unchanged, got %q", out["KEY"])
	}
}

func TestRedactor_RedactString_ReplacesMatch(t *testing.T) {
	r, _ := NewRedactor([]string{`\d{4}-\d{4}-\d{4}-\d{4}`})
	input := "card: 1234-5678-9012-3456"
	got := r.RedactString(input)
	expected := "card: ***REDACTED***"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestRedactKeys_MasksKeysBySubstring(t *testing.T) {
	secrets := map[string]string{
		"DB_PASSWORD": "hunter2",
		"API_TOKEN":   "tok_abc",
		"APP_ENV":     "production",
	}
	out := RedactKeys(secrets, DefaultSensitivePatterns())
	if out["DB_PASSWORD"] != "***REDACTED***" {
		t.Errorf("expected DB_PASSWORD redacted, got %q", out["DB_PASSWORD"])
	}
	if out["API_TOKEN"] != "***REDACTED***" {
		t.Errorf("expected API_TOKEN redacted, got %q", out["API_TOKEN"])
	}
	if out["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV unchanged, got %q", out["APP_ENV"])
	}
}

func TestRedactKeys_EmptyPatterns_ReturnsUnchanged(t *testing.T) {
	secrets := map[string]string{"SECRET_KEY": "abc123"}
	out := RedactKeys(secrets, nil)
	if out["SECRET_KEY"] != "abc123" {
		t.Errorf("expected unchanged, got %q", out["SECRET_KEY"])
	}
}

func TestDefaultSensitivePatterns_NotEmpty(t *testing.T) {
	patterns := DefaultSensitivePatterns()
	if len(patterns) == 0 {
		t.Error("expected at least one default pattern")
	}
}
