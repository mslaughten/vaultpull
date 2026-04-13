package sync

import (
	"bytes"
	"strings"
	"testing"
)

func TestCompare_InSync(t *testing.T) {
	c := NewComparer(nil)
	local := map[string]string{"A": "1", "B": "2"}
	vault := map[string]string{"A": "1", "B": "2"}
	res := c.Compare("test.env", local, vault)
	if res.HasDrift() {
		t.Fatalf("expected no drift, got %+v", res)
	}
	if len(res.InSync) != 2 {
		t.Fatalf("expected 2 in-sync keys, got %d", len(res.InSync))
	}
}

func TestCompare_OnlyLocal(t *testing.T) {
	c := NewComparer(nil)
	local := map[string]string{"A": "1", "LOCAL": "x"}
	vault := map[string]string{"A": "1"}
	res := c.Compare("test.env", local, vault)
	if !res.HasDrift() {
		t.Fatal("expected drift")
	}
	if len(res.OnlyLocal) != 1 || res.OnlyLocal[0] != "LOCAL" {
		t.Fatalf("unexpected OnlyLocal: %v", res.OnlyLocal)
	}
}

func TestCompare_OnlyVault(t *testing.T) {
	c := NewComparer(nil)
	local := map[string]string{"A": "1"}
	vault := map[string]string{"A": "1", "VAULT_ONLY": "v"}
	res := c.Compare("test.env", local, vault)
	if len(res.OnlyVault) != 1 || res.OnlyVault[0] != "VAULT_ONLY" {
		t.Fatalf("unexpected OnlyVault: %v", res.OnlyVault)
	}
}

func TestCompare_Diverged(t *testing.T) {
	c := NewComparer(nil)
	local := map[string]string{"KEY": "old"}
	vault := map[string]string{"KEY": "new"}
	res := c.Compare("test.env", local, vault)
	if len(res.Diverged) != 1 || res.Diverged[0] != "KEY" {
		t.Fatalf("unexpected Diverged: %v", res.Diverged)
	}
}

func TestCompare_Summary_NoDrift(t *testing.T) {
	c := NewComparer(nil)
	res := c.Compare("x.env", map[string]string{"K": "v"}, map[string]string{"K": "v"})
	if !strings.Contains(res.Summary(), "in sync") {
		t.Fatalf("expected 'in sync' in summary, got: %s", res.Summary())
	}
}

func TestCompare_Summary_WithDrift(t *testing.T) {
	c := NewComparer(nil)
	res := c.Compare("x.env", map[string]string{"A": "1"}, map[string]string{"B": "2"})
	s := res.Summary()
	if !strings.Contains(s, "drift detected") {
		t.Fatalf("expected drift detected, got: %s", s)
	}
	if !strings.Contains(s, "only-local") || !strings.Contains(s, "only-vault") {
		t.Fatalf("expected both sections in summary, got: %s", s)
	}
}

func TestCompare_Print_WritesToWriter(t *testing.T) {
	var buf bytes.Buffer
	c := NewComparer(&buf)
	res := c.Compare("f.env", map[string]string{"X": "1"}, map[string]string{"X": "1"})
	c.Print(res)
	if buf.Len() == 0 {
		t.Fatal("expected output written to writer")
	}
}
