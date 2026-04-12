package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestSchemaCmd_RegisteredOnRoot(t *testing.T) {
	var found *cobra.Command
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "schema <vault-path>" {
			found = sub
			break
		}
	}
	if found == nil {
		t.Fatal("schema command not registered on root")
	}
}

func TestSchemaCmd_DefaultFlags(t *testing.T) {
	var found *cobra.Command
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "schema <vault-path>" {
			found = sub
			break
		}
	}
	if found == nil {
		t.Fatal("schema command not found")
	}

	schemaFile, err := found.Flags().GetString("schema-file")
	if err != nil {
		t.Fatalf("schema-file flag missing: %v", err)
	}
	if schemaFile != ".vaultschema.json" {
		t.Errorf("expected default schema-file '.vaultschema.json', got %q", schemaFile)
	}

	strict, err := found.Flags().GetBool("strict")
	if err != nil {
		t.Fatalf("strict flag missing: %v", err)
	}
	if strict {
		t.Error("expected strict to default to false")
	}
}

func TestSchemaCmd_RequiresOneArg(t *testing.T) {
	var found *cobra.Command
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "schema <vault-path>" {
			found = sub
			break
		}
	}
	if found == nil {
		t.Fatal("schema command not found")
	}

	err := found.Args(found, []string{})
	if err == nil {
		t.Error("expected error when no args provided")
	}

	err = found.Args(found, []string{"secret/myapp"})
	if err != nil {
		t.Errorf("unexpected error with one arg: %v", err)
	}

	err = found.Args(found, []string{"secret/myapp", "extra"})
	if err == nil {
		t.Error("expected error when two args provided")
	}
}
