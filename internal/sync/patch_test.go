package sync

import (
	"testing"
)

func TestNewPatcher_ValidRules(t *testing.T) {
	p, err := NewPatcher([]string{"set:FOO=bar", "delete:BAZ", "append:MSG= world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.ops) != 3 {
		t.Fatalf("expected 3 ops, got %d", len(p.ops))
	}
}

func TestNewPatcher_InvalidFormat_MissingColon(t *testing.T) {
	_, err := NewPatcher([]string{"setFOO=bar"})
	if err == nil {
		t.Fatal("expected error for missing colon")
	}
}

func TestNewPatcher_UnknownOp(t *testing.T) {
	_, err := NewPatcher([]string{"upsert:FOO=bar"})
	if err == nil {
		t.Fatal("expected error for unknown op")
	}
}

func TestNewPatcher_SetMissingEquals(t *testing.T) {
	_, err := NewPatcher([]string{"set:FOO"})
	if err == nil {
		t.Fatal("expected error for missing =value")
	}
}

func TestNewPatcher_EmptyKey(t *testing.T) {
	_, err := NewPatcher([]string{"delete:"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestPatcher_Apply_Set(t *testing.T) {
	p, _ := NewPatcher([]string{"set:FOO=newval"})
	out, err := p.Apply(map[string]string{"FOO": "oldval", "BAR": "keep"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "newval" {
		t.Errorf("expected FOO=newval, got %q", out["FOO"])
	}
	if out["BAR"] != "keep" {
		t.Errorf("expected BAR=keep, got %q", out["BAR"])
	}
}

func TestPatcher_Apply_Delete(t *testing.T) {
	p, _ := NewPatcher([]string{"delete:REMOVE_ME"})
	out, err := p.Apply(map[string]string{"REMOVE_ME": "gone", "KEEP": "yes"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["REMOVE_ME"]; ok {
		t.Error("expected REMOVE_ME to be deleted")
	}
	if out["KEEP"] != "yes" {
		t.Errorf("expected KEEP=yes, got %q", out["KEEP"])
	}
}

func TestPatcher_Apply_Append(t *testing.T) {
	p, _ := NewPatcher([]string{"append:GREETING= world"})
	out, err := p.Apply(map[string]string{"GREETING": "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["GREETING"] != "hello world" {
		t.Errorf("expected 'hello world', got %q", out["GREETING"])
	}
}

func TestPatcher_Apply_Append_MissingKey(t *testing.T) {
	p, _ := NewPatcher([]string{"append:NEW=value"})
	out, err := p.Apply(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["NEW"] != "value" {
		t.Errorf("expected NEW=value, got %q", out["NEW"])
	}
}

func TestPatcher_Apply_DoesNotMutateInput(t *testing.T) {
	p, _ := NewPatcher([]string{"set:A=changed"})
	in := map[string]string{"A": "original"}
	p.Apply(in) //nolint:errcheck
	if in["A"] != "original" {
		t.Error("input map was mutated")
	}
}
