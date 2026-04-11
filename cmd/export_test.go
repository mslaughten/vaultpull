package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestExportCmd_RegisteredOnRoot(t *testing.T) {
	var found bool
	for _, sub := range RootCmd.Commands() {
		if sub.Use == "export <secret-path>" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("export command not registered on root")
	}
}

func TestExportCmd_DefaultFlags(t *testing.T) {
	f := exportCmd.Flags()

	format, err := f.GetString("format")
	if err != nil || format != "dotenv" {
		t.Errorf("expected default format=dotenv, got %q err=%v", format, err)
	}

	mount, err := f.GetString("mount")
	if err != nil || mount != "secret" {
		t.Errorf("expected default mount=secret, got %q err=%v", mount, err)
	}

	ns, err := f.GetString("namespace")
	if err != nil || ns != "" {
		t.Errorf("expected default namespace empty, got %q err=%v", ns, err)
	}
}

func TestExportCmd_RequiresOneArg(t *testing.T) {
	var buf bytes.Buffer
	RootCmd.SetOut(&buf)
	RootCmd.SetErr(&buf)
	RootCmd.SetArgs([]string{"export"})

	err := RootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when no args provided")
	}
	if !strings.Contains(err.Error(), "accepts 1 arg") {
		t.Errorf("unexpected error: %v", err)
	}
}
