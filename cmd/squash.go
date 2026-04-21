package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/envfile"
	"github.com/yourusername/vaultpull/internal/sync"
)

func init() {
	var prefix, outKey, sep, strategy string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "squash <file>",
		Short: "Collapse keys sharing a prefix into a single key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSquash(args[0], prefix, outKey, sep, strategy, dryRun)
		},
	}

	cmd.Flags().StringVar(&prefix, "prefix", "", "Key prefix to match (required)")
	cmd.Flags().StringVar(&outKey, "out", "", "Output key name (required)")
	cmd.Flags().StringVar(&sep, "sep", ",", "Separator for concat strategy")
	cmd.Flags().StringVar(&strategy, "strategy", "concat", "Squash strategy: concat|first|last")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print result without writing")

	_ = cmd.MarkFlagRequired("prefix")
	_ = cmd.MarkFlagRequired("out")

	rootCmd.AddCommand(cmd)
}

func runSquash(file, prefix, outKey, sep, strategyStr string, dryRun bool) error {
	strategy, err := sync.SquashStrategyFromString(strategyStr)
	if err != nil {
		return err
	}

	squasher, err := sync.NewSquasher(prefix, outKey, sep, strategy)
	if err != nil {
		return err
	}

	r := envfile.NewReader(file)
	secrets, err := r.Read()
	if err != nil {
		return fmt.Errorf("squash: read %s: %w", file, err)
	}

	result, err := squasher.Apply(secrets)
	if err != nil {
		return fmt.Errorf("squash: %w", err)
	}

	if dryRun {
		w := envfile.NewWriter(os.Stdout)
		return w.Write(result)
	}

	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("squash: open %s: %w", file, err)
	}
	defer f.Close()

	return envfile.NewWriter(f).Write(result)
}
