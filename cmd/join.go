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
	var separator string
	var dryRun bool
	var output string

	cmd := &cobra.Command{
		Use:   "join <primary.env> <secondary.env>",
		Short: "Merge two .env files using a join strategy",
		Long: `Merge two .env files together using concat, first, or last strategy.

concat  – values present in both files are joined with the separator.
first   – values in the primary file take precedence.
last    – values in the secondary file overwrite the primary.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runJoin(args[0], args[1], strategy, separator, dryRun, output)
		},
	}

	cmd.Flags().StringVarP(&strategy, "strategy", "s", "first", "join strategy: concat|first|last")
	cmd.Flags().StringVar(&separator, "separator", ",", "separator used by concat strategy")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print merged result without writing files")
	cmd.Flags().StringVarP(&output, "output", "o", "", "write merged result to this file (default: overwrite primary)")

	rootCmd.AddCommand(cmd)
}

func runJoin(primary, secondary, strategy, separator string, dryRun bool, output string) error {
	strat, err := sync.JoinStrategyFromString(strategy)
	if err != nil {
		return err
	}
	joiner, err := sync.NewJoiner(strat, separator)
	if err != nil {
		return err
	}

	dst, err := envfile.NewReader(primary).Read()
	if err != nil {
		return fmt.Errorf("reading primary file: %w", err)
	}
	src, err := envfile.NewReader(secondary).Read()
	if err != nil {
		return fmt.Errorf("reading secondary file: %w", err)
	}

	merged := joiner.Apply(dst, src)

	if dryRun {
		for _, k := range sortedKeys(merged) {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, merged[k])
		}
		return nil
	}

	dest := primary
	if output != "" {
		dest = output
	}
	return envfile.NewWriter(dest).Write(merged)
}
