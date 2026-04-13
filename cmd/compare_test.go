package cmd

import (
	"testing"
)

func TestCompareCmd_RegisteredOnRoot(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Name() == "compare" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("compare command not registered on root")
	}
}

func TestCompareCmd_DefaultFlags(t *testing.T) {
	var cmd *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Name() == "compare" {
			cmd = c
			break
		}
	}
	if cmd == nil {
		t.Fatal("compare command not found")
	}
	f := cmd.Flags().Lookup("env-file")
	if f == nil {
		t.Fatal("expected --env-file flag")
	}
	if f.DefValue != ".env" {
		t.Fatalf("expected default '.env', got %q", f.DefValue)
	}
}

func TestCompareCmd_RequiresOneArg(t *testing.T) {
	var cmd *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Name() == "compare" {
			cmd = c
			break
		}
	}
	if cmd == nil {
		t.Fatal("compare command not found")
	}
	cmd.SetArgs([]string{})
	if err := cmd.Args(cmd, []string{}); err == nil {
		t.Fatal("expected error for missing argument")
	}
	if err := cmd.Args(cmd, []string{"one", "two"}); err == nil {
		t.Fatal("expected error for too many arguments")
	}
	if err := cmd.Args(cmd, []string{"one"}); err != nil {
		t.Fatalf("expected no error for single argument, got: %v", err)
	}
}
