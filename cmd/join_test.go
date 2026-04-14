package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestJoinCmd_RegisteredOnRoot(t *testing.T) {
	for _, c := range rootCmd.Commands() {
		if c.Use == "join <primary.env> <secondary.env>" {
			return
		}
	}
	t.Fatal("join command not registered on root")
}

func TestJoinCmd_DefaultFlags(t *testing.T) {
	var joinCmd = findCommand(rootCmd, "join")
	if joinCmd == nil {
		t.Fatal("join command not found")
	}
	if f := joinCmd.Flags().Lookup("strategy"); f == nil || f.DefValue != "first" {
		t.Errorf("expected strategy default 'first', got %v", f)
	}
	if f := joinCmd.Flags().Lookup("separator"); f == nil || f.DefValue != "," {
		t.Errorf("expected separator default ',', got %v", f)
	}
	if f := joinCmd.Flags().Lookup("dry-run"); f == nil || f.DefValue != "false" {
		t.Errorf("expected dry-run default false")
	}
}

func TestJoinCmd_RequiresExactlyTwoArgs(t *testing.T) {
	rootCmd.SetArgs([]string{"join", "only-one.env"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error with one arg")
	}
}

func TestJoinCmd_DryRun_PrintsMerged(t *testing.T) {
	dir := t.TempDir()
	primary := filepath.Join(dir, "primary.env")
	secondary := filepath.Join(dir, "secondary.env")
	os.WriteFile(primary, []byte("A=from-primary\nB=only-primary\n"), 0o644)
	os.WriteFile(secondary, []byte("A=from-secondary\nC=only-secondary\n"), 0o644)

	buf := &strings.Builder{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"join", "--dry-run", "--strategy", "last", primary, secondary})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestJoinCmd_InvalidStrategy_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	primary := filepath.Join(dir, "a.env")
	secondary := filepath.Join(dir, "b.env")
	os.WriteFile(primary, []byte("X=1\n"), 0o644)
	os.WriteFile(secondary, []byte("Y=2\n"), 0o644)

	rootCmd.SetArgs([]string{"join", "--strategy", "invalid", primary, secondary})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid strategy")
	}
}
