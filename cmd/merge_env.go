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
	var keys []string
	var dryRun bool
	var output string

	cmd := &cobra.Command{
		Use:   "merge-env <dst> <src>",
		Short: "Merge two .env files using a named strategy",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMergeEnv(args[0], args[1], strategy, keys, dryRun, output)
		},
	}

	cmd.Flags().StringVarP(&strategy, "strategy", "s", "overwrite", "merge strategy: overwrite|keep|vault")
	cmd.Flags().StringSliceVarP(&keys, "keys", "k", nil, "comma-separated keys to merge (default: all)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print result without writing")
	cmd.Flags().StringVarP(&output, "output", "o", "", "write result to this file instead of <dst>")

	rootCmd.AddCommand(cmd)
}

func runMergeEnv(dstPath, srcPath, strategyStr string, keys []string, dryRun bool, output string) error {
	strategy, err := sync.MergeEnvStrategyFromString(strategyStr)
	if err != nil {
		return err
	}

	dstReader := envfile.NewReader(dstPath)
	dstMap, err := dstReader.Read()
	if err != nil {
		return fmt.Errorf("reading dst %s: %w", dstPath, err)
	}

	srcReader := envfile.NewReader(srcPath)
	srcMap, err := srcReader.Read()
	if err != nil {
		return fmt.Errorf("reading src %s: %w", srcPath, err)
	}

	merger := sync.NewEnvMerger(strategy, keys)
	result := merger.Apply(dstMap, srcMap)

	if dryRun {
		w := envfile.NewWriter(os.Stdout)
		return w.Write(result)
	}

	dest := dstPath
	if output != "" {
		dest = output
	}

	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer f.Close()

	w := envfile.NewWriter(f)
	return w.Write(result)
}
