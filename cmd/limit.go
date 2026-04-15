package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/envfile"
	"vaultpull/internal/sync"
)

var limitCmd = &cobra.Command{
	Use:   "limit <env-file>",
	Short: "Reduce an env file to at most N entries",
	Args:  cobra.ExactArgs(1),
	RunE:  runLimit,
}

func init() {
	limitCmd.Flags().IntP("max", "n", 10, "maximum number of keys to keep")
	limitCmd.Flags().StringP("strategy", "s", "alpha", "selection strategy: first|last|alpha")
	limitCmd.Flags().Bool("dry-run", false, "print result without writing")
	rootCmd.AddCommand(limitCmd)
}

func runLimit(cmd *cobra.Command, args []string) error {
	max, _ := cmd.Flags().GetInt("max")
	strategyStr, _ := cmd.Flags().GetString("strategy")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	strategy, err := sync.LimitStrategyFromString(strategyStr)
	if err != nil {
		return err
	}

	limiter, err := sync.NewLimiter(max, strategy)
	if err != nil {
		return err
	}

	reader, err := envfile.NewReader(args[0])
	if err != nil {
		return fmt.Errorf("limit: open %s: %w", args[0], err)
	}
	data, err := reader.Read()
	if err != nil {
		return fmt.Errorf("limit: read %s: %w", args[0], err)
	}

	out := limiter.Apply(data)

	if dryRun {
		for k, v := range out {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
		}
		return nil
	}

	w, err := envfile.NewWriter(args[0])
	if err != nil {
		return fmt.Errorf("limit: open writer %s: %w", args[0], err)
	}
	return w.Write(out)
}
