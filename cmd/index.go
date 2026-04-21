package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/envfile"
	"vaultpull/internal/sync"
)

func init() {
	var strategy string

	indexCmd := &cobra.Command{
		Use:   "index <env-file>",
		Short: "Print a positional index of keys in an .env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			strat, err := sync.IndexStrategyFromString(strategy)
			if err != nil {
				return err
			}

			r := envfile.NewReader(args[0])
			m, err := r.Read()
			if err != nil {
				return fmt.Errorf("reading %s: %w", args[0], err)
			}
			if len(m) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "(empty)")
				return nil
			}

			ix := sync.NewIndexer(strat)
			entries := ix.Build(m)
			w := cmd.OutOrStdout()
			for _, e := range entries {
				fmt.Fprintf(w, "%4d  %s\n", e.Position, e.Key)
			}
			return nil
		},
	}

	indexCmd.Flags().StringVar(&strategy, "strategy", "alpha",
		"Ordering strategy: alpha or insertion")

	_ = os.Stderr // satisfy import
	rootCmd.AddCommand(indexCmd)
}
