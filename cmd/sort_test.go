package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSortCmd_RegisteredOnRoot(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "sort <env-file>" {
			return
		}
	}
	t.Fatal("sort command not registered on root")
}

func TestSortCmd_DefaultFlags(t *testing.T) {
	f := sortCmd.Flags().Lookup("strategy")
	if f == nil {
		t.Fatal("expected --strategy flag")
	}
	if f.DefValue != "alpha" {
		t.Errorf("default strategy = %q, want \"alpha\"", f.DefValue)
	}
	dr := sortCmd.Flags().Lookup("dry-run")
	if dr == nil {
		t.Fatal("expected --dry-run flag")
	}
}

func TestSortCmd_RequiresOneArg(t *testing.T) {
	rootCmd.SetArgs([]string{"sort"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error when no arg provided")
	}
}

func TestSortCmd_DryRun_PrintsAlphaOrder(t *testing.T) {
	dir := t.TempDir()
	env := filepath.Join(dir, ".env")
	_ = os.WriteFile(env, []byte("ZEBRA=1\nAPPLE=2\nMango=3\n"), 0o600)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"sort", "--dry-run", "--strategy", "alpha", env})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d: %q", len(lines), out)
	}
	// alpha order: APPLE < Mango < ZEBRA (uppercase before lower in ASCII, but Go sort is byte-wise)
	if !strings.HasPrefix(lines[0], "APPLE") {
		t.Errorf("first line should start with APPLE, got %q", lines[0])
	}
}

func TestSortCmd_InvalidStrategy_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	env := filepath.Join(dir, ".env")
	_ = os.WriteFile(env, []byte("FOO=bar\n"), 0o600)

	rootCmd.SetArgs([]string{"sort", "--strategy", "bogus", env})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error for invalid strategy")
	}
}
