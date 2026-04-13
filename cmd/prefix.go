package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/envfile"
	"vaultpull/internal/sync"
)

var prefixCmd = &cobra.Command{
	Use:   "prefix <env-file>",
	Short: "Add or strip a prefix from all keys in an env file",
	Args:  cobra.ExactArgs(1),
	RunE:  runPrefix,
}

func init() {
	prefixCmd.Flags().String("prefix", "", "Prefix string to add or strip (required)")
	prefixCmd.Flags().String("strategy", "add", "Strategy: add|strip")
	prefixCmd.Flags().Bool("dry-run", false, "Print result as JSON without writing")
	_ = prefixCmd.MarkFlagRequired("prefix")
	rootCmd.AddCommand(prefixCmd)
}

func runPrefix(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	prefix, _ := cmd.Flags().GetString("prefix")
	strategy, _ := cmd.Flags().GetString("strategy")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	p, err := sync.NewPrefixer(prefix, strategy)
	if err != nil {
		return fmt.Errorf("prefix: %w", err)
	}

	r := envfile.NewReader(filePath)
	secrets, err := r.Read()
	if err != nil {
		return fmt.Errorf("prefix: read %s: %w", filePath, err)
	}

	result, err := p.Apply(secrets)
	if err != nil {
		return fmt.Errorf("prefix: apply: %w", err)
	}

	if dryRun {
		return json.NewEncoder(os.Stdout).Encode(result)
	}

	w := envfile.NewWriter(filePath)
	if err := w.Write(result); err != nil {
		return fmt.Errorf("prefix: write %s: %w", filePath, err)
	}
	fmt.Fprintf(os.Stderr, "prefix: wrote %d keys to %s\n", len(result), filePath)
	return nil
}
