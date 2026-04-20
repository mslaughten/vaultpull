package sync

import (
	"testing"
)

func TestPivotStrategyFromString_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  PivotStrategy
	}{
		{"key_to_value", PivotStrategyKeyToValue},
		{"ktv", PivotStrategyKeyToValue},
		{"value_to_key", PivotStrategyValueToKey},
		{"vtk", PivotStrategyValueToKey},
		{"KEY_TO_VALUE", PivotStrategyKeyToValue},
	}
	for _, tc := range cases {
		got, err := PivotStrategyFromString(tc.input)
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("input %q: got %d, want %d", tc.input, got, tc.want)
		}
	}
}

func TestPivotStrategyFromString_Invalid(t *testing.T) {
	_, err := PivotStrategyFromString("bogus")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestPivoter_KeyToValue_SwapsKeyAndValue(t *testing.T) {
	p := NewPivoter(PivotStrategyKeyToValue, "", true)
	src := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	got, err := p.Apply(map[string]string{}, src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["localhost"] != "DB_HOST" {
		t.Errorf("expected localhost -> DB_HOST, got %q", got["localhost"])
	}
	if got["5432"] != "DB_PORT" {
		t.Errorf("expected 5432 -> DB_PORT, got %q", got["5432"])
	}
}

func TestPivoter_ValueToKey_PreservesMapping(t *testing.T) {
	p := NewPivoter(PivotStrategyValueToKey, "ENV_", true)
	src := map[string]string{"HOST": "localhost"}
	got, err := p.Apply(map[string]string{}, src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["ENV_HOST"] != "localhost" {
		t.Errorf("expected ENV_HOST=localhost, got %q", got["ENV_HOST"])
	}
}

func TestPivoter_NoOverwrite_KeepsExisting(t *testing.T) {
	p := NewPivoter(PivotStrategyValueToKey, "", false)
	dst := map[string]string{"HOST": "original"}
	src := map[string]string{"HOST": "new"}
	got, err := p.Apply(dst, src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["HOST"] != "original" {
		t.Errorf("expected original value preserved, got %q", got["HOST"])
	}
}

func TestPivoter_WithPrefix_PrependedToNewKeys(t *testing.T) {
	p := NewPivoter(PivotStrategyKeyToValue, "PIVOT_", true)
	src := map[string]string{"FOO": "bar"}
	got, err := p.Apply(map[string]string{}, src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["PIVOT_bar"] != "FOO" {
		t.Errorf("expected PIVOT_bar -> FOO, got %v", got)
	}
}

func TestPivoter_EmptySrc_ReturnsDstUnchanged(t *testing.T) {
	p := NewPivoter(PivotStrategyValueToKey, "", true)
	dst := map[string]string{"KEEP": "me"}
	got, err := p.Apply(dst, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got["KEEP"] != "me" {
		t.Errorf("expected dst unchanged, got %v", got)
	}
}
