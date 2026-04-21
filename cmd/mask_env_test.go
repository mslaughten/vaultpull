package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempMaskEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempMaskEnv: %v", err)
	}
	return p
}

func TestMaskEnvCmd_RegisteredOnRoot(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "mask-env <file>" {
			return
		}
	}
	t.Fatal("mask-env command not registered on root")
}

func TestMaskEnvCmd_RequiresOneArg(t *testing.T) {
	rootCmd.SetArgs([]string{"mask-env"})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error with no args")
	}
}

func TestMaskEnvCmd_DefaultFlags(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use != "mask-env <file>" {
			continue
		}
		if f := sub.Flags().Lookup("mode"); f == nil || f.DefValue != "all" {
			t.Errorf("expected --mode default=all")
		}
		if f := sub.Flags().Lookup("symbol"); f == nil || f.DefValue != "*" {
			t.Errorf("expected --symbol default=*")
		}
		if f := sub.Flags().Lookup("reveal"); f == nil || f.DefValue != "0" {
			t.Errorf("expected --reveal default=0")
		}
		return
	}
	t.Fatal("mask-env command not found")
}

func TestMaskEnvCmd_MasksAllValues(t *testing.T) {
	p := writeTempMaskEnv(t, "SECRET=abc123\nPLAIN=hello\n")
	out := captureOutput(t, func() {
		rootCmd.SetArgs([]string{"mask-env", "--mode", "all", p})
		_ = rootCmd.Execute()
	})
	if strings.Contains(out, "abc123") {
		t.Errorf("expected abc123 to be masked, got: %q", out)
	}
	if !strings.Contains(out, "SECRET=") {
		t.Errorf("expected SECRET key in output, got: %q", out)
	}
}

func TestMaskEnvCmd_InvalidMode_ReturnsError(t *testing.T) {
	p := writeTempMaskEnv(t, "KEY=value\n")
	rootCmd.SetArgs([]string{"mask-env", "--mode", "middle", p})
	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error for unknown mode")
	}
}
