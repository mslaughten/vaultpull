package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempPromoteEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestPromoteCmd_RegisteredOnRoot(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Name() == "promote" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("promote command not registered on root")
	}
}

func TestPromoteCmd_RequiresExactlyTwoArgs(t *testing.T) {
	rootCmd.SetArgs([]string{"promote", "only-one"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error for wrong arg count")
	}
}

func TestPromoteCmd_DefaultFlags(t *testing.T) {
	f := promoteCmd.Flags().Lookup("strategy")
	if f == nil {
		t.Fatal("missing --strategy flag")
	}
	if f.DefValue != "missing" {
		t.Errorf("default strategy: got %q, want \"missing\"", f.DefValue)
	}
	if promoteCmd.Flags().Lookup("dry-run") == nil {
		t.Fatal("missing --dry-run flag")
	}
}

func TestPromoteCmd_DryRun_PrintsPromotedKeys(t *testing.T) {
	src := writeTempPromoteEnv(t, "NEW_KEY=hello\nSHARED=from_src\n")
	dst := writeTempPromoteEnv(t, "SHARED=from_dst\nEXISTING=keep\n")

	out := &strings.Builder{}
	rootCmd.SetOut(out)
	rootCmd.SetArgs([]string{"promote", "--dry-run", "--strategy=missing", src, dst})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := out.String()
	if !strings.Contains(result, "+ NEW_KEY") {
		t.Errorf("expected NEW_KEY promoted, got:\n%s", result)
	}
	if !strings.Contains(result, "~ SHARED (skipped)") {
		t.Errorf("expected SHARED skipped, got:\n%s", result)
	}
}

func TestPromoteCmd_All_OverwritesDst(t *testing.T) {
	src := writeTempPromoteEnv(t, "KEY=new_value\n")
	dstDir := t.TempDir()
	dst := filepath.Join(dstDir, "dst.env")
	if err := os.WriteFile(dst, []byte("KEY=old_value\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	rootCmd.SetArgs([]string{"promote", "--strategy=all", src, dst})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "KEY=new_value") {
		t.Errorf("expected KEY=new_value in dst, got:\n%s", string(data))
	}
}
