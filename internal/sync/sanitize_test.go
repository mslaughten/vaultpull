package sync

import (
	"testing"
)

func TestNewSanitizer_DefaultRules(t *testing.T) {
	s, err := NewSanitizer(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.rules) != 1 {
		t.Fatalf("expected 1 default rule, got %d", len(s.rules))
	}
}

func TestNewSanitizer_ValidRules(t *testing.T) {
	s, err := NewSanitizer([]string{`-=_`, `\.=_`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(s.rules))
	}
}

func TestNewSanitizer_MissingSeparator(t *testing.T) {
	_, err := NewSanitizer([]string{"no-separator"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestNewSanitizer_InvalidPattern(t *testing.T) {
	_, err := NewSanitizer([]string{`[invalid=_`})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestSanitizer_Apply_DefaultRule(t *testing.T) {
	s, _ := NewSanitizer(nil)
	input := map[string]string{
		"my-key":    "val1",
		"foo.bar":   "val2",
		"GOOD_KEY":  "val3",
	}
	out := s.Apply(input)
	if out["my_key"] != "val1" {
		t.Errorf("expected my_key=val1, got %q", out["my_key"])
	}
	if out["foo_bar"] != "val2" {
		t.Errorf("expected foo_bar=val2, got %q", out["foo_bar"])
	}
	if out["GOOD_KEY"] != "val3" {
		t.Errorf("expected GOOD_KEY=val3, got %q", out["GOOD_KEY"])
	}
}

func TestSanitizer_Apply_CustomRule(t *testing.T) {
	s, err := NewSanitizer([]string{`-=DASH`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := s.Apply(map[string]string{"my-key": "v"})
	if out["myDASHkey"] != "v" {
		t.Errorf("expected myDASHkey, got keys: %v", out)
	}
}

func TestSanitizer_Apply_ValuesUnchanged(t *testing.T) {
	s, _ := NewSanitizer(nil)
	out := s.Apply(map[string]string{"k": "hello-world.value"})
	if out["k"] != "hello-world.value" {
		t.Errorf("value should not be modified, got %q", out["k"])
	}
}

func TestSanitizer_Apply_EmptyMap(t *testing.T) {
	s, _ := NewSanitizer(nil)
	out := s.Apply(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
