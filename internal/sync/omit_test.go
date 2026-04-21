package sync

import (
	"testing"
)

func TestNewOmitter_ValidPatterns(t *testing.T) {
	o, err := NewOmitter([]string{`^secret`, `password`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(o.patterns) != 2 {
		t.Fatalf("expected 2 patterns, got %d", len(o.patterns))
	}
}

func TestNewOmitter_EmptyPattern_ReturnsError(t *testing.T) {
	_, err := NewOmitter([]string{""})
	if err == nil {
		t.Fatal("expected error for empty pattern, got nil")
	}
}

func TestNewOmitter_InvalidRegex_ReturnsError(t *testing.T) {
	_, err := NewOmitter([]string{`[invalid`})
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestNewOmitter_NoPatterns_IsNoop(t *testing.T) {
	o, err := NewOmitter(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := map[string]string{"KEY": "value", "OTHER": "data"}
	out := o.Apply(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestOmitter_Apply_RemovesMatchingValues(t *testing.T) {
	o, _ := NewOmitter([]string{`^REDACTED`, `\*\*\*`})
	input := map[string]string{
		"API_KEY":  "REDACTED-abc123",
		"DB_PASS":  "***hidden***",
		"APP_NAME": "vaultpull",
	}
	out := o.Apply(input)
	if _, ok := out["API_KEY"]; ok {
		t.Error("API_KEY should have been omitted")
	}
	if _, ok := out["DB_PASS"]; ok {
		t.Error("DB_PASS should have been omitted")
	}
	if v, ok := out["APP_NAME"]; !ok || v != "vaultpull" {
		t.Errorf("APP_NAME should be kept with value 'vaultpull', got %q", v)
	}
}

func TestOmitter_Apply_EmptyMap_ReturnsEmpty(t *testing.T) {
	o, _ := NewOmitter([]string{`.*`})
	out := o.Apply(map[string]string{})
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(out))
	}
}

func TestOmitter_Apply_NoMatchesKeepsAll(t *testing.T) {
	o, _ := NewOmitter([]string{`^NEVER_MATCHES_XYZ$`})
	input := map[string]string{"A": "1", "B": "2"}
	out := o.Apply(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}
