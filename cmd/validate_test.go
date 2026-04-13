package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	return f.Name()
}

func TestValidateCmd_RegisteredOnRoot(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "validate <env-file>" {
			return
		}
	}
	t.Fatal("validate command not registered on root")
}

func TestValidateCmd_RequiresOneArg(t *testing.T) {
	rootCmd.SetArgs([]string{"validate"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when no arg provided")
	}
}

func TestValidateCmd_PassesWithMatchingRule(t *testing.T) {
	envFile := writeTempEnv(t, "PORT=8080\n")
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"validate", envFile, "--rule", "PORT=^[0-9]+$"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "passed") {
		t.Errorf("expected success message, got: %q", buf.String())
	}
}

func TestValidateCmd_FailsOnMissingKey(t *testing.T) {
	envFile := writeTempEnv(t, "OTHER=value\n")
	rootCmd.SetArgs([]string{"validate", envFile, "--rule", "REQUIRED_KEY="})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "REQUIRED_KEY") {
		t.Errorf("expected error to mention REQUIRED_KEY, got: %v", err)
	}
}

func TestValidateCmd_DefaultFlags(t *testing.T) {
	envFile := filepath.Join(t.TempDir(), "empty.env")
	_ = os.WriteFile(envFile, []byte(""), 0o644)
	rootCmd.SetArgs([]string{"validate", envFile})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error with no rules: %v", err)
	}
}
