package sync

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewExporter_InvalidFormat(t *testing.T) {
	_, err := NewExporter("xml", nil)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestNewExporter_NilWriter_DefaultsToStdout(t *testing.T) {
	ex, err := NewExporter(ExportFormatJSON, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ex.out == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestExporter_WriteJSON(t *testing.T) {
	var buf bytes.Buffer
	ex, _ := NewExporter(ExportFormatJSON, &buf)

	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := ex.Write(secrets); err != nil {
		t.Fatalf("Write: %v", err)
	}

	var got map[string]string
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if got["FOO"] != "bar" || got["BAZ"] != "qux" {
		t.Errorf("unexpected values: %v", got)
	}
}

func TestExporter_WriteDotEnv_Sorted(t *testing.T) {
	var buf bytes.Buffer
	ex, _ := NewExporter(ExportFormatDotEnv, &buf)

	secrets := map[string]string{"ZEBRA": "z", "ALPHA": "a"}
	if err := ex.Write(secrets); err != nil {
		t.Fatalf("Write: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "ALPHA=") {
		t.Errorf("expected ALPHA first, got %q", lines[0])
	}
	if !strings.HasPrefix(lines[1], "ZEBRA=") {
		t.Errorf("expected ZEBRA second, got %q", lines[1])
	}
}

func TestExporter_WriteDotEnv_QuotesValues(t *testing.T) {
	var buf bytes.Buffer
	ex, _ := NewExporter(ExportFormatDotEnv, &buf)

	secrets := map[string]string{"KEY": "hello world"}
	_ = ex.Write(secrets)

	if !strings.Contains(buf.String(), `"hello world"`) {
		t.Errorf("expected quoted value, got: %s", buf.String())
	}
}
