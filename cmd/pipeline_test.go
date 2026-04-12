package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempJSON(t *testing.T, data map[string]string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "secrets-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(data); err != nil {
		t.Fatalf("encode JSON: %v", err)
	}
	return f.Name()
}

func TestPipelineCmd_RegisteredOnRoot(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use == "pipeline <path>" {
			return
		}
	}
	t.Fatal("pipeline command not registered on root")
}

func TestPipelineCmd_RequiresOneArg(t *testing.T) {
	input := writeTempJSON(t, map[string]string{"K": "v"})
	rootCmd.SetArgs([]string{"pipeline", "--input", input})
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing positional arg")
	}
}

func TestPipelineCmd_DefaultFlags(t *testing.T) {
	for _, sub := range rootCmd.Commands() {
		if sub.Use != "pipeline <path>" {
			continue
		}
		if f := sub.Flags().Lookup("stage"); f == nil {
			t.Error("missing --stage flag")
		}
		if f := sub.Flags().Lookup("input"); f == nil {
			t.Error("missing --input flag")
		}
		if f := sub.Flags().Lookup("output"); f == nil {
			t.Error("missing --output flag")
		} else if f.DefValue != "dotenv" {
			t.Errorf("expected default output=dotenv, got %s", f.DefValue)
		}
		return
	}
	t.Fatal("pipeline command not found")
}

func TestPipelineCmd_UnknownStage_ReturnsError(t *testing.T) {
	input := writeTempJSON(t, map[string]string{"K": "v"})
	rootCmd.SetArgs([]string{"pipeline", "secret/app", "--input", input, "--stage", "nonexistent"})
	var buf bytes.Buffer
	rootCmd.SetErr(&buf)
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for unknown stage")
	}
	if !strings.Contains(err.Error(), "nonexistent") {
		t.Errorf("error should mention stage name, got: %v", err)
	}
}

func TestPipelineCmd_MissingInputFile_ReturnsError(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope.json")
	rootCmd.SetArgs([]string{"pipeline", "secret/app", "--input", missing})
	var buf bytes.Buffer
	rootCmd.SetErr(&buf)
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing input file")
	}
}
