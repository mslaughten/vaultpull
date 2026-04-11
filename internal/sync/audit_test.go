package sync

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func fixedEntry(err error) AuditEntry {
	return AuditEntry{
		Timestamp: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
		Path:      "secret/app/prod",
		File:      ".env.prod",
		Added:     2,
		Removed:   1,
		Changed:   3,
		Unchanged: 5,
		Err:       err,
	}
}

func TestAuditLog_Record_Success(t *testing.T) {
	var buf bytes.Buffer
	log := NewAuditLog(&buf)
	log.Record(fixedEntry(nil))

	line := buf.String()
	for _, want := range []string{
		"2024-01-15T12:00:00Z",
		`path="secret/app/prod"`,
		`file=".env.prod"`,
		"added=2",
		"removed=1",
		"changed=3",
		"unchanged=5",
		"status=ok",
	} {
		if !strings.Contains(line, want) {
			t.Errorf("expected %q in output %q", want, line)
		}
	}
}

func TestAuditLog_Record_Error(t *testing.T) {
	var buf bytes.Buffer
	log := NewAuditLog(&buf)
	log.Record(fixedEntry(errors.New("vault unreachable")))

	if !strings.Contains(buf.String(), "error: vault unreachable") {
		t.Errorf("expected error message in output, got: %s", buf.String())
	}
}

func TestNewAuditLog_NilWriter_Discards(t *testing.T) {
	log := NewAuditLog(nil)
	// should not panic
	log.Record(fixedEntry(nil))
}

func TestNewAuditLogToFile_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	log, f, err := NewAuditLogToFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer f.Close()

	log.Record(fixedEntry(nil))
	f.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read audit file: %v", err)
	}
	if !strings.Contains(string(data), "status=ok") {
		t.Errorf("audit file missing expected content: %s", data)
	}
}

func TestEntryFromDiff_PopulatesFields(t *testing.T) {
	d := Diff{
		Added:     map[string]string{"A": "1", "B": "2"},
		Removed:   map[string]string{"C": "3"},
		Changed:   map[string]DiffChange{"D": {}},
		Unchanged: map[string]string{"E": "5", "F": "6", "G": "7"},
	}
	e := EntryFromDiff("secret/svc", ".env", d, nil)

	if e.Added != 2 || e.Removed != 1 || e.Changed != 1 || e.Unchanged != 3 {
		t.Errorf("counts mismatch: %+v", e)
	}
	if e.Path != "secret/svc" || e.File != ".env" {
		t.Errorf("metadata mismatch: %+v", e)
	}
	if e.Err != nil {
		t.Errorf("expected nil error, got %v", e.Err)
	}
}
