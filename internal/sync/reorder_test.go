package sync

import (
	"testing"
)

func TestReorderStrategyFromString_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  ReorderStrategy
	}{
		{"explicit", ReorderStrategyExplicit},
		{"reverse", ReorderStrategyReverse},
		{"EXPLICIT", ReorderStrategyExplicit},
		{"REVERSE", ReorderStrategyReverse},
	}
	for _, tc := range cases {
		got, err := ReorderStrategyFromString(tc.input)
		if err != nil {
			t.Errorf("input %q: unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("input %q: got %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestReorderStrategyFromString_Invalid(t *testing.T) {
	_, err := ReorderStrategyFromString("bogus")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestNewReorderer_ExplicitRequiresKeys(t *testing.T) {
	_, err := NewReorderer(ReorderStrategyExplicit, nil)
	if err == nil {
		t.Fatal("expected error when no keys provided for explicit strategy")
	}
}

func TestReorderer_Explicit_PlacesListedKeysFirst(t *testing.T) {
	r, err := NewReorderer(ReorderStrategyExplicit, []string{"Z_KEY", "A_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := map[string]string{
		"A_KEY": "1",
		"B_KEY": "2",
		"Z_KEY": "3",
	}
	_, ordered, err := r.Apply(m)
	if err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	if ordered[0] != "Z_KEY" || ordered[1] != "A_KEY" {
		t.Errorf("expected Z_KEY then A_KEY first, got %v", ordered[:2])
	}
	if len(ordered) != 3 {
		t.Errorf("expected 3 keys, got %d", len(ordered))
	}
}

func TestReorderer_Explicit_SkipsMissingKeys(t *testing.T) {
	r, _ := NewReorderer(ReorderStrategyExplicit, []string{"MISSING", "A_KEY"})
	m := map[string]string{"A_KEY": "1", "B_KEY": "2"}
	_, ordered, err := r.Apply(m)
	if err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	if ordered[0] != "A_KEY" {
		t.Errorf("expected A_KEY first, got %v", ordered[0])
	}
	if len(ordered) != 2 {
		t.Errorf("expected 2 keys, got %d", len(ordered))
	}
}

func TestReorderer_Reverse_ReversesAlphaOrder(t *testing.T) {
	r, err := NewReorderer(ReorderStrategyReverse, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := map[string]string{"A": "1", "B": "2", "C": "3"}
	_, ordered, err := r.Apply(m)
	if err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	if ordered[0] != "C" || ordered[1] != "B" || ordered[2] != "A" {
		t.Errorf("expected [C B A], got %v", ordered)
	}
}

func TestReorderer_PreservesAllValues(t *testing.T) {
	r, _ := NewReorderer(ReorderStrategyReverse, nil)
	m := map[string]string{"X": "foo", "Y": "bar"}
	out, _, err := r.Apply(m)
	if err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	if out["X"] != "foo" || out["Y"] != "bar" {
		t.Errorf("values changed unexpectedly: %v", out)
	}
}
