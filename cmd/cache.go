package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	syncp "github.com/vaultpull/internal/sync"
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage the local secret cache",
}

var cacheGetCmd = &cobra.Command{
	Use:   "get <vault-path>",
	Short: "Print cached secrets for a Vault path",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, _ := cmd.Flags().GetString("cache-dir")
		ttl, _ := cmd.Flags().GetDuration("ttl")
		c, err := syncp.NewSecretCache(dir, ttl)
		if err != nil {
			return err
		}
		entry, ok := c.Get(args[0])
		if !ok {
			fmt.Fprintln(cmd.OutOrStdout(), "cache miss")
			return nil
		}
		for k, v := range entry.Secrets {
			fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, v)
		}
		return nil
	},
}

var cacheInvalidateCmd = &cobra.Command{
	Use:   "invalidate <vault-path>",
	Short: "Remove the cached entry for a Vault path",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, _ := cmd.Flags().GetString("cache-dir")
		ttl, _ := cmd.Flags().GetDuration("ttl")
		c, err := syncp.NewSecretCache(dir, ttl)
		if err != nil {
			return err
		}
		if err := c.Invalidate(args[0]); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "invalidated")
		return nil
	},
}

var cacheClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Remove all entries from the local secret cache",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, _ := cmd.Flags().GetString("cache-dir")
		ttl, _ := cmd.Flags().GetDuration("ttl")
		c, err := syncp.NewSecretCache(dir, ttl)
		if err != nil {
			return err
		}
		if err := c.Clear(); err != nil {
			return fmt.Errorf("clearing cache: %w", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), "cache cleared")
		return nil
	},
}

func init() {
	defaultCacheDir := os.ExpandEnv("$HOME/.vaultpull/cache")

	for _, sub := range []*cobra.Command{cacheGetCmd, cacheInvalidateCmd, cacheClearCmd} {
		sub.Flags().String("cache-dir", defaultCacheDir, "directory for cached secrets")
		sub.Flags().Duration("ttl", 5*time.Minute, "cache entry time-to-live")
		cacheCmd.AddCommand(sub)
	}

	rootCmd.AddCommand(cacheCmd)
}
