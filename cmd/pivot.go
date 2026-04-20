package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"

	"vaultpull/internal/envfile"
	"vaultpull/internal/sync"
)

var pivotCmd = &cobra.Command{
	Use:   "pivot <env-file>",
	Short: "Pivot secret keys and values using a chosen strategy",
	Long: `Reads a .env file and restructures it by swapping or remapping keys
and values according to the selected strategy.

Strategies:
  key_to_value (ktv)  - each key becomes the value; the original value becomes the key
  value_to_key (vtk)  - each value becomes the key; the original key becomes the value`,
	Args: cobra.ExactArgs(1),
	RunE: runPivot,
}

func init() {
	pivotCmd.Flags().StringP("strategy", "s", "value_to_key", "pivot strategy (key_to_value|value_to_key)")
	pivotCmd.Flags().StringP("prefix", "p", "", "prefix prepended to generated keys")
	pivotCmd.Flags().Bool("overwrite", false, "overwrite existing keys in the output")
	pivotCmd.Flags().Bool("dry-run", false, "print result without writing to disk")
	rootCmd.AddCommand(pivotCmd)
}

func runPivot(cmd *cobra.Command, args []string) error {
	path := args[0]
	strategyStr, _ := cmd.Flags().GetString("strategy")
	prefix, _ := cmd.Flags().GetString("prefix")
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	strategy, err := sync.PivotStrategyFromString(strategyStr)
	if err != nil {
		return err
	}

	r := envfile.NewReader(path)
	existing, err := r.Read()
	if err != nil {
		return fmt.Errorf("pivot: read %s: %w", path, err)
	}

	pivoter := sync.NewPivoter(strategy, prefix, overwrite)
	result, err := pivoter.Apply(map[string]string{}, existing)
	if err != nil {
		return fmt.Errorf("pivot: %w", err)
	}

	if dryRun {
		keys := make([]string, 0, len(result))
		for k := range result {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, result[k])
		}
		return nil
	}

	w := envfile.NewWriter(path)
	return w.Write(result)
}
