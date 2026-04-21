package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/envfile"
	"github.com/yourusername/vaultpull/internal/sync"
)

func init() {
	var strategy string
	var fillValue string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "align <base.env> <ref.env>",
		Short: "Align keys between two .env files using a set strategy",
		Long: `Align compares two .env files and produces a result according to the
chosen strategy:
  intersection  keep only keys present in both files
  union         keep all keys; fill missing with --fill
  left          keep all base keys; add missing ref keys with --fill`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAlign(args[0], args[1], strategy, fillValue, dryRun)
		},
	}

	cmd.Flags().StringVarP(&strategy, "strategy", "s", "intersection", "align strategy (intersection|union|left)")
	cmd.Flags().StringVar(&fillValue, "fill", "", "fill value for missing keys (union/left strategies)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print result without writing")

	rootCmd.AddCommand(cmd)
}

func runAlign(basePath, refPath, strategy, fillValue string, dryRun bool) error {
	strat, err := sync.AlignStrategyFromString(strategy)
	if err != nil {
		return err
	}

	baseReader := envfile.NewReader(basePath)
	baseMap, err := baseReader.Read()
	if err != nil {
		return fmt.Errorf("reading base file: %w", err)
	}

	refReader := envfile.NewReader(refPath)
	refMap, err := refReader.Read()
	if err != nil {
		return fmt.Errorf("reading ref file: %w", err)
	}

	aligner := sync.NewAligner(strat, fillValue)
	out := aligner.Apply(baseMap, refMap)

	if dryRun {
		keys := make([]string, 0, len(out))
		for k := range out {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Printf("%s=%s\n", k, out[k])
		}
		return nil
	}

	w := envfile.NewWriter(basePath)
	if err := w.Write(out); err != nil {
		return fmt.Errorf("writing aligned env file: %w", err)
	}

	fmt.Fprintf(os.Stderr, "aligned %d keys into %s\n", len(out), basePath)
	return nil
}
