package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/envfile"
	"github.com/yourusername/vaultpull/internal/sync"
)

var labelmapCmd = &cobra.Command{
	Use:   "labelmap <env-file>",
	Short: "Rename keys in a .env file using label=newkey rules",
	Args:  cobra.ExactArgs(1),
	RunE:  runLabelMap,
}

func init() {
	labelmapCmd.Flags().StringSliceP("rule", "r", nil, "Mapping rules in label=newkey format (repeatable)")
	labelmapCmd.Flags().Bool("dry-run", false, "Print result without writing to file")
	rootCmd.AddCommand(labelmapCmd)
}

func runLabelMap(cmd *cobra.Command, args []string) error {
	path := args[0]
	rules, _ := cmd.Flags().GetStringSlice("rule")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	lm, err := sync.NewLabelMapper(rules)
	if err != nil {
		return fmt.Errorf("labelmap: %w", err)
	}

	reader := envfile.NewReader(path)
	secrets, err := reader.Read()
	if err != nil {
		return fmt.Errorf("labelmap: read %s: %w", path, err)
	}

	mapped, err := lm.Apply(secrets)
	if err != nil {
		return fmt.Errorf("labelmap: apply: %w", err)
	}

	if dryRun {
		w := envfile.NewWriter(os.Stdout)
		return w.Write(mapped)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("labelmap: open %s: %w", path, err)
	}
	defer f.Close()
	return envfile.NewWriter(f).Write(mapped)
}
