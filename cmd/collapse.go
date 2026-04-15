package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/envfile"
	"vaultpull/internal/sync"
)

func init() {
	var (
		prefix   string
		outKey   string
		sep      string
		strategy string
		dryRun   bool
	)

	cmd := &cobra.Command{
		Use:   "collapse <envfile>",
		Short: "Collapse keys sharing a prefix into a single key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCollapse(args[0], prefix, outKey, sep, strategy, dryRun)
		},
	}

	cmd.Flags().StringVar(&prefix, "prefix", "", "Key prefix to match (required)")
	cmd.Flags().StringVar(&outKey, "out", "", "Output key name (required)")
	cmd.Flags().StringVar(&sep, "sep", ",", "Separator for concat strategy")
	cmd.Flags().StringVar(&strategy, "strategy", "first", "Collapse strategy: first|last|concat")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print result without writing")
	_ = cmd.MarkFlagRequired("prefix")
	_ = cmd.MarkFlagRequired("out")

	rootCmd.AddCommand(cmd)
}

func runCollapse(path, prefix, outKey, sep, strategy string, dryRun bool) error {
	strat, err := sync.CollapseStrategyFromString(strategy)
	if err != nil {
		return err
	}

	col, err := sync.NewCollapser(prefix, outKey, sep, strat)
	if err != nil {
		return err
	}

	r := envfile.NewReader(path)
	m, err := r.Read()
	if err != nil {
		return fmt.Errorf("collapse: read %s: %w", path, err)
	}

	result, err := col.Apply(m)
	if err != nil {
		return err
	}

	if dryRun {
		for k, v := range result {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
		}
		return nil
	}

	w := envfile.NewWriter(path)
	return w.Write(result)
}
