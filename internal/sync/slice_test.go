package sync

import (
	"testing"
)

func TestSliceStrategyFromString_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  SliceStrategy
	}{
		{"index", SliceStrategyIndex},
		{"first", SliceStrategyFirst},
		{"last", SliceStrategyLast},
		{"join", SliceStrategyJoin},
		{"INDEX", SliceStrategyIndex},
	}
	for _, tc := range cases {
		got, err := SliceStrategyFromString(tc.input)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("input %q: got %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestSliceStrategyFromString_Invalid(t *testing.T) {
	_, err := SliceStrategyFromString("unknown")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestSlicer_Index_ExpandsToNumberedKeys(t *testing.T) {
	s, _ := NewSlicer(SlicerOptions{Strategy: SliceStrategyIndex, Delimiter: ","})
	out, err := s.Apply(map[string]string{"HOSTS": "a, b, c"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["HOSTS_0"] != "a" || out["HOSTS_1"] != "b" || out["HOSTS_2"] != "c" {
		t.Errorf("unexpected output: %v", out)
	}
	if _, ok := out["HOSTS"]; ok {
		t.Error("original key should be removed on index strategy")
	}
}

func TestSlicer_First_KeepsFirstElement(t *testing.T) {
	s, _ := NewSlicer(SlicerOptions{Strategy: SliceStrategyFirst, Delimiter: ","})
	out, _ := s.Apply(map[string]string{"HOSTS": "alpha, beta, gamma"})
	if out["HOSTS"] != "alpha" {
		t.Errorf("got %q, want %q", out["HOSTS"], "alpha")
	}
}

func TestSlicer_Last_KeepsLastElement(t *testing.T) {
	s, _ := NewSlicer(SlicerOptions{Strategy: SliceStrategyLast, Delimiter: ";"})
	out, _ := s.Apply(map[string]string{"K": "x;y;z"})
	if out["K"] != "z" {
		t.Errorf("got %q, want %q", out["K"], "z")
	}
}

func TestSlicer_Join_CombinesWithSeparator(t *testing.T) {
	s, _ := NewSlicer(SlicerOptions{Strategy: SliceStrategyJoin, Delimiter: ",", Separator: "|"})
	out, _ := s.Apply(map[string]string{"ADDRS": "1.2.3.4, 5.6.7.8"})
	if out["ADDRS"] != "1.2.3.4|5.6.7.8" {
		t.Errorf("got %q, want %q", out["ADDRS"], "1.2.3.4|5.6.7.8")
	}
}

func TestSlicer_NoDelimiterInValue_PassesThrough(t *testing.T) {
	s, _ := NewSlicer(SlicerOptions{Strategy: SliceStrategyIndex, Delimiter: ","})
	out, _ := s.Apply(map[string]string{"KEY": "single"})
	if out["KEY"] != "single" {
		t.Errorf("got %q, want %q", out["KEY"], "single")
	}
}

func TestSlicer_DefaultDelimiter_IsComma(t *testing.T) {
	s, err := NewSlicer(SlicerOptions{Strategy: SliceStrategyFirst})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, _ := s.Apply(map[string]string{"K": "a,b,c"})
	if out["K"] != "a" {
		t.Errorf("got %q, want %q", out["K"], "a")
	}
}
