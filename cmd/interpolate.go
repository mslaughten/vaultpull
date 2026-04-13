package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/envfile"
	"github.com/yourusername/vaultpull/internal/sync"
)

var (
	interpolateDryRun bool
	interpolateStrict bool
)

var interpolateCmd = &cobra.Command{
	Use:   "interpolate <env-file>",
	Short: "Resolve variable references inside a .env file",
	Long: `Reads a .env file, expands ${VAR} and $VAR references using the
values already present in the file, and writes the result back in-place.

Use --dry-run to preview substitutions without modifying the file.
Use --strict to fail when a referenced variable is not defined.`,
	Args: cobra.ExactArgs(1),
	RunE: runInterpolate,
}

func init() {
	interpolateCmd.Flags().BoolVar(&interpolateDryRun, "dry-run", false, "print result without writing")
	interpolateCmd.Flags().BoolVar(&interpolateStrict, "strict", false, "error on undefined variables")
	rootCmd.AddCommand(interpolateCmd)
}

func runInterpolate(cmd *cobra.Command, args []string) error {
	path := args[0]

	reader := envfile.NewReader(path)
	secrets, err := reader.Read()
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	opts := []sync.InterpolatorOption{}
	if interpolateStrict {
		opts = append(opts, sync.WithStrictInterpolation())
	}

	ip := sync.NewInterpolator(secrets, opts...)
	resolved, err := ip.Apply(secrets)
	if err != nil {
		return err
	}

	if interpolateDryRun {
		for k, v := range resolved {
			fmt.Fprintf(cmd.OutOrStdout(), "%s=%s\n", k, v)
		}
		return nil
	}

	w := envfile.NewWriter(path)
	if err := w.Write(resolved); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	fmt.Fprintf(os.Stdout, "interpolated %d keys in %s\n", len(resolved), path)
	return nil
}
