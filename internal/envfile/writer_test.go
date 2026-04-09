package envfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWrite_CreatesFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, ".env")

	writer := NewWriter(filePath)
	secrets := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"API_KEY":      "secret123",
	}

	err := writer.Write(secrets)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	if !writer.Exists() {
		t.Error("Expected file to exist after Write()")
	}
}

func TestWrite_SortsKeys(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, ".env")

	writer := NewWriter(filePath)
	secrets := map[string]string{
		"ZEBRA": "last",
		"ALPHA": "first",
		"BETA":  "second",
	}

	err := writer.Write(secrets)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	// Skip header comments
	var dataLines []string
	for _, line := range lines {
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		dataLines = append(dataLines, line)
	}

	if len(dataLines) != 3 {
		t.Fatalf("Expected 3 data lines, got %d", len(dataLines))
	}

	if !strings.HasPrefix(dataLines[0], "ALPHA=") {
		t.Errorf("Expected first key to be ALPHA, got %s", dataLines[0])
	}
	if !strings.HasPrefix(dataLines[1], "BETA=") {
		t.Errorf("Expected second key to be BETA, got %s", dataLines[1])
	}
	if !strings.HasPrefix(dataLines[2], "ZEBRA=") {
		t.Errorf("Expected third key to be ZEBRA, got %s", dataLines[2])
	}
}

func TestWrite_EscapesSpecialCharacters(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, ".env")

	writer := NewWriter(filePath)
	secrets := map[string]string{
		"KEY_WITH_QUOTES": `value with "quotes"`,
		"KEY_WITH_NEWLINE": "line1\nline2",
	}

	err := writer.Write(secrets)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if !strings.Contains(string(content), `\\"`) {
		t.Error("Expected escaped quotes in output")
	}
}
