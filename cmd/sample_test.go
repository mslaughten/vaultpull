package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempSampleEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestSampleCmd_RegisteredOnRoot(t *testing.T) {
	for _, sub := range root.Commands() {
		if sub.Use == "sample <env-file>" {
			return
		}
	}
	t.Fatal("sample command not registered on root")
}

func TestSampleCmd_DefaultFlags(t *testing.T) {
	var sampleCmd = findCommand(root, "sample")
	if sampleCmd == nil {
		t.Fatal("sample command not found")
	}
	n, err := sampleCmd.Flags().GetInt("n")
	if err != nil || n != 5 {
		t.Errorf("expected default n=5, got %d", n)
	}
	strategy, err := sampleCmd.Flags().GetString("strategy")
	if err != nil || strategy != "first" {
		t.Errorf("expected default strategy=first, got %s", strategy)
	}
}

func TestSampleCmd_RequiresOneArg(t *testing.T) {
	root.SetArgs([]string{"sample"})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error when no args provided")
	}
}

func TestSampleCmd_DryRun_PrintsSampledKeys(t *testing.T) {
	path := writeTempSampleEnv(t, "ALPHA=1\nBETA=2\nCHARLIE=3\n")

	buf := &strings.Builder{}
	root.SetOut(buf)
	root.SetArgs([]string{"sample", "--n=2", "--strategy=first", "--dry-run", path})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSampleCmd_InvalidStrategy_ReturnsError(t *testing.T) {
	path := writeTempSampleEnv(t, "ALPHA=1\n")
	root.SetArgs([]string{"sample", "--strategy=bogus", path})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for invalid strategy")
	}
}
