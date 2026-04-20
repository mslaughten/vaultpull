package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultpull/internal/envfile"
	"github.com/your-org/vaultpull/internal/sync"
)

var reorderCmd = &cobra.Command{
	Use:   "reorder <file>",
	Short: "Reorder keys in a .env file",
	Long: `Reorder keys in a .env file using a chosen strategy.

Strategies:
  explicit  Place listed keys first (--keys), remaining keys follow in alpha order.
  reverse   Reverse alphabetical order of all keys.
`,
	Args: cobra.ExactArgs(1),
	RunE: runReorder,
}

func init() {
	reorderCmd.Flags().String("strategy", "reverse", "reorder strategy: explicit|reverse")
	reorderCmd.Flags().StringSlice("keys", nil, "ordered key list for explicit strategy")
	reorderCmd.Flags().Bool("dry-run", false, "print result without writing")
	rootCmd.AddCommand(reorderCmd)
}

func runReorder(cmd *cobra.Command, args []string) error {
	path := args[0]
	strategyStr, _ := cmd.Flags().GetString("strategy")
	keyList, _ := cmd.Flags().GetStringSlice("keys")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	strategy, err := sync.ReorderStrategyFromString(strategyStr)
	if err != nil {
		return err
	}

	r, err := sync.NewReorderer(strategy, keyList)
	if err != nil {
		return err
	}

	reader := envfile.NewReader(path)
	m, err := reader.Read()
	if err != nil {
		return fmt.Errorf("reorder: read %s: %w", path, err)
	}

	out, ordered, err := r.Apply(m)
	if err != nil {
		return fmt.Errorf("reorder: %w", err)
	}

	if dryRun {
		for _, k := range ordered {
			fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, out[k])
		}
		return nil
	}

	var sb strings.Builder
	for _, k := range ordered {
		v := out[k]
		if strings.ContainsAny(v, " \t\n#") {
			v = `"` + v + `"`
		}
		sb.WriteString(k + "=" + v + "\n")
	}

	if err := os.WriteFile(path, []byte(sb.String()), 0o644); err != nil {
		return fmt.Errorf("reorder: write %s: %w", path, err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "reordered %d keys in %s\n", len(ordered), path)
	return nil
}
