package sync

import (
	"testing"
)

func TestSortStrategyFromString_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  SortStrategy
	}{
		{"alpha", SortStrategyAlpha},
		{"asc", SortStrategyAlpha},
		{"", SortStrategyAlpha},
		{"alpha-desc", SortStrategyAlphaDesc},
		{"desc", SortStrategyAlphaDesc},
		{"length", SortStrategyLength},
		{"len", SortStrategyLength},
		{"length-desc", SortStrategyLengthDesc},
		{"len-desc", SortStrategyLengthDesc},
	}
	for _, tc := range cases {
		got, err := SortStrategyFromString(tc.input)
		if err != nil {
			t.Errorf("SortStrategyFromString(%q) unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("SortStrategyFromString(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestSortStrategyFromString_Invalid(t *testing.T) {
	_, err := SortStrategyFromString("random")
	if err == nil {
		t.Fatal("expected error for unknown strategy, got nil")
	}
}

func TestSorter_Alpha_OrdersAscending(t *testing.T) {
	s := NewSorter(SortStrategyAlpha)
	secrets := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	_, keys := s.Apply(secrets)
	if keys[0] != "APPLE" || keys[1] != "MANGO" || keys[2] != "ZEBRA" {
		t.Errorf("unexpected order: %v", keys)
	}
}

func TestSorter_AlphaDesc_OrdersDescending(t *testing.T) {
	s := NewSorter(SortStrategyAlphaDesc)
	secrets := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	_, keys := s.Apply(secrets)
	if keys[0] != "ZEBRA" || keys[1] != "MANGO" || keys[2] != "APPLE" {
		t.Errorf("unexpected order: %v", keys)
	}
}

func TestSorter_Length_ShortestFirst(t *testing.T) {
	s := NewSorter(SortStrategyLength)
	secrets := map[string]string{"LONGKEY": "a", "K": "b", "MED": "c"}
	_, keys := s.Apply(secrets)
	if keys[0] != "K" || keys[1] != "MED" || keys[2] != "LONGKEY" {
		t.Errorf("unexpected order: %v", keys)
	}
}

func TestSorter_LengthDesc_LongestFirst(t *testing.T) {
	s := NewSorter(SortStrategyLengthDesc)
	secrets := map[string]string{"LONGKEY": "a", "K": "b", "MED": "c"}
	_, keys := s.Apply(secrets)
	if keys[0] != "LONGKEY" || keys[1] != "MED" || keys[2] != "K" {
		t.Errorf("unexpected order: %v", keys)
	}
}

func TestSorter_Apply_PreservesValues(t *testing.T) {
	s := NewSorter(SortStrategyAlpha)
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, _ := s.Apply(secrets)
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("values were mutated: %v", out)
	}
}

func TestSorter_Apply_EmptyMap(t *testing.T) {
	s := NewSorter(SortStrategyAlpha)
	out, keys := s.Apply(map[string]string{})
	if len(out) != 0 || len(keys) != 0 {
		t.Errorf("expected empty results, got out=%v keys=%v", out, keys)
	}
}
