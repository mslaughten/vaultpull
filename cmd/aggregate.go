package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/envfile"
	"github.com/yourusername/vaultpull/internal/sync"
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
		Use:   "aggregate <file>",
		Short: "Combine keys sharing a prefix into a single output key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAggregate(args[0], prefix, outKey, sep, strategy, dryRun)
		},
	}

	cmd.Flags().StringVar(&prefix, "prefix", "", "key prefix to match (required)")
	cmd.Flags().StringVar(&outKey, "out-key", "", "destination key for aggregated value (required)")
	cmd.Flags().StringVar(&sep, "sep", ",", "separator used by concat and unique strategies")
	cmd.Flags().StringVar(&strategy, "strategy", "concat", "aggregation strategy: concat|count|unique")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print result without writing to file")

	_ = cmd.MarkFlagRequired("prefix")
	_ = cmd.MarkFlagRequired("out-key")

	rootCmd.AddCommand(cmd)
}

func runAggregate(file, prefix, outKey, sep, strategyStr string, dryRun bool) error {
	strategy, err := sync.AggregateStrategyFromString(strategyStr)
	if err != nil {
		return err
	}

	agg, err := sync.NewAggregator(prefix, outKey, sep, strategy)
	if err != nil {
		return err
	}

	reader := envfile.NewReader(file)
	m, err := reader.Read()
	if err != nil {
		return fmt.Errorf("aggregate: read %s: %w", file, err)
	}

	result, err := agg.Apply(m)
	if err != nil {
		return fmt.Errorf("aggregate: %w", err)
	}

	if dryRun {
		for k, v := range result {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
		}
		return nil
	}

	w := envfile.NewWriter(file)
	return w.Write(result)
}
