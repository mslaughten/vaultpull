package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/sync"
)

var checkpointCmd = &cobra.Command{
	Use:   "checkpoint",
	Short: "Manage sync checkpoints",
}

var checkpointShowCmd = &cobra.Command{
	Use:   "show <env-file>",
	Short: "Show the last checkpoint for an env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, _ := cmd.Flags().GetString("checkpoint-dir")
		store, err := sync.NewCheckpointStore(dir)
		if err != nil {
			return err
		}
		cp, err := store.Load(args[0])
		if err != nil {
			return err
		}
		if cp == nil {
			fmt.Fprintln(cmd.OutOrStdout(), "no checkpoint found")
			return nil
		}
		fmt.Fprintf(cmd.OutOrStdout(), "path:      %s\n", cp.Path)
		fmt.Fprintf(cmd.OutOrStdout(), "timestamp: %s\n", cp.Timestamp.Format("2006-01-02T15:04:05Z"))
		fmt.Fprintf(cmd.OutOrStdout(), "keys:      %d\n", len(cp.Hashes))
		return nil
	},
}

var checkpointDeleteCmd = &cobra.Command{
	Use:   "delete <env-file>",
	Short: "Delete the checkpoint for an env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, _ := cmd.Flags().GetString("checkpoint-dir")
		store, err := sync.NewCheckpointStore(dir)
		if err != nil {
			return err
		}
		if err := store.Delete(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "checkpoint deleted")
		return nil
	},
}

func init() {
	const defaultDir = ".vaultpull/checkpoints"
	checkpointShowCmd.Flags().String("checkpoint-dir", defaultDir, "directory for checkpoint files")
	checkpointDeleteCmd.Flags().String("checkpoint-dir", defaultDir, "directory for checkpoint files")
	checkpointCmd.AddCommand(checkpointShowCmd)
	checkpointCmd.AddCommand(checkpointDeleteCmd)
	rootCmd.AddCommand(checkpointCmd)
}
