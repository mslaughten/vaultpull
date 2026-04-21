package sync

import (
	"testing"
)

func TestIndexStrategyFromString_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  IndexStrategy
	}{
		{"alpha", IndexStrategyAlpha},
		{"Alpha", IndexStrategyAlpha},
		{"insertion", IndexStrategyInsertion},
		{"INSERTION", IndexStrategyInsertion},
	}
	for _, tc := range cases {
		got, err := IndexStrategyFromString(tc.input)
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("input %q: got %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestIndexStrategyFromString_Invalid(t *testing.T) {
	_, err := IndexStrategyFromString("random")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestIndexer_Alpha_OrdersKeys(t *testing.T) {
	m := map[string]string{
		"ZEBRA": "1",
		"APPLE": "2",
		"MANGO": "3",
	}
	ix := NewIndexer(IndexStrategyAlpha)
	entries := ix.Build(m)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Key != "APPLE" || entries[1].Key != "MANGO" || entries[2].Key != "ZEBRA" {
		t.Errorf("unexpected order: %v", entries)
	}
	for i, e := range entries {
		if e.Position != i {
			t.Errorf("entry %d: position mismatch, got %d", i, e.Position)
		}
	}
}

func TestIndexer_Insertion_PreservesOrder(t *testing.T) {
	m := map[string]string{"A": "1", "B": "2"}
	ix := NewIndexer(IndexStrategyInsertion)
	entries := ix.Build(m)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestLookup_Found(t *testing.T) {
	m := map[string]string{"FOO": "bar", "BAZ": "qux"}
	ix := NewIndexer(IndexStrategyAlpha)
	entries := ix.Build(m)
	pos := Lookup(entries, "FOO")
	if pos < 0 {
		t.Errorf("expected non-negative position, got %d", pos)
	}
}

func TestLookup_Missing(t *testing.T) {
	m := map[string]string{"FOO": "bar"}
	ix := NewIndexer(IndexStrategyAlpha)
	entries := ix.Build(m)
	pos := Lookup(entries, "MISSING")
	if pos != -1 {
		t.Errorf("expected -1, got %d", pos)
	}
}

func TestIndexer_EmptyMap(t *testing.T) {
	ix := NewIndexer(IndexStrategyAlpha)
	entries := ix.Build(map[string]string{})
	if len(entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(entries))
	}
}
