package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/envfile"
	"vaultpull/internal/sync"
)

var cloneCmd = &cobra.Command{
	Use:   "clone <env-file>",
	Short: "Clone secrets matching a key prefix to a new prefix inside an env file",
	Args:  cobra.ExactArgs(1),
	RunE:  runClone,
}

var (
	cloneSrc    string
	cloneDst    string
	cloneDryRun bool
)

func init() {
	cloneCmd.Flags().StringVar(&cloneSrc, "src", "", "source key prefix to match (required)")
	cloneCmd.Flags().StringVar(&cloneDst, "dst", "", "destination key prefix to write (required)")
	cloneCmd.Flags().BoolVar(&cloneDryRun, "dry-run", false, "print changes without writing to file")
	_ = cloneCmd.MarkFlagRequired("src")
	_ = cloneCmd.MarkFlagRequired("dst")
	rootCmd.AddCommand(cloneCmd)
}

func runClone(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	r, err := envfile.NewReader(filePath)
	if err != nil {
		return fmt.Errorf("clone: open %s: %w", filePath, err)
	}
	secrets, err := r.Read()
	if err != nil {
		return fmt.Errorf("clone: read %s: %w", filePath, err)
	}

	cloner, err := sync.NewCloner(cloneSrc, cloneDst, cloneDryRun)
	if err != nil {
		return err
	}

	out, results := cloner.Apply(secrets)

	fmt.Fprint(cmd.OutOrStdout(), sync.CloneSummary(results))

	if cloneDryRun {
		return nil
	}

	w, err := envfile.NewWriter(filePath)
	if err != nil {
		return fmt.Errorf("clone: create writer for %s: %w", filePath, err)
	}
	if err := w.Write(out); err != nil {
		fmt.Fprintf(os.Stderr, "clone: write error: %v\n", err)
		return err
	}
	return nil
}
