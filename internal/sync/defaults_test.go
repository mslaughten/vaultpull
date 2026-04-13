package sync

import (
	"testing"
)

func TestDefaultsStrategyFromString_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  DefaultsStrategy
	}{
		{"", DefaultsSkip},
		{"skip", DefaultsSkip},
		{"apply", DefaultsApply},
	}
	for _, tc := range cases {
		got, err := DefaultsStrategyFromString(tc.input)
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("input %q: got %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestDefaultsStrategyFromString_Invalid(t *testing.T) {
	_, err := DefaultsStrategyFromString("unknown")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestNewDefaulter_NilMap_ReturnsError(t *testing.T) {
	_, err := NewDefaulter(nil, DefaultsApply)
	if err == nil {
		t.Fatal("expected error for nil defaults map")
	}
}

func TestDefaulter_Apply_SkipStrategy_DoesNotFill(t *testing.T) {
	d, err := NewDefaulter(map[string]string{"FOO": "bar", "BAZ": "qux"}, DefaultsSkip)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	secrets := map[string]string{"EXISTING": "val"}
	out := d.Apply(secrets)
	if _, ok := out["FOO"]; ok {
		t.Error("FOO should not be present with skip strategy")
	}
	if out["EXISTING"] != "val" {
		t.Errorf("EXISTING should be preserved")
	}
}

func TestDefaulter_Apply_ApplyStrategy_FillsMissing(t *testing.T) {
	d, err := NewDefaulter(map[string]string{"FOO": "default_foo", "BAR": "default_bar"}, DefaultsApply)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	secrets := map[string]string{"FOO": "existing_foo"}
	out := d.Apply(secrets)
	if out["FOO"] != "existing_foo" {
		t.Errorf("FOO should not be overwritten: got %q", out["FOO"])
	}
	if out["BAR"] != "default_bar" {
		t.Errorf("BAR should be filled from defaults: got %q", out["BAR"])
	}
}

func TestDefaulter_Apply_DoesNotMutateInput(t *testing.T) {
	d, _ := NewDefaulter(map[string]string{"NEW": "val"}, DefaultsApply)
	original := map[string]string{"A": "1"}
	d.Apply(original)
	if _, ok := original["NEW"]; ok {
		t.Error("Apply must not mutate the input map")
	}
}

func TestDefaulter_Keys_ReturnsDefaultKeys(t *testing.T) {
	d, _ := NewDefaulter(map[string]string{"X": "1", "Y": "2"}, DefaultsSkip)
	keys := d.Keys()
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}
