package sync

import (
	"testing"
)

func TestBuildPlan_BasicPaths(t *testing.T) {
	paths := []string{
		"secret/myapp/database",
		"secret/myapp/cache",
	}
	plan := BuildPlan(paths, "")

	if len(plan.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(plan.Entries))
	}
	if plan.Entries[0].EnvFile != "database.env" {
		t.Errorf("expected database.env, got %s", plan.Entries[0].EnvFile)
	}
	if plan.Entries[1].EnvFile != "cache.env" {
		t.Errorf("expected cache.env, got %s", plan.Entries[1].EnvFile)
	}
}

func TestBuildPlan_StripsNamespace(t *testing.T) {
	paths := []string{
		"myteam/production/database",
		"myteam/production/api",
	}
	plan := BuildPlan(paths, "myteam/production")

	if len(plan.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(plan.Entries))
	}
	if plan.Entries[0].VaultPath != "myteam/production/database" {
		t.Errorf("VaultPath should remain unchanged, got %s", plan.Entries[0].VaultPath)
	}
	if plan.Entries[0].EnvFile != "database.env" {
		t.Errorf("expected database.env, got %s", plan.Entries[0].EnvFile)
	}
}

func TestBuildPlan_DeduplicatesFileNames(t *testing.T) {
	paths := []string{
		"ns/a/config",
		"ns/b/config",
	}
	plan := BuildPlan(paths, "")

	if len(plan.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(plan.Entries))
	}
	if plan.Entries[0].EnvFile != "config.env" {
		t.Errorf("first entry should be config.env, got %s", plan.Entries[0].EnvFile)
	}
	if plan.Entries[1].EnvFile != "config_1.env" {
		t.Errorf("second entry should be config_1.env, got %s", plan.Entries[1].EnvFile)
	}
}

func TestBuildPlan_EmptyPaths(t *testing.T) {
	plan := BuildPlan([]string{}, "")
	if len(plan.Entries) != 0 {
		t.Errorf("expected empty plan, got %d entries", len(plan.Entries))
	}
}

func TestPlanEntry_String(t *testing.T) {
	e := PlanEntry{VaultPath: "secret/app/db", EnvFile: "db.env"}
	got := e.String()
	expected := "secret/app/db -> db.env"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
