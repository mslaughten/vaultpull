package sync

import (
	"strings"
	"testing"
)

func TestNewCloner_EmptySrc_ReturnsError(t *testing.T) {
	_, err := NewCloner("", "DST_", false)
	if err == nil {
		t.Fatal("expected error for empty src")
	}
}

func TestNewCloner_EmptyDst_ReturnsError(t *testing.T) {
	_, err := NewCloner("SRC_", "", false)
	if err == nil {
		t.Fatal("expected error for empty dst")
	}
}

func TestNewCloner_Valid(t *testing.T) {
	c, err := NewCloner("SRC_", "DST_", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Cloner")
	}
}

func TestCloner_Apply_ClonesMatchingKeys(t *testing.T) {
	c, _ := NewCloner("OLD_", "NEW_", false)
	secrets := map[string]string{
		"OLD_HOST": "localhost",
		"OLD_PORT": "5432",
		"KEEP_ME":  "yes",
	}
	out, results := c.Apply(secrets)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if out["NEW_HOST"] != "localhost" {
		t.Errorf("expected NEW_HOST=localhost, got %q", out["NEW_HOST"])
	}
	if out["NEW_PORT"] != "5432" {
		t.Errorf("expected NEW_PORT=5432, got %q", out["NEW_PORT"])
	}
	if out["KEEP_ME"] != "yes" {
		t.Errorf("original key KEEP_ME should be preserved")
	}
	// originals still present
	if out["OLD_HOST"] != "localhost" {
		t.Errorf("original OLD_HOST should still be in map")
	}
}

func TestCloner_Apply_DryRun_DoesNotWriteNewKeys(t *testing.T) {
	c, _ := NewCloner("OLD_", "NEW_", true)
	secrets := map[string]string{"OLD_HOST": "localhost"}
	out, results := c.Apply(secrets)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if _, ok := out["NEW_HOST"]; ok {
		t.Error("dry-run should not write NEW_HOST into output map")
	}
}

func TestCloner_Apply_NoMatch_ReturnsEmpty(t *testing.T) {
	c, _ := NewCloner("X_", "Y_", false)
	out, results := c.Apply(map[string]string{"A": "1", "B": "2"})
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
	if len(out) != 2 {
		t.Errorf("output map should retain original keys")
	}
}

func TestCloneSummary_NoResults(t *testing.T) {
	s := CloneSummary(nil)
	if !strings.Contains(s, "no keys matched") {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestCloneSummary_WithResults(t *testing.T) {
	results := []CloneResult{
		{SrcKey: "OLD_A", DstKey: "NEW_A"},
		{SrcKey: "OLD_B", DstKey: "NEW_B"},
	}
	s := CloneSummary(results)
	if !strings.Contains(s, "OLD_A -> NEW_A") {
		t.Errorf("summary missing OLD_A -> NEW_A: %q", s)
	}
	if !strings.Contains(s, "OLD_B -> NEW_B") {
		t.Errorf("summary missing OLD_B -> NEW_B: %q", s)
	}
}
