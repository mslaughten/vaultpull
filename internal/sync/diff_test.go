package sync

import (
	"testing"
)

func TestDiff_Added(t *testing.T) {
	existing := map[string]string{"FOO": "bar"}
	incoming := map[string]string{"FOO": "bar", "NEW_KEY": "value"}

	result := Diff(existing, incoming)

	if len(result.Added) != 1 || result.Added[0] != "NEW_KEY" {
		t.Errorf("expected Added=[NEW_KEY], got %v", result.Added)
	}
	if len(result.Removed) != 0 {
		t.Errorf("expected no removals, got %v", result.Removed)
	}
	if len(result.Changed) != 0 {
		t.Errorf("expected no changes, got %v", result.Changed)
	}
}

func TestDiff_Removed(t *testing.T) {
	existing := map[string]string{"FOO": "bar", "OLD_KEY": "gone"}
	incoming := map[string]string{"FOO": "bar"}

	result := Diff(existing, incoming)

	if len(result.Removed) != 1 || result.Removed[0] != "OLD_KEY" {
		t.Errorf("expected Removed=[OLD_KEY], got %v", result.Removed)
	}
	if len(result.Added) != 0 {
		t.Errorf("expected no additions, got %v", result.Added)
	}
}

func TestDiff_Changed(t *testing.T) {
	existing := map[string]string{"FOO": "old"}
	incoming := map[string]string{"FOO": "new"}

	result := Diff(existing, incoming)

	if len(result.Changed) != 1 || result.Changed[0] != "FOO" {
		t.Errorf("expected Changed=[FOO], got %v", result.Changed)
	}
	if len(result.Unchanged) != 0 {
		t.Errorf("expected no unchanged, got %v", result.Unchanged)
	}
}

func TestDiff_Unchanged(t *testing.T) {
	existing := map[string]string{"FOO": "bar", "BAZ": "qux"}
	incoming := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := Diff(existing, incoming)

	if result.HasChanges() {
		t.Errorf("expected no changes, but HasChanges() returned true")
	}
	if len(result.Unchanged) != 2 {
		t.Errorf("expected 2 unchanged keys, got %d", len(result.Unchanged))
	}
}

func TestDiff_Summary_NoChanges(t *testing.T) {
	result := Diff(map[string]string{"A": "1"}, map[string]string{"A": "1"})
	if result.Summary() != "no changes detected" {
		t.Errorf("unexpected summary: %s", result.Summary())
	}
}

func TestDiff_Summary_WithChanges(t *testing.T) {
	existing := map[string]string{"OLD": "v"}
	incoming := map[string]string{"NEW": "v"}
	result := Diff(existing, incoming)
	expected := "+1 added, -1 removed, ~0 changed"
	if result.Summary() != expected {
		t.Errorf("expected %q, got %q", expected, result.Summary())
	}
}
