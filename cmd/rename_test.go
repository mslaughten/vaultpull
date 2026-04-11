package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenameCmd_RegisteredOnRoot(t *testing.T) {
	for _, sub := range RootCmd.Commands() {
		if sub.Use == "rename <env-file>" {
			return
		}
	}
	t.Fatal("rename command not registered on root")
}

func TestRenameCmd_RequiresOneArg(t *testing.T) {
	cmd := renameCmd
	cmd.SetArgs([]string{})
	var buf bytes.Buffer
	cmd.SetErr(&buf)
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when no args provided")
	}
}

func TestRenameCmd_DryRun_PrintsRenamedKeys(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	if err := os.WriteFile(envPath, []byte("DB_HOST=localhost\nDB_PORT=5432\n"), 0600); err != nil {
		t.Fatalf("setup: %v", err)
	}

	var out bytes.Buffer
	renameCmd.SetOut(&out)

	renameRulesRaw = []string{`{"pattern":"^DB_(.+)$","replacement":"DATABASE_$1"}`}
	renameDryRun = true
	t.Cleanup(func() {
		renameRulesRaw = nil
		renameDryRun = false
	})

	if err := runRename(renameCmd, []string{envPath}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := out.String()
	if !strings.Contains(got, "DATABASE_HOST") {
		t.Errorf("expected DATABASE_HOST in output, got:\n%s", got)
	}
}

func TestRenameCmd_InvalidRuleJSON_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	_ = os.WriteFile(envPath, []byte("FOO=bar\n"), 0600)

	renameRulesRaw = []string{`not-json`}
	t.Cleanup(func() { renameRulesRaw = nil })

	if err := runRename(renameCmd, []string{envPath}); err == nil {
		t.Fatal("expected error for invalid JSON rule")
	}
}

func TestRenameCmd_DefaultFlags(t *testing.T) {
	f := renameCmd.Flags().Lookup("dry-run")
	if f == nil {
		t.Fatal("dry-run flag not registered")
	}
	if f.DefValue != "false" {
		t.Errorf("expected default false, got %s", f.DefValue)
	}
}
