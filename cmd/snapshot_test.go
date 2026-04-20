package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestSnapshotCmd_RegisteredOnRoot(t *testing.T) {
	var found bool
	for _, c := range rootCmd.Commands() {
		if c.Name() == "snapshot" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("snapshot command not registered on root")
	}
}

func TestSnapshotSaveCmd_RequiresTwoArgs(t *testing.T) {
	rootCmd.SetArgs([]string{"snapshot", "save", "only-one"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing second arg")
	}
}

func TestSnapshotLoadCmd_RequiresOneArg(t *testing.T) {
	rootCmd.SetArgs([]string{"snapshot", "load"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing label arg")
	}
}

func TestSnapshotSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, "test.env")
	snapDir := filepath.Join(dir, "snaps")

	_ = os.WriteFile(envPath, []byte("FOO=bar\nBAZ=qux\n"), 0o600)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"snapshot", "save", "--dir", snapDir, "mysnap", envPath})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("save: %v", err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"snapshot", "load", "--dir", snapDir, "mysnap"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("load: %v", err)
	}
	out := buf.String()
	if out == "" {
		t.Fatal("expected output from load")
	}
}

func TestSnapshotLoad_Missing_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	rootCmd.SetArgs([]string{"snapshot", "load", "--dir", dir, "ghost"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}
