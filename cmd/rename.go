package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultpull/vaultpull/internal/sync"
)

// renameRuleFlag is a JSON-encoded RenameRule for flag parsing.
type renameRuleFlag struct {
	Pattern     string `json:"pattern"`
	Replacement string `json:"replacement"`
}

var renameCmd = &cobra.Command{
	Use:   "rename <env-file>",
	Short: "Rename keys in an existing .env file using regex rules",
	Long: `Apply rename rules to keys in a .env file.
Rules are supplied as JSON objects: --rule '{"pattern":"^DB_(.+)","replacement":"DATABASE_$1"}'`,
	Args: cobra.ExactArgs(1),
	RunE: runRename,
}

var renameRulesRaw []string
var renameDryRun bool

func init() {
	renameCmd.Flags().StringArrayVar(&renameRulesRaw, "rule", nil, "rename rule as JSON {\"pattern\":\"...\",\"replacement\":\"...\"}")
	renameCmd.Flags().BoolVar(&renameDryRun, "dry-run", false, "print renamed keys without writing to file")
	RootCmd.AddCommand(renameCmd)
}

func runRename(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	var rules []sync.RenameRule
	for _, raw := range renameRulesRaw {
		var rf renameRuleFlag
		if err := json.Unmarshal([]byte(raw), &rf); err != nil {
			return fmt.Errorf("invalid rule JSON %q: %w", raw, err)
		}
		rules = append(rules, sync.RenameRule{Pattern: rf.Pattern, Replacement: rf.Replacement})
	}

	rn, err := sync.NewRenamer(rules)
	if err != nil {
		return err
	}

	reader, err := sync.NewEnvFileReader(filePath)
	if err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}
	existing, err := reader.Read()
	if err != nil {
		return fmt.Errorf("parsing env file: %w", err)
	}

	renamed := rn.Apply(existing)

	if renameDryRun {
		for k, v := range renamed {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, v)
		}
		return nil
	}

	w, err := sync.NewEnvFileWriter(filePath)
	if err != nil {
		return fmt.Errorf("opening env file for write: %w", err)
	}
	return w.Write(renamed)
}
