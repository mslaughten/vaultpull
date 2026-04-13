package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"vaultpull/internal/envfile"
	"vaultpull/internal/sync"
)

func init() {
	var strategy string
	var delimiter string
	var outputDir string
	var dryRun bool

	groupCmd := &cobra.Command{
		Use:   "group <env-file>",
		Short: "Split a .env file into grouped files by key prefix or namespace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGroup(args[0], sync.GroupStrategy(strategy), delimiter, outputDir, dryRun)
		},
	}

	groupCmd.Flags().StringVar(&strategy, "strategy", "prefix", "grouping strategy: prefix, namespace, flat")
	groupCmd.Flags().StringVar(&delimiter, "delimiter", "_", "delimiter used to split keys into group and local name")
	groupCmd.Flags().StringVar(&outputDir, "output-dir", ".", "directory to write grouped .env files into")
	groupCmd.Flags().BoolVar(&dryRun, "dry-run", false, "print groups without writing files")

	rootCmd.AddCommand(groupCmd)
}

func runGroup(src string, strategy sync.GroupStrategy, delimiter, outputDir string, dryRun bool) error {
	reader, err := envfile.NewReader(src)
	if err != nil {
		return fmt.Errorf("group: open %s: %w", src, err)
	}
	secrets, err := reader.Read()
	if err != nil {
		return fmt.Errorf("group: read %s: %w", src, err)
	}

	grouper, err := sync.NewGrouper(strategy, delimiter)
	if err != nil {
		return err
	}

	entries := grouper.Apply(secrets)
	for _, entry := range entries {
		fileName := entry.Name + ".env"
		if dryRun {
			fmt.Fprintf(os.Stdout, "[%s]\n", entry.Name)
			for k, v := range entry.Values {
				fmt.Fprintf(os.Stdout, "  %s=%s\n", k, v)
			}
			continue
		}
		dest := filepath.Join(outputDir, fileName)
		w, err := envfile.NewWriter(dest)
		if err != nil {
			return fmt.Errorf("group: create %s: %w", dest, err)
		}
		if err := w.Write(entry.Values); err != nil {
			return fmt.Errorf("group: write %s: %w", dest, err)
		}
	}
	return nil
}
