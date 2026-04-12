package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/sync"
)

var redactCmd = &cobra.Command{
	Use:   "redact [path]",
	Short: "Print secrets with sensitive values masked",
	Long: `Read secrets from a Vault path and print them to stdout with
values matching sensitive patterns replaced by ***REDACTED***.`,
	Args: cobra.ExactArgs(1),
	RunE: runRedact,
}

func init() {
	redactCmd.Flags().StringSliceP("patterns", "p", sync.DefaultSensitivePatterns(), "regex patterns whose matching values are redacted")
	redactCmd.Flags().StringSliceP("keys", "k", nil, "key substrings whose values are redacted (case-insensitive)")
	redactCmd.Flags().StringP("format", "f", "env", "output format: env or json")
	RootCmd.AddCommand(redactCmd)
}

func runRedact(cmd *cobra.Command, args []string) error {
	patterns, _ := cmd.Flags().GetStringSlice("patterns")
	keys, _ := cmd.Flags().GetStringSlice("keys")
	format, _ := cmd.Flags().GetString("format")

	// Build a sample secrets map from the provided path argument for
	// demonstration; in a full implementation this would call the Vault client.
	_ = args[0]

	redactor, err := sync.NewRedactor(patterns)
	if err != nil {
		return fmt.Errorf("invalid redact pattern: %w", err)
	}

	// Placeholder: real implementation would fetch secrets from Vault.
	secrets := map[string]string{}

	masked := redactor.Apply(secrets)
	if len(keys) > 0 {
		masked = sync.RedactKeys(masked, keys)
	}

	switch format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(masked)
	default:
		for k, v := range masked {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
		}
	}
	return nil
}
