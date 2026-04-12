package sync

import (
	"strings"
	"testing"
)

func TestNewLinter_DefaultRules(t *testing.T) {
	l := NewLinter(nil)
	if len(l.rules) == 0 {
		t.Fatal("expected default rules to be populated")
	}
}

func TestLinter_Check_ValidKeys(t *testing.T) {
	l := NewLinter(nil)
	secrets := map[string]string{
		"DATABASE_URL": "postgres://localhost",
		"API_KEY":      "abc123",
		"TIMEOUT_MS":   "3000",
	}
	violations := l.Check(secrets)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d: %v", len(violations), violations)
	}
}

func TestLinter_Check_LowercaseKey(t *testing.T) {
	l := NewLinter(nil)
	secrets := map[string]string{"db_password": "secret"}
	violations := l.Check(secrets)
	if len(violations) == 0 {
		t.Fatal("expected violation for lowercase key")
	}
	found := false
	for _, v := range violations {
		if v.Key == "db_password" && v.Rule == "uppercase" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected uppercase rule violation for db_password, got %v", violations)
	}
}

func TestLinter_Check_LeadingDigit(t *testing.T) {
	l := NewLinter(nil)
	secrets := map[string]string{"1_BAD_KEY": "value"}
	violations := l.Check(secrets)
	ruleNames := make(map[string]bool)
	for _, v := range violations {
		ruleNames[v.Rule] = true
	}
	if !ruleNames["no-leading-digit"] {
		t.Errorf("expected no-leading-digit violation, got rules: %v", ruleNames)
	}
}

func TestLinter_Check_DoubleUnderscore(t *testing.T) {
	l := NewLinter(nil)
	secrets := map[string]string{"BAD__KEY": "value"}
	violations := l.Check(secrets)
	ruleNames := make(map[string]bool)
	for _, v := range violations {
		ruleNames[v.Rule] = true
	}
	if !ruleNames["no-double-underscore"] {
		t.Errorf("expected no-double-underscore violation, got rules: %v", ruleNames)
	}
}

func TestLinter_Summary_NoViolations(t *testing.T) {
	l := NewLinter(nil)
	summary := l.Summary(nil)
	if summary != "lint: all keys passed" {
		t.Errorf("unexpected summary: %q", summary)
	}
}

func TestLinter_Summary_WithViolations(t *testing.T) {
	l := NewLinter(nil)
	violations := []LintViolation{
		{Key: "bad_key", Rule: "uppercase", Message: "must be uppercase"},
	}
	summary := l.Summary(violations)
	if !strings.Contains(summary, "1 violation") {
		t.Errorf("expected violation count in summary, got: %q", summary)
	}
	if !strings.Contains(summary, "bad_key") {
		t.Errorf("expected key name in summary, got: %q", summary)
	}
}

func TestLintViolation_Error(t *testing.T) {
	v := LintViolation{Key: "foo", Rule: "uppercase", Message: "must be uppercase"}
	err := v.Error()
	if !strings.Contains(err, "foo") || !strings.Contains(err, "uppercase") {
		t.Errorf("unexpected error string: %q", err)
	}
}
