package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/envfile"
	"github.com/yourusername/vaultpull/internal/sync"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot <label> <envfile>",
	Short: "Capture or restore a named snapshot of an env file",
	Args:  cobra.ExactArgs(2),
}

var snapshotSaveCmd = &cobra.Command{
	Use:   "save <label> <envfile>",
	Short: "Save a snapshot of an env file under the given label",
	Args:  cobra.ExactArgs(2),
	RunE:  runSnapshotSave,
}

var snapshotLoadCmd = &cobra.Command{
	Use:   "load <label>",
	Short: "Print a previously saved snapshot",
	Args:  cobra.ExactArgs(1),
	RunE:  runSnapshotLoad,
}

var snapshotDir string

func init() {
	snapshotCmd.PersistentFlags().StringVar(&snapshotDir, "dir", ".vaultpull/snapshots", "directory for snapshot storage")
	snapshotCmd.AddCommand(snapshotSaveCmd)
	snapshotCmd.AddCommand(snapshotLoadCmd)
	rootCmd.AddCommand(snapshotCmd)
}

func runSnapshotSave(cmd *cobra.Command, args []string) error {
	label, envPath := args[0], args[1]

	reader := envfile.NewReader(envPath)
	secrets, err := reader.Read()
	if err != nil {
		return fmt.Errorf("read env file: %w", err)
	}

	ss, err := sync.NewSnapshotter(snapshotDir)
	if err != nil {
		return err
	}
	if err := ss.Save(label, secrets); err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "snapshot %q saved (%d keys)\n", label, len(secrets))
	return nil
}

func runSnapshotLoad(cmd *cobra.Command, args []string) error {
	label := args[0]

	ss, err := sync.NewSnapshotter(snapshotDir)
	if err != nil {
		return err
	}
	entry, err := ss.Load(label)
	if err != nil {
		return err
	}
	if entry == nil {
		fmt.Fprintf(os.Stderr, "snapshot %q not found\n", label)
		return fmt.Errorf("snapshot not found: %s", label)
	}
	ss.Print(entry, cmd.OutOrStdout())
	return nil
}
