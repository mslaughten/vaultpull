package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/sync"
)

var transformCmd = &cobra.Command{
	Use:   "transform <env-file>",
	Short: "Apply value transformations to an existing .env file",
	Long: `Read a .env file, apply named transformations to specified keys,
and write the result back (or print with --dry-run).

Available transforms: upper, lower, trimspace, base64

Example:
  vaultpull transform .env --rule DB_PASS=upper --rule API_KEY=trimspace`,
	Args: cobra.ExactArgs(1),
	RunE: runTransform,
}

func init() {
	transformCmd.Flags().StringArrayP("rule", "r", nil, "Transform rule in KEY=transform format (repeatable)")
	transformCmd.Flags().Bool("dry-run", false, "Print transformed values without writing to file")
	RootCmd.AddCommand(transformCmd)
}

func runTransform(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	rules, _ := cmd.Flags().GetStringArray("rule")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	if len(rules) == 0 {
		return fmt.Errorf("at least one --rule is required")
	}

	transformer, err := sync.NewTransformer(rules)
	if err != nil {
		return err
	}

	reader := envfile.NewReader(filePath)
	secrets, err := reader.Read()
	if err != nil {
		return fmt.Errorf("reading %s: %w", filePath, err)
	}

	transformed, err := transformer.Apply(secrets)
	if err != nil {
		return err
	}

	if dryRun {
		for k, v := range transformed {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
		}
		return nil
	}

	writer := envfile.NewWriter(filePath)
	if err := writer.Write(transformed); err != nil {
		return fmt.Errorf("writing %s: %w", filePath, err)
	}
	fmt.Fprintf(os.Stdout, "transformed %d key(s) in %s\n", len(rules), filePath)
	return nil
}
