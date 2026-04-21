package sync

import (
	"bytes"
	"testing"
)

func TestObserveStrategyFromString_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  ObserveStrategy
	}{
		{"all", ObserveAll},
		{"changed", ObserveChanged},
		{"missing", ObserveMissing},
		{"ALL", ObserveAll},
	}
	for _, tc := range cases {
		got, err := ObserveStrategyFromString(tc.input)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("got %q, want %q", got, tc.want)
		}
	}
}

func TestObserveStrategyFromString_Invalid(t *testing.T) {
	_, err := ObserveStrategyFromString("unknown")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNewObserver_NilReference_ReturnsError(t *testing.T) {
	_, err := NewObserver(nil, ObserveAll, nil)
	if err == nil {
		t.Fatal("expected error for nil reference")
	}
}

func TestObserver_All_ReportsAllKeys(t *testing.T) {
	ref := map[string]string{"A": "1", "B": "2"}
	curr := map[string]string{"A": "1", "B": "changed", "C": "new"}

	var buf bytes.Buffer
	obs, err := NewObserver(ref, ObserveAll, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	results, err := obs.Observe(curr)
	if err != nil {
		t.Fatalf("observe error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	statuses := map[string]string{}
	for _, r := range results {
		statuses[r.Key] = r.Status
	}
	if statuses["A"] != "unchanged" {
		t.Errorf("A: want unchanged, got %s", statuses["A"])
	}
	if statuses["B"] != "changed" {
		t.Errorf("B: want changed, got %s", statuses["B"])
	}
	if statuses["C"] != "added" {
		t.Errorf("C: want added, got %s", statuses["C"])
	}
}

func TestObserver_Changed_FiltersToChangedOnly(t *testing.T) {
	ref := map[string]string{"A": "1", "B": "2"}
	curr := map[string]string{"A": "1", "B": "99"}

	var buf bytes.Buffer
	obs, _ := NewObserver(ref, ObserveChanged, &buf)
	results, _ := obs.Observe(curr)

	if len(results) != 1 || results[0].Key != "B" {
		t.Errorf("expected only B changed, got %+v", results)
	}
}

func TestObserver_Missing_ReportsMissingKeys(t *testing.T) {
	ref := map[string]string{"A": "1", "B": "2", "C": "3"}
	curr := map[string]string{"A": "1"}

	var buf bytes.Buffer
	obs, _ := NewObserver(ref, ObserveMissing, &buf)
	results, _ := obs.Observe(curr)

	if len(results) != 2 {
		t.Fatalf("expected 2 missing, got %d", len(results))
	}
	for _, r := range results {
		if r.Status != "missing" {
			t.Errorf("expected missing status, got %s", r.Status)
		}
	}
}

func TestObserver_NilWriter_DefaultsToStdout(t *testing.T) {
	obs, err := NewObserver(map[string]string{"X": "1"}, ObserveAll, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obs.w == nil {
		t.Error("expected non-nil writer")
	}
}
