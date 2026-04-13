package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/envfile"
	"vaultpull/internal/sync"
	"vaultpull/internal/vault"
)

var compareCmd = &cobra.Command{
	Use:   "compare <vault-path>",
	Short: "Compare a local .env file against Vault secrets and report drift",
	Args:  cobra.ExactArgs(1),
	RunE:  runCompare,
}

var compareEnvFile string

func init() {
	compareCmd.Flags().StringVar(&compareEnvFile, "env-file", ".env", "Path to the local .env file")
	rootCmd.AddCommand(compareCmd)
}

func runCompare(cmd *cobra.Command, args []string) error {
	cfg, err := config.FromEnv()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("config validation: %w", err)
	}

	client, err := vault.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	vaultSecrets, err := client.ReadSecret(cmd.Context(), args[0])
	if err != nil {
		return fmt.Errorf("read vault secret: %w", err)
	}

	reader := envfile.NewReader(compareEnvFile)
	localSecrets, err := reader.Read()
	if err != nil {
		return fmt.Errorf("read env file: %w", err)
	}

	comparer := sync.NewComparer(os.Stdout)
	result := comparer.Compare(compareEnvFile, localSecrets, vaultSecrets)
	comparer.Print(result)

	if result.HasDrift() {
		os.Exit(1)
	}
	return nil
}
