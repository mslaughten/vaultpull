package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/envfile"
	"vaultpull/internal/sync"
)

var observeCmd = &cobra.Command{
	Use:   "observe <env-file>",
	Short: "Compare a .env file against a reference and report key-level changes",
	Args:  cobra.ExactArgs(1),
	RunE:  runObserve,
}

func init() {
	observeCmd.Flags().StringP("reference", "r", "", "path to reference .env file (required)")
	observeCmd.Flags().StringP("strategy", "s", "all", "observe strategy: all|changed|missing")
	_ = observeCmd.MarkFlagRequired("reference")
	RootCmd.AddCommand(observeCmd)
}

func runObserve(cmd *cobra.Command, args []string) error {
	currentPath := args[0]
	refPath, _ := cmd.Flags().GetString("reference")
	strategyStr, _ := cmd.Flags().GetString("strategy")

	strategy, err := sync.ObserveStrategyFromString(strategyStr)
	if err != nil {
		return err
	}

	refReader, err := envfile.NewReader(refPath)
	if err != nil {
		return fmt.Errorf("observe: open reference: %w", err)
	}
	refMap, err := refReader.Read()
	if err != nil {
		return fmt.Errorf("observe: read reference: %w", err)
	}

	currReader, err := envfile.NewReader(currentPath)
	if err != nil {
		return fmt.Errorf("observe: open current: %w", err)
	}
	currMap, err := currReader.Read()
	if err != nil {
		return fmt.Errorf("observe: read current: %w", err)
	}

	obs, err := sync.NewObserver(refMap, strategy, os.Stdout)
	if err != nil {
		return err
	}

	results, err := obs.Observe(currMap)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Fprintln(os.Stdout, "no matching observations")
	}
	return nil
}
