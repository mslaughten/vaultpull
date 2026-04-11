package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"vaultpull/internal/config"
	"vaultpull/internal/sync"
	"vaultpull/internal/vault"
)

var (
	filterInclude []string
	filterExclude []string
	filterPath    string
)

var filterCmd = &cobra.Command{
	Use:   "filter <mount>",
	Short: "Print secrets from a Vault path filtered by key patterns",
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

		mount := args[0]
		path := strings.TrimPrefix(filterPath, "/")

		secrets, err := client.ReadSecret(mount, path)
		if err != nil {
			return fmt.Errorf("read secret: %w", err)
		}

		f, err := sync.NewKeyFilter(filterInclude, filterExclude)
		if err != nil {
			return fmt.Errorf("build filter: %w", err)
		}

		filtered := f.Apply(secrets)
		for k, v := range filtered {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
		}
		return nil
	},
}

func init() {
	filterCmd.Flags().StringSliceVar(&filterInclude, "include", nil, "Key patterns to include (supports * wildcard)")
	filterCmd.Flags().StringSliceVar(&filterExclude, "exclude", nil, "Key patterns to exclude (supports * wildcard)")
	filterCmd.Flags().StringVar(&filterPath, "path", "", "Sub-path within the mount to read")
	rootCmd.AddCommand(filterCmd)
}
