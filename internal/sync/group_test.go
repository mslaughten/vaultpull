package sync

import (
	"testing"
)

func TestNewGrouper_ValidStrategies(t *testing.T) {
	for _, s := range []GroupStrategy{GroupByPrefix, GroupByNamespace, GroupFlat} {
		_, err := NewGrouper(s, "_")
		if err != nil {
			t.Errorf("expected no error for strategy %q, got %v", s, err)
		}
	}
}

func TestNewGrouper_InvalidStrategy_ReturnsError(t *testing.T) {
	_, err := NewGrouper("bogus", "_")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestGrouper_Flat_ReturnsSingleGroup(t *testing.T) {
	g, _ := NewGrouper(GroupFlat, "_")
	secrets := map[string]string{"A": "1", "B": "2"}
	entries := g.Apply(secrets)
	if len(entries) != 1 {
		t.Fatalf("expected 1 group, got %d", len(entries))
	}
	if entries[0].Name != "default" {
		t.Errorf("expected group name 'default', got %q", entries[0].Name)
	}
	if entries[0].Values["A"] != "1" {
		t.Errorf("expected A=1")
	}
}

func TestGrouper_ByPrefix_SplitsOnDelimiter(t *testing.T) {
	g, _ := NewGrouper(GroupByPrefix, "_")
	secrets := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_HOST":  "db.local",
		"NOPREFIX": "value",
	}
	entries := g.Apply(secrets)
	if len(entries) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(entries))
	}
	// entries are sorted by name: APP, DB, default
	if entries[0].Name != "APP" {
		t.Errorf("expected first group APP, got %q", entries[0].Name)
	}
	if entries[0].Values["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost in APP group")
	}
	if entries[1].Name != "DB" {
		t.Errorf("expected second group DB, got %q", entries[1].Name)
	}
	if entries[2].Name != "default" {
		t.Errorf("expected third group default, got %q", entries[2].Name)
	}
}

func TestGrouper_ByNamespace_CustomDelimiter(t *testing.T) {
	g, _ := NewGrouper(GroupByNamespace, "/")
	secrets := map[string]string{
		"prod/API_KEY":    "secret",
		"prod/DB_PASS":    "pass",
		"staging/API_KEY": "dev-secret",
	}
	entries := g.Apply(secrets)
	if len(entries) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(entries))
	}
	if entries[0].Name != "prod" {
		t.Errorf("expected first group prod, got %q", entries[0].Name)
	}
	if entries[0].Values["API_KEY"] != "secret" {
		t.Errorf("expected API_KEY=secret in prod group")
	}
}

func TestGrouper_NoDelimiterInKey_FallsToDefault(t *testing.T) {
	g, _ := NewGrouper(GroupByPrefix, "_")
	secrets := map[string]string{"SIMPLE": "val"}
	entries := g.Apply(secrets)
	if len(entries) != 1 || entries[0].Name != "default" {
		t.Errorf("expected single default group")
	}
	if entries[0].Values["SIMPLE"] != "val" {
		t.Errorf("expected SIMPLE=val in default group")
	}
}
