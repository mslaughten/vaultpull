package sync

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTakeSnapshot_ExistingFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	original := []byte("KEY=value\n")
	if err := os.WriteFile(p, original, 0o600); err != nil {
		t.Fatal(err)
	}

	snap, err := TakeSnapshot(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !snap.exists {
		t.Error("expected exists=true")
	}
	if string(snap.Content) != string(original) {
		t.Errorf("content mismatch: got %q want %q", snap.Content, original)
	}
}

func TestTakeSnapshot_MissingFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")

	snap, err := TakeSnapshot(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.exists {
		t.Error("expected exists=false for missing file")
	}
}

func TestRestore_WritesContentBack(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	original := []byte("FOO=bar\n")
	if err := os.WriteFile(p, original, 0o600); err != nil {
		t.Fatal(err)
	}

	snap, _ := TakeSnapshot(p)

	// Overwrite the file to simulate a sync.
	_ = os.WriteFile(p, []byte("FOO=changed\n"), 0o600)

	if err := snap.Restore(); err != nil {
		t.Fatalf("Restore error: %v", err)
	}
	got, _ := os.ReadFile(p)
	if string(got) != string(original) {
		t.Errorf("got %q, want %q", got, original)
	}
}

func TestRestore_RemovesFileIfDidNotExist(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")

	snap, _ := TakeSnapshot(p) // file doesn't exist yet

	// Create the file to simulate a sync write.
	_ = os.WriteFile(p, []byte("NEW=1\n"), 0o600)

	if err := snap.Restore(); err != nil {
		t.Fatalf("Restore error: %v", err)
	}
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		t.Error("expected file to be removed after rollback")
	}
}

func TestRollbackAll_CollectsErrors(t *testing.T) {
	// Provide a snapshot pointing at a path inside a non-existent directory
	// so WriteFile will fail.
	s := &Snapshot{
		Path:    "/nonexistent/dir/.env",
		Content: []byte("X=1"),
		exists:  true,
	}
	errs := RollbackAll([]*Snapshot{s})
	if len(errs) == 0 {
		t.Error("expected at least one error from RollbackAll")
	}
}
