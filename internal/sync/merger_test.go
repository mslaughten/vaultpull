package sync

import (
	"testing"
)

func TestMerge_Overwrite(t *testing.T) {
	existing := map[string]string{"A": "local", "B": "local"}
	incoming := map[string]string{"A": "vault", "C": "vault"}

	got := Merge(existing, incoming, MergeStrategyOverwrite)

	if got["A"] != "vault" {
		t.Errorf("expected A=vault, got %s", got["A"])
	}
	if _, ok := got["B"]; ok {
		t.Error("expected B to be absent in overwrite mode")
	}
	if got["C"] != "vault" {
		t.Errorf("expected C=vault, got %s", got["C"])
	}
}

func TestMerge_KeepExisting(t *testing.T) {
	existing := map[string]string{"A": "local", "B": "local"}
	incoming := map[string]string{"A": "vault", "C": "vault"}

	got := Merge(existing, incoming, MergeStrategyKeepExisting)

	if got["A"] != "local" {
		t.Errorf("expected A=local (keep-existing), got %s", got["A"])
	}
	if got["B"] != "local" {
		t.Errorf("expected B=local, got %s", got["B"])
	}
	if got["C"] != "vault" {
		t.Errorf("expected C=vault (new key), got %s", got["C"])
	}
}

func TestMerge_VaultWins(t *testing.T) {
	existing := map[string]string{"A": "local", "B": "local"}
	incoming := map[string]string{"A": "vault", "C": "vault"}

	got := Merge(existing, incoming, MergeStrategyVaultWins)

	if got["A"] != "vault" {
		t.Errorf("expected A=vault (vault-wins), got %s", got["A"])
	}
	if got["B"] != "local" {
		t.Errorf("expected B=local (only in existing), got %s", got["B"])
	}
	if got["C"] != "vault" {
		t.Errorf("expected C=vault, got %s", got["C"])
	}
}

func TestMergeStrategyFromString(t *testing.T) {
	cases := []struct {
		input    string
		want     MergeStrategy
		wantOK   bool
	}{
		{"overwrite", MergeStrategyOverwrite, true},
		{"", MergeStrategyOverwrite, true},
		{"keep-existing", MergeStrategyKeepExisting, true},
		{"keep", MergeStrategyKeepExisting, true},
		{"vault-wins", MergeStrategyVaultWins, true},
		{"vault", MergeStrategyVaultWins, true},
		{"unknown", MergeStrategyOverwrite, false},
	}

	for _, tc := range cases {
		got, ok := MergeStrategyFromString(tc.input)
		if ok != tc.wantOK {
			t.Errorf("input %q: ok=%v, want %v", tc.input, ok, tc.wantOK)
		}
		if got != tc.want {
			t.Errorf("input %q: strategy=%v, want %v", tc.input, got, tc.want)
		}
	}
}
