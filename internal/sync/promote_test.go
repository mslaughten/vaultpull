package sync

import (
	"testing"
)

func TestPromoteStrategyFromString_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  PromoteStrategy
	}{
		{"missing", PromoteStrategyMissing},
		{"", PromoteStrategyMissing},
		{"all", PromoteStrategyAll},
	}
	for _, tc := range cases {
		got, err := PromoteStrategyFromString(tc.input)
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("input %q: got %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestPromoteStrategyFromString_Invalid(t *testing.T) {
	_, err := PromoteStrategyFromString("overwrite")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestPromoter_Missing_SkipsExistingKeys(t *testing.T) {
	p := NewPromoter(PromoteStrategyMissing)
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{"A": "original", "C": "3"}

	out, summary, err := p.Apply(src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "original" {
		t.Errorf("A should be unchanged, got %q", out["A"])
	}
	if out["B"] != "2" {
		t.Errorf("B should be promoted, got %q", out["B"])
	}
	if len(summary.Promoted) != 1 || summary.Promoted[0] != "B" {
		t.Errorf("promoted: %v", summary.Promoted)
	}
	if len(summary.Skipped) != 1 || summary.Skipped[0] != "A" {
		t.Errorf("skipped: %v", summary.Skipped)
	}
}

func TestPromoter_All_OverwritesExistingKeys(t *testing.T) {
	p := NewPromoter(PromoteStrategyAll)
	src := map[string]string{"A": "new", "B": "2"}
	dst := map[string]string{"A": "old"}

	out, summary, err := p.Apply(src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["A"] != "new" {
		t.Errorf("A should be overwritten, got %q", out["A"])
	}
	if len(summary.Promoted) != 2 {
		t.Errorf("expected 2 promoted, got %d", len(summary.Promoted))
	}
	if len(summary.Skipped) != 0 {
		t.Errorf("expected 0 skipped, got %d", len(summary.Skipped))
	}
}

func TestPromoter_EmptySrc_ReturnsUnchangedDst(t *testing.T) {
	p := NewPromoter(PromoteStrategyMissing)
	dst := map[string]string{"X": "1"}
	out, summary, err := p.Apply(map[string]string{}, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["X"] != "1" {
		t.Errorf("expected X=1, got %q", out["X"])
	}
	if len(summary.Promoted) != 0 {
		t.Errorf("expected 0 promoted")
	}
}

func TestPromoteSummary_String(t *testing.T) {
	s := PromoteSummary{Promoted: []string{"A", "B"}, Skipped: []string{"C"}}
	got := s.String()
	want := "promoted=2 skipped=1"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
