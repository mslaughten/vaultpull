package sync

import (
	"fmt"
	"io"
	"os"
	"time"
)

// AuditEntry records the outcome of a single secret sync operation.
type AuditEntry struct {
	Timestamp time.Time
	Path      string
	File      string
	Added     int
	Removed   int
	Changed   int
	Unchanged int
	Err       error
}

// AuditLog collects audit entries produced during a sync run.
type AuditLog struct {
	w io.Writer
}

// NewAuditLog returns an AuditLog that writes to w.
// Pass nil to discard all output.
func NewAuditLog(w io.Writer) *AuditLog {
	if w == nil {
		w = io.Discard
	}
	return &AuditLog{w: w}
}

// NewAuditLogToFile opens (or creates) the file at path and returns an
// AuditLog backed by it. The caller is responsible for closing the file.
func NewAuditLogToFile(path string) (*AuditLog, *os.File, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, nil, fmt.Errorf("audit: open %s: %w", path, err)
	}
	return NewAuditLog(f), f, nil
}

// Record writes a single entry to the underlying writer.
func (a *AuditLog) Record(e AuditEntry) {
	status := "ok"
	if e.Err != nil {
		status = fmt.Sprintf("error: %v", e.Err)
	}
	fmt.Fprintf(
		a.w,
		"%s path=%q file=%q added=%d removed=%d changed=%d unchanged=%d status=%s\n",
		e.Timestamp.UTC().Format(time.RFC3339),
		e.Path,
		e.File,
		e.Added,
		e.Removed,
		e.Changed,
		e.Unchanged,
		status,
	)
}

// EntryFromDiff builds an AuditEntry from a Diff and associated metadata.
func EntryFromDiff(path, file string, d Diff, err error) AuditEntry {
	return AuditEntry{
		Timestamp: time.Now(),
		Path:      path,
		File:      file,
		Added:     len(d.Added),
		Removed:   len(d.Removed),
		Changed:   len(d.Changed),
		Unchanged: len(d.Unchanged),
		Err:       err,
	}
}
