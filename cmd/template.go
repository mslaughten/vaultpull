package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/sync"
	"vaultpull/internal/vault"
)

var templateCmd = &cobra.Command{
	Use:   "template <vault-path> <template-file>",
	Short: "Render a Go template using secrets from a Vault path",
	Long: `Reads secrets from a Vault KV path and renders a Go text/template.

The template receives a map[string]string of secret key/value pairs.
Use {{ index . "KEY" }} to reference individual secrets.

Example:
  vaultpull template secret/myapp configs/app.tmpl`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.FromEnv()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("invalid config: %w", err)
		}

		vaultPath := args[0]
		tmplFile := args[1]

		tmplSrc, err := os.ReadFile(tmplFile)
		if err != nil {
			return fmt.Errorf("read template file %q: %w", tmplFile, err)
		}

		client, err := vault.NewClient(cfg)
		if err != nil {
			return fmt.Errorf("create vault client: %w", err)
		}

		secrets, err := client.ReadSecretMap(cmd.Context(), vaultPath)
		if err != nil {
			return fmt.Errorf("read secrets from %q: %w", vaultPath, err)
		}

		renderer, err := sync.NewTemplateRenderer(string(tmplSrc), os.Stdout)
		if err != nil {
			return fmt.Errorf("parse template: %w", err)
		}

		if err := renderer.Render(secrets); err != nil {
			return fmt.Errorf("render template: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)
}
