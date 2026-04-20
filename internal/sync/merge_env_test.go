package sync

import (
	"testing"
)

func TestMergeEnvStrategyFromString_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  MergeStrategy
	}{
		{"overwrite", MergeOverwrite},
		{"keep", MergeKeepExisting},
		{"vault", MergeVaultWins},
	}
	for _, tc := range cases {
		got, err := MergeEnvStrategyFromString(tc.input)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("%s: got %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestMergeEnvStrategyFromString_Invalid(t *testing.T) {
	_, err := MergeEnvStrategyFromString("unknown")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestEnvMerger_Overwrite_ReplacesExisting(t *testing.T) {
	m := NewEnvMerger(MergeOverwrite, nil)
	dst := map[string]string{"A": "old", "B": "keep"}
	src := map[string]string{"A": "new", "C": "added"}
	out := m.Apply(dst, src)
	if out["A"] != "new" {
		t.Errorf("A: got %q, want \"new\"", out["A"])
	}
	if out["B"] != "keep" {
		t.Errorf("B: got %q, want \"keep\"", out["B"])
	}
	if out["C"] != "added" {
		t.Errorf("C: got %q, want \"added\"", out["C"])
	}
}

func TestEnvMerger_KeepExisting_DoesNotReplace(t *testing.T) {
	m := NewEnvMerger(MergeKeepExisting, nil)
	dst := map[string]string{"A": "original"}
	src := map[string]string{"A": "new", "B": "added"}
	out := m.Apply(dst, src)
	if out["A"] != "original" {
		t.Errorf("A: got %q, want \"original\"", out["A"])
	}
	if out["B"] != "added" {
		t.Errorf("B: got %q, want \"added\"", out["B"])
	}
}

func TestEnvMerger_VaultWins_PrefersNonEmpty(t *testing.T) {
	m := NewEnvMerger(MergeVaultWins, nil)
	dst := map[string]string{"A": "local", "B": "local"}
	src := map[string]string{"A": "vault", "B": ""}
	out := m.Apply(dst, src)
	if out["A"] != "vault" {
		t.Errorf("A: got %q, want \"vault\"", out["A"])
	}
	if out["B"] != "local" {
		t.Errorf("B: got %q, want \"local\"", out["B"])
	}
}

func TestEnvMerger_FilterKeys_OnlyMergesAllowed(t *testing.T) {
	m := NewEnvMerger(MergeOverwrite, []string{"A"})
	dst := map[string]string{"A": "old", "B": "old"}
	src := map[string]string{"A": "new", "B": "new", "C": "added"}
	out := m.Apply(dst, src)
	if out["A"] != "new" {
		t.Errorf("A: got %q, want \"new\"", out["A"])
	}
	if out["B"] != "old" {
		t.Errorf("B: got %q, want \"old\"", out["B"])
	}
	if _, ok := out["C"]; ok {
		t.Error("C should not be present")
	}
}

func TestEnvMerger_DoesNotMutateDst(t *testing.T) {
	m := NewEnvMerger(MergeOverwrite, nil)
	dst := map[string]string{"A": "original"}
	src := map[string]string{"A": "changed"}
	m.Apply(dst, src)
	if dst["A"] != "original" {
		t.Error("dst was mutated")
	}
}
