package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultpull/vaultpull/internal/config"
)

var cfg = &config.Config{}

var rootCmd = &cobra.Command{
	Use:   "vaultpull",
	Short: "Sync HashiCorp Vault secrets into local .env files",
	Long: `vaultpull connects to a HashiCorp Vault instance and pulls secrets
from a specified KV mount path, optionally filtering by namespace,
and writes them to a local .env file.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfg.FromEnv()
		return cfg.Validate()
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfg.VaultAddr, "vault-addr", "", "Vault server address (env: VAULT_ADDR)")
	rootCmd.PersistentFlags().StringVar(&cfg.VaultToken, "vault-token", "", "Vault token (env: VAULT_TOKEN)")
	rootCmd.PersistentFlags().StringVar(&cfg.Namespace, "namespace", "", "Vault namespace prefix to filter secrets (env: VAULT_NAMESPACE)")
	rootCmd.PersistentFlags().StringVar(&cfg.MountPath, "mount", "secret", "KV secrets engine mount path")
	rootCmd.PersistentFlags().StringVar(&cfg.OutputFile, "output", ".env", "Output .env file path")
	rootCmd.PersistentFlags().IntVar(&cfg.KVVersion, "kv-version", 2, "KV engine version (1 or 2)")
}
