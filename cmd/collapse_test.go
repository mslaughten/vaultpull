package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCollapseCmd_RegisteredOnRoot(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "collapse <envfile>" {
			return
		}
	}
	t.Fatal("collapse command not registered on root")
}

func TestCollapseCmd_RequiresOneArg(t *testing.T) {
	rootCmd.SetArgs([]string{"collapse", "--prefix=DB_", "--out=DATABASE"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when no file arg provided")
	}
}

func TestCollapseCmd_DefaultFlags(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use != "collapse <envfile>" {
			continue
		}
		if f := sub.Flags().Lookup("strategy"); f == nil || f.DefValue != "first" {
			t.Error("expected --strategy default to be 'first'")
		}
		if f := sub.Flags().Lookup("sep"); f == nil || f.DefValue != "," {
			t.Error("expected --sep default to be ','")
		}
		if f := sub.Flags().Lookup("dry-run"); f == nil || f.DefValue != "false" {
			t.Error("expected --dry-run default to be false")
		}
		return
	}
	t.Fatal("collapse command not found")
}

func TestCollapseCmd_DryRun_PrintsCollapsed(t *testing.T) {
	dir := t.TempDir()
	env := filepath.Join(dir, ".env")
	if err := os.WriteFile(env, []byte("DB_HOST=localhost\nDB_PORT=5432\nOTHER=x\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	buf := &strings.Builder{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{
		"collapse", env,
		"--prefix=DB_",
		"--out=DATABASE",
		"--strategy=concat",
		"--sep=|",
		"--dry-run",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCollapseCmd_InvalidStrategy_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	env := filepath.Join(dir, ".env")
	_ = os.WriteFile(env, []byte("DB_HOST=localhost\n"), 0o644)

	rootCmd.SetArgs([]string{
		"collapse", env,
		"--prefix=DB_",
		"--out=DATABASE",
		"--strategy=bogus",
	})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}
