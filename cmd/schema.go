package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/sync"
	"github.com/yourusername/vaultpull/internal/vault"
)

var schemaCmd = &cobra.Command{
	Use:   "schema <vault-path>",
	Short: "Validate secrets at a Vault path against a local schema file",
	Args:  cobra.ExactArgs(1),
	RunE:  runSchema,
}

func init() {
	schemaCmd.Flags().String("schema-file", sync.DefaultSchemaPath, "path to JSON schema file")
	schemaCmd.Flags().Bool("strict", false, "exit non-zero if any violations are found")
	rootCmd.AddCommand(schemaCmd)
}

func runSchema(cmd *cobra.Command, args []string) error {
	cfg, err := config.FromEnv()
	if err != nil {
		return err
	}
	if err := cfg.Validate(); err != nil {
		return err
	}

	schemaFile, _ := cmd.Flags().GetString("schema-file")
	strict, _ := cmd.Flags().GetBool("strict")

	v, err := sync.LoadSchemaFile(schemaFile)
	if err != nil {
		return err
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return err
	}

	secrets, err := client.ReadSecret(cmd.Context(), args[0])
	if err != nil {
		return fmt.Errorf("reading secret: %w", err)
	}

	violations := v.Validate(secrets)
	ok := v.WriteReport(os.Stdout, violations)

	if strict && !ok {
		return fmt.Errorf("schema validation failed with %d violation(s)", len(violations))
	}
	return nil
}
