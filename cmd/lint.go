package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/envfile"
	"vaultpull/internal/sync"
)

var lintCmd = &cobra.Command{
	Use:   "lint <env-file>",
	Short: "Check secret keys in a .env file against naming conventions",
	Args:  cobra.ExactArgs(1),
	RunE:  runLint,
}

var lintStrict bool

func init() {
	lintCmd.Flags().BoolVar(&lintStrict, "strict", false, "exit non-zero when any violations are found")
	rootCmd.AddCommand(lintCmd)
}

func runLint(cmd *cobra.Command, args []string) error {
	path := args[0]

	reader := envfile.NewReader(path)
	secrets, err := reader.Read()
	if err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}

	if len(secrets) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "lint: file is empty, nothing to check")
		return nil
	}

	linter := sync.NewLinter(nil)
	violations := linter.Check(secrets)
	summary := linter.Summary(violations)
	fmt.Fprintln(cmd.OutOrStdout(), summary)

	if lintStrict && len(violations) > 0 {
		os.Exit(1)
	}
	return nil
}
