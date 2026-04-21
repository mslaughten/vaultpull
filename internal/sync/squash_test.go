package sync

import (
	"testing"
)

func TestSquashStrategyFromString_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  SquashStrategy
	}{
		{"concat", SquashStrategyConcat},
		{"first", SquashStrategyFirst},
		{"last", SquashStrategyLast},
		{"CONCAT", SquashStrategyConcat},
	}
	for _, tc := range cases {
		got, err := SquashStrategyFromString(tc.input)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("got %q, want %q", got, tc.want)
		}
	}
}

func TestSquashStrategyFromString_Invalid(t *testing.T) {
	_, err := SquashStrategyFromString("unknown")
	if err == nil {
		t.Fatal("expected error for unknown strategy"n
func TestNewSquasher_EmptyPrefix_ReturnsError(t *testing.T) {
	_, err := NewSquasher("", "OUT", ",", SquashStrategyConcat)
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestNewSquasher_EmptyOutKey_ReturnsError(t *testing.T) {
	_, err := NewSquasher("PREFIX_", "", ",", SquashStrategyConcat)
	if err == nil {
		t.Fatal("expected error for empty outKey")
	}
}

func TestSquasher_Concat_JoinsValues(t *testing.T) {
	s, _ := NewSquasher("TAG_", "TAGS", ",", SquashStrategyConcat)
	input := map[string]string{"TAG_A": "alpha", "TAG_B": "beta", "OTHER": "x"}
	out, err := s.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["TAGS"] != "alpha,beta" {
		t.Errorf("got %q, want %q", out["TAGS"], "alpha,beta")
	}
	if out["OTHER"] != "x" {
		t.Errorf("OTHER key should be preserved")
	}
	if _, ok := out["TAG_A"]; ok {
		t.Error("TAG_A should be removed")
	}
}

func TestSquasher_First_KeepsFirstValue(t *testing.T) {
	s, _ := NewSquasher("TAG_", "TAGS", ",", SquashStrategyFirst)
	input := map[string]string{"TAG_A": "alpha", "TAG_B": "beta"}
	out, _ := s.Apply(input)
	if out["TAGS"] != "alpha" {
		t.Errorf("got %q, want %q", out["TAGS"], "alpha")
	}
}

func TestSquasher_Last_KeepsLastValue(t *testing.T) {
	s, _ := NewSquasher("TAG_", "TAGS", ",", SquashStrategyLast)
	input := map[string]string{"TAG_A": "alpha", "TAG_B": "beta"}
	out, _ := s.Apply(input)
	if out["TAGS"] != "beta" {
		t.Errorf("got %q, want %q", out["TAGS"], "beta")
	}
}

func TestSquasher_NoMatchingKeys_ReturnsUnchanged(t *testing.T) {
	s, _ := NewSquasher("NOPE_", "OUT", ",", SquashStrategyConcat)
	input := map[string]string{"FOO": "bar"}
	out, _ := s.Apply(input)
	if out["FOO"] != "bar" {
		t.Error("expected unchanged map")
	}
	if _, ok := out["OUT"]; ok {
		t.Error("OUT key should not be created when no match")
	}
}
