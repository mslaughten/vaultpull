package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/envfile"
	"vaultpull/internal/sync"
)

var validateCmd = &cobra.Command{
	Use:   "validate <env-file>",
	Short: "Validate a .env file against a set of rules",
	Args:  cobra.ExactArgs(1),
	RunE:  runValidate,
}

func init() {
	vlidateRules := validateCmd.Flags().StringSliceP(
		"rule", "r", nil,
		`validation rules in the form KEY=pattern (empty pattern = required present)`,
	)
	_ = vlidateRules
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) error {
	rulesRaw, _ := cmd.Flags().GetStringSlice("rule")

	validator, err := sync.NewValidator(rulesRaw)
	if err != nil {
		return fmt.Errorf("building validator: %w", err)
	}

	reader := envfile.NewReader(args[0])
	secrets, err := reader.Read()
	if err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}

	if err := validator.Validate(secrets); err != nil {
		ve, ok := err.(*sync.ValidationError)
		if ok {
			fmt.Fprintln(os.Stderr, "Validation errors:")
			for _, v := range ve.Violations {
				fmt.Fprintf(os.Stderr, "  - %s\n", v)
			}
		}
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout(), "All validation rules passed.")
	return nil
}
