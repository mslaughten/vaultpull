package cmd

import (
	"testing"
)

func TestTransformCmd_RegisteredOnRoot(t *testing.T) {
	var found bool
	for _, c := range RootCmd.Commands() {
		if c.Name() == "transform" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'transform' command to be registered on root")
	}
}

func TestTransformCmd_DefaultFlags(t *testing.T) {
	cmd := transformCmd

	dryRun, err := cmd.Flags().GetBool("dry-run")
	if err != nil {
		t.Fatalf("dry-run flag not found: %v", err)
	}
	if dryRun {
		t.Error("expected dry-run default to be false")
	}

	rules, err := cmd.Flags().GetStringArray("rule")
	if err != nil {
		t.Fatalf("rule flag not found: %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("expected empty rules by default, got %v", rules)
	}
}

func TestTransformCmd_RequiresOneArg(t *testing.T) {
	cmd := transformCmd
	cmd.SetArgs([]string{})
	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error when no args provided")
	}
}

func TestTransformCmd_RequiresExactlyOneArg(t *testing.T) {
	cmd := transformCmd
	err := cmd.Args(cmd, []string{"file1", "file2"})
	if err == nil {
		t.Error("expected error when two args provided")
	}
}

func TestTransformCmd_AcceptsOneArg(t *testing.T) {
	cmd := transformCmd
	err := cmd.Args(cmd, []string{".env"})
	if err != nil {
		t.Errorf("unexpected error for single arg: %v", err)
	}
}
