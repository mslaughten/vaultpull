package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/envfile"
	"vaultpull/internal/sync"
)

var sortCmd = &cobra.Command{
	Use:   "sort <env-file>",
	Short: "Re-order keys in a .env file by a chosen strategy",
	Args:  cobra.ExactArgs(1),
	RunE:  runSort,
}

func init() {
	sortCmd.Flags().StringP("strategy", "s", "alpha", "Sort strategy: alpha, alpha-desc, length, length-desc")
	sortCmd.Flags().BoolP("dry-run", "n", false, "Print sorted output without writing the file")
	rootCmd.AddCommand(sortCmd)
}

func runSort(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	strategyStr, _ := cmd.Flags().GetString("strategy")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	strategy, err := sync.SortStrategyFromString(strategyStr)
	if err != nil {
		return err
	}

	reader := envfile.NewReader(filePath)
	secrets, err := reader.Read()
	if err != nil {
		return fmt.Errorf("reading %s: %w", filePath, err)
	}

	sorter := sync.NewSorter(strategy)
	sorted, keys := sorter.Apply(secrets)

	if dryRun {
		for _, k := range keys {
			fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, sorted[k])
		}
		return nil
	}

	writer := envfile.NewWriter(filePath)
	if err := writer.Write(sorted); err != nil {
		return fmt.Errorf("writing %s: %w", filePath, err)
	}

	fmt.Fprintf(os.Stdout, "sorted %d keys in %s using strategy %q\n", len(keys), filePath, strategyStr)
	return nil
}
