package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/envfile"
	"vaultpull/internal/sync"
)

func init() {
	var strategy string
	var delimiter string
	var dryRun bool

	splitCmd := &cobra.Command{
		Use:   "split <env-file>",
		Short: "Partition a .env file into named buckets by key prefix or delimiter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSplit(args[0], strategy, delimiter, dryRun)
		},
	}

	splitCmd.Flags().StringVarP(&strategy, "strategy", "s", "prefix", "split strategy: prefix or delimiter")
	splitCmd.Flags().StringVarP(&delimiter, "delimiter", "d", "_", "delimiter character used to split keys")
	splitCmd.Flags().BoolVar(&dryRun, "dry-run", false, "print buckets as JSON without writing files")

	rootCmd.AddCommand(splitCmd)
}

func runSplit(path, strategy, delimiter string, dryRun bool) error {
	r := envfile.NewReader(path)
	secrets, err := r.Read()
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	strat, err := sync.SplitStrategyFromString(strategy)
	if err != nil {
		return err
	}

	sp, err := sync.NewSplitter(strat, delimiter)
	if err != nil {
		return err
	}

	buckets := sp.Apply(secrets)

	if dryRun {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(buckets)
	}

	for name, kv := range buckets {
		outPath := name + ".env"
		w := envfile.NewWriter(outPath)
		if err := w.Write(kv); err != nil {
			return fmt.Errorf("write %s: %w", outPath, err)
		}
		fmt.Fprintf(os.Stdout, "wrote %s (%d keys)\n", outPath, len(kv))
	}
	return nil
}
