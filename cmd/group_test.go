package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempGroupEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestGroupCmd_RegisteredOnRoot(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "group <env-file>" {
			return
		}
	}
	t.Fatal("group command not registered on root")
}

func TestGroupCmd_RequiresOneArg(t *testing.T) {
	rootCmd.SetArgs([]string{"group"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when no arg provided")
	}
}

func TestGroupCmd_DefaultFlags(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use != "group <env-file>" {
			continue
		}
		if f := sub.Flags().Lookup("strategy"); f == nil || f.DefValue != "prefix" {
			t.Errorf("expected --strategy default 'prefix'")
		}
		if f := sub.Flags().Lookup("delimiter"); f == nil || f.DefValue != "_" {
			t.Errorf("expected --delimiter default '_'")
		}
		if f := sub.Flags().Lookup("dry-run"); f == nil || f.DefValue != "false" {
			t.Errorf("expected --dry-run default false")
		}
		return
	}
	t.Fatal("group command not found")
}

func TestGroupCmd_DryRun_PrintsGroups(t *testing.T) {
	env := writeTempGroupEnv(t, "APP_HOST=localhost\nAPP_PORT=8080\nDB_HOST=db\n")

	out := &strings.Builder{}
	rootCmd.SetOut(out)
	rootCmd.SetArgs([]string{"group", "--dry-run", "--strategy", "prefix", "--delimiter", "_", env})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGroupCmd_WritesFiles(t *testing.T) {
	dir := t.TempDir()
	env := writeTempGroupEnv(t, "APP_HOST=localhost\nDB_HOST=db\n")

	rootCmd.SetArgs([]string{"group", "--strategy", "prefix", "--delimiter", "_", "--output-dir", dir, env})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, name := range []string{"APP.env", "DB.env"} {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected file %s to exist: %v", name, err)
		}
	}
}
