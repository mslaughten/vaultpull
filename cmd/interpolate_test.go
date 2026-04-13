package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInterpolateCmd_RegisteredOnRoot(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "interpolate <env-file>" {
			return
		}
	}
	t.Fatal("interpolate command not registered on root")
}

func TestInterpolateCmd_RequiresOneArg(t *testing.T) {
	rootCmd.SetArgs([]string{"interpolate"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error for missing argument")
	}
}

func TestInterpolateCmd_DefaultFlags(t *testing.T) {
	f := interpolateCmd.Flags()
	if f.Lookup("dry-run") == nil {
		t.Error("expected --dry-run flag")
	}
	if f.Lookup("strict") == nil {
		t.Error("expected --strict flag")
	}
}

func TestInterpolateCmd_DryRun_PrintsResolved(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	content := "HOST=localhost\nPORT=5432\nDSN=postgres://${HOST}:${PORT}/mydb\n"
	if err := os.WriteFile(envPath, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	buf := &strings.Builder{}
	interpolateCmd.SetOut(buf)

	rootCmd.SetArgs([]string{"interpolate", "--dry-run", envPath})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "DSN") {
		t.Errorf("expected DSN in output, got: %s", out)
	}
}

func TestInterpolateCmd_Strict_MissingVar_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	content := "URL=http://${UNDEFINED_HOST}/path\n"
	if err := os.WriteFile(envPath, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	rootCmd.SetArgs([]string{"interpolate", "--strict", envPath})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error for undefined variable in strict mode")
	}
}
