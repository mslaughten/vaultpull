package sync

import (
	"testing"
)

func TestNewMasker_InvalidMode(t *testing.T) {
	_, err := NewMasker(MaskMode(99))
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestNewMasker_NegativeReveal(t *testing.T) {
	_, err := NewMasker(MaskPartial, WithRevealChars(-1))
	if err == nil {
		t.Fatal("expected error for negative reveal chars")
	}
}

func TestMasker_MaskFull(t *testing.T) {
	m, err := NewMasker(MaskFull)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := m.Apply("supersecret")
	if got != "***********" {
		t.Errorf("expected 11 asterisks, got %q", got)
	}
}

func TestMasker_MaskFull_EmptyValue(t *testing.T) {
	m, _ := NewMasker(MaskFull)
	got := m.Apply("")
	if got != "*" {
		t.Errorf("expected single asterisk for empty value, got %q", got)
	}
}

func TestMasker_MaskNone(t *testing.T) {
	m, _ := NewMasker(MaskNone)
	got := m.Apply("plaintext")
	if got != "plaintext" {
		t.Errorf("expected unchanged value, got %q", got)
	}
}

func TestMasker_MaskPartial_RevealsTrailing(t *testing.T) {
	m, _ := NewMasker(MaskPartial, WithRevealChars(4))
	got := m.Apply("supersecret")
	// "supersecret" => 11 chars, reveal last 4 = "cret", mask 7 = "*******cret"
	if got != "*******cret" {
		t.Errorf("got %q", got)
	}
}

func TestMasker_MaskPartial_ShortValue(t *testing.T) {
	m, _ := NewMasker(MaskPartial, WithRevealChars(4))
	// value shorter than reveal => fully masked
	got := m.Apply("abc")
	if got != "***" {
		t.Errorf("expected full mask for short value, got %q", got)
	}
}

func TestMasker_CustomSymbol(t *testing.T) {
	m, _ := NewMasker(MaskFull, WithMaskSymbol("#"))
	got := m.Apply("hi")
	if got != "##" {
		t.Errorf("expected ##, got %q", got)
	}
}

func TestMasker_FixedLength(t *testing.T) {
	m, _ := NewMasker(MaskFull, WithFixedLength(6))
	got := m.Apply("anyvalue")
	if got != "******" {
		t.Errorf("expected 6 asterisks, got %q", got)
	}
}

func TestMasker_ApplyMap(t *testing.T) {
	m, _ := NewMasker(MaskFull)
	secrets := map[string]string{
		"KEY_A": "value1",
		"KEY_B": "value2",
	}
	masked := m.ApplyMap(secrets)
	for k, v := range masked {
		if v == secrets[k] {
			t.Errorf("key %s was not masked", k)
		}
	}
	if len(masked) != len(secrets) {
		t.Errorf("expected %d entries, got %d", len(secrets), len(masked))
	}
}
