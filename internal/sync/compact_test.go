package sync

import (
	"testing"
)

func TestCompactStrategyFromString_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  CompactStrategy
	}{
		{"empty", CompactEmpty},
		{"blank", CompactBlank},
	}
	for _, tc := range cases {
		got, err := CompactStrategyFromString(tc.input)
		if err != nil {
			t.Fatalf("%q: unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("%q: got %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestCompactStrategyFromString_Invalid(t *testing.T) {
	_, err := CompactStrategyFromString("unknown")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestCompacter_Apply_RemovesEmptyValues(t *testing.T) {
	c, err := NewCompacter("empty")
	if err != nil {
		t.Fatal(err)
	}
	input := map[string]string{
		"KEY_A": "value",
		"KEY_B": "",
		"KEY_C": "  ",
	}
	out := c.Apply(input)
	if _, ok := out["KEY_B"]; ok {
		t.Error("KEY_B should have been removed")
	}
	if _, ok := out["KEY_A"]; !ok {
		t.Error("KEY_A should be retained")
	}
	// blank-but-not-empty should be kept under 'empty' strategy
	if _, ok := out["KEY_C"]; !ok {
		t.Error("KEY_C should be retained under empty strategy")
	}
}

func TestCompacter_Apply_RemovesBlankValues(t *testing.T) {
	c, err := NewCompacter("blank")
	if err != nil {
		t.Fatal(err)
	}
	input := map[string]string{
		"KEY_A": "value",
		"KEY_B": "",
		"KEY_C": "   ",
		"KEY_D": "\t\n",
	}
	out := c.Apply(input)
	for _, k := range []string{"KEY_B", "KEY_C", "KEY_D"} {
		if _, ok := out[k]; ok {
			t.Errorf("%s should have been removed", k)
		}
	}
	if _, ok := out["KEY_A"]; !ok {
		t.Error("KEY_A should be retained")
	}
}

func TestCompacter_Apply_EmptyMap_ReturnsEmpty(t *testing.T) {
	c, _ := NewCompacter("empty")
	out := c.Apply(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestNewCompacter_InvalidStrategy_ReturnsError(t *testing.T) {
	_, err := NewCompacter("nope")
	if err == nil {
		t.Fatal("expected error for invalid strategy")
	}
}
