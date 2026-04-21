package sync

import (
	"testing"
)

func TestAlignStrategyFromString_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  AlignStrategy
	}{
		{"intersection", AlignIntersection},
		{"union", AlignUnion},
		{"left", AlignLeft},
	}
	for _, tc := range cases {
		got, err := AlignStrategyFromString(tc.input)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("%s: got %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestAlignStrategyFromString_Invalid(t *testing.T) {
	_, err := AlignStrategyFromString("bogus")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestAligner_Intersection_KeepsCommonKeys(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2", "C": "3"}
	ref := map[string]string{"B": "x", "C": "y", "D": "z"}

	a := NewAligner(AlignIntersection, "")
	out := a.Apply(base, ref)

	if _, ok := out["A"]; ok {
		t.Error("A should have been removed")
	}
	if out["B"] != "2" || out["C"] != "3" {
		t.Errorf("unexpected values: %v", out)
	}
	if _, ok := out["D"]; ok {
		t.Error("D should not be present")
	}
}

func TestAligner_Union_FillsMissingWithFillValue(t *testing.T) {
	base := map[string]string{"A": "1"}
	ref := map[string]string{"B": "2"}

	a := NewAligner(AlignUnion, "MISSING")
	out := a.Apply(base, ref)

	if out["A"] != "1" {
		t.Errorf("A: got %q, want %q", out["A"], "1")
	}
	if out["B"] != "MISSING" {
		t.Errorf("B: got %q, want %q", out["B"], "MISSING")
	}
}

func TestAligner_Left_KeepsBaseAndFillsRefExtras(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	ref := map[string]string{"B": "x", "C": "y"}

	a := NewAligner(AlignLeft, "")
	out := a.Apply(base, ref)

	if out["A"] != "1" || out["B"] != "2" {
		t.Errorf("base keys should be preserved: %v", out)
	}
	if out["C"] != "" {
		t.Errorf("C should be filled with empty string, got %q", out["C"])
	}
}

func TestAligner_Intersection_EmptyBase(t *testing.T) {
	base := map[string]string{}
	ref := map[string]string{"A": "1"}

	a := NewAligner(AlignIntersection, "")
	out := a.Apply(base, ref)

	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestAligner_Union_BothEmpty(t *testing.T) {
	a := NewAligner(AlignUnion, "x")
	out := a.Apply(map[string]string{}, map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
