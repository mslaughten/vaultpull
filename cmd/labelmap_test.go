package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLabelmapCmd_RegisteredOnRoot(t *testing.T) {
	var found bool
	for _, c := range rootCmd.Commands() {
		if c.Use == "labelmap <env-file>" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("labelmap command not registered on root")
	}
}

func TestLabelmapCmd_RequiresOneArg(t *testing.T) {
	rootCmd.SetArgs([]string{"labelmap"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error when no arg provided")
	}
}

func TestLabelmapCmd_DefaultFlags(t *testing.T) {
	cmd := labelmapCmd
	if cmd.Flag("dry-run") == nil {
		t.Error("expected --dry-run flag")
	}
	if cmd.Flag("rule") == nil {
		t.Error("expected --rule flag")
	}
}

func writeTempLabelEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLabelmapCmd_DryRun_PrintsMappedKeys(t *testing.T) {
	p := writeTempLabelEnv(t, "DB_PASS=secret\nAPI=tok\n")

	var sb strings.Builder
	rootCmd.SetOut(&sb)
	rootCmd.SetArgs([]string{
		"labelmap", p,
		"--rule", "DB_PASS=DATABASE_PASSWORD",
		"--dry-run",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "DATABASE_PASSWORD") {
		t.Errorf("expected DATABASE_PASSWORD in output, got: %s", out)
	}
}

func TestLabelmapCmd_InvalidRule_ReturnsError(t *testing.T) {
	p := writeTempLabelEnv(t, "KEY=val\n")
	rootCmd.SetArgs([]string{"labelmap", p, "--rule", "BADFORMAT"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error for invalid rule format")
	}
}
