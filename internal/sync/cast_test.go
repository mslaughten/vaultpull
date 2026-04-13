package sync

import (
	"testing"
)

func TestNewCastFormatter_ValidRules(t *testing.T) {
	c, err := NewCastFormatter([]string{"FOO=upper", "BAR=lower", "BAZ=title"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(c.rules) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(c.rules))
	}
}

func TestNewCastFormatter_InvalidRule_MissingSeparator(t *testing.T) {
	_, err := NewCastFormatter([]string{"FOOupper"})
	if err == nil {
		t.Fatal("expected error for missing separator")
	}
}

func TestNewCastFormatter_EmptyKey_ReturnsError(t *testing.T) {
	_, err := NewCastFormatter([]string{"=upper"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestNewCastFormatter_UnknownFormat_ReturnsError(t *testing.T) {
	_, err := NewCastFormatter([]string{"FOO=camel"})
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestCastFormatter_Apply_Upper(t *testing.T) {
	c, _ := NewCastFormatter([]string{"FOO=upper"})
	out, err := c.Apply(map[string]string{"FOO": "hello", "BAR": "world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "HELLO" {
		t.Errorf("expected HELLO, got %s", out["FOO"])
	}
	if out["BAR"] != "world" {
		t.Errorf("expected world unchanged, got %s", out["BAR"])
	}
}

func TestCastFormatter_Apply_Lower(t *testing.T) {
	c, _ := NewCastFormatter([]string{"FOO=lower"})
	out, _ := c.Apply(map[string]string{"FOO": "HELLO"})
	if out["FOO"] != "hello" {
		t.Errorf("expected hello, got %s", out["FOO"])
	}
}

func TestCastFormatter_Apply_Title(t *testing.T) {
	c, _ := NewCastFormatter([]string{"FOO=title"})
	out, _ := c.Apply(map[string]string{"FOO": "hello world"})
	if out["FOO"] != "Hello World" {
		t.Errorf("expected 'Hello World', got %s", out["FOO"])
	}
}

func TestCastFormatter_Apply_MissingKey_Skips(t *testing.T) {
	c, _ := NewCastFormatter([]string{"MISSING=upper"})
	out, err := c.Apply(map[string]string{"FOO": "bar"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" {
		t.Errorf("expected bar unchanged, got %s", out["FOO"])
	}
}
