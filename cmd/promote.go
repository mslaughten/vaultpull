package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/envfile"
	"github.com/yourusername/vaultpull/internal/sync"
)

var promoteCmd = &cobra.Command{
	Use:   "promote <src-env-file> <dst-env-file>",
	Short: "Promote secrets from one env file into another",
	Long: `Reads secrets from the source .env file and promotes them into the
destination .env file according to the chosen strategy.

Strategies:
  missing  (default) only add keys absent in the destination
  all      overwrite all keys in the destination`,
	Args: cobra.ExactArgs(2),
	RunE: runPromote,
}

func init() {
	promoteCmd.Flags().StringP("strategy", "s", "missing", "promote strategy: missing|all")
	promoteCmd.Flags().BoolP("dry-run", "n", false, "print result without writing")
	rootCmd.AddCommand(promoteCmd)
}

func runPromote(cmd *cobra.Command, args []string) error {
	srcPath := args[0]
	dstPath := args[1]

	stratStr, _ := cmd.Flags().GetString("strategy")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	strat, err := sync.PromoteStrategyFromString(stratStr)
	if err != nil {
		return err
	}

	srcMap, err := envfile.NewReader(srcPath).Read()
	if err != nil {
		return fmt.Errorf("reading src: %w", err)
	}
	dstMap, err := envfile.NewReader(dstPath).Read()
	if err != nil {
		return fmt.Errorf("reading dst: %w", err)
	}

	promoter := sync.NewPromoter(strat)
	out, summary, err := promoter.Apply(srcMap, dstMap)
	if err != nil {
		return err
	}

	if dryRun {
		for _, k := range summary.Promoted {
			fmt.Fprintf(cmd.OutOrStdout(), "+ %s=%s\n", k, out[k])
		}
		for _, k := range summary.Skipped {
			fmt.Fprintf(cmd.OutOrStdout(), "~ %s (skipped)\n", k)
		}
		fmt.Fprintln(cmd.OutOrStdout(), summary.String())
		return nil
	}

	if err := envfile.NewWriter(dstPath).Write(out); err != nil {
		return fmt.Errorf("writing dst: %w", err)
	}
	fmt.Fprintln(os.Stderr, summary.String())
	return nil
}
