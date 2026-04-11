package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/sync"
	"github.com/yourusername/vaultpull/internal/vault"
)

var (
	exportFormat    string
	exportNamespace string
	exportMount     string
)

var exportCmd = &cobra.Command{
	Use:   "export <secret-path>",
	Short: "Print Vault secrets to stdout in JSON or dotenv format",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
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

		secretPath := args[0]
		secrets, err := client.ReadSecret(cmd.Context(), exportMount, secretPath)
		if err != nil {
			return fmt.Errorf("read secret: %w", err)
		}

		if exportNamespace != "" {
			secrets = vault.FilterByNamespace(secrets, exportNamespace)
		}

		ex, err := sync.NewExporter(sync.ExportFormat(exportFormat), os.Stdout)
		if err != nil {
			return err
		}
		return ex.Write(secrets)
	},
}

func init() {
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "dotenv", "Output format: json or dotenv")
	exportCmd.Flags().StringVarP(&exportNamespace, "namespace", "n", "", "Filter keys by namespace prefix")
	exportCmd.Flags().StringVarP(&exportMount, "mount", "m", "secret", "KV mount path")
	RootCmd.AddCommand(exportCmd)
}
