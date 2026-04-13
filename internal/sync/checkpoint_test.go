package sync

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewCheckpointStore_CreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "checkpoints")
	_, err := NewCheckpointStore(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("expected dir to exist: %v", err)
	}
}

func TestCheckpointStore_SaveAndLoad(t *testing.T) {
	store, _ := NewCheckpointStore(t.TempDir())
	hashes := map[string]string{"DB_PASS": "abc123", "API_KEY": "def456"}

	envPath := "/app/.env"
	if err := store.Save(envPath, hashes); err != nil {
		t.Fatalf("Save: %v", err)
	}

	cp, err := store.Load(envPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cp == nil {
		t.Fatal("expected non-nil checkpoint")
	}
	if cp.Path != envPath {
		t.Errorf("path: got %q, want %q", cp.Path, envPath)
	}
	if cp.Hashes["DB_PASS"] != "abc123" {
		t.Errorf("hash mismatch for DB_PASS")
	}
	if cp.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestCheckpointStore_Load_MissingFile_ReturnsNil(t *testing.T) {
	store, _ := NewCheckpointStore(t.TempDir())
	cp, err := store.Load("/nonexistent/.env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cp != nil {
		t.Errorf("expected nil checkpoint, got %+v", cp)
	}
}

func TestCheckpointStore_Delete_RemovesFile(t *testing.T) {
	store, _ := NewCheckpointStore(t.TempDir())
	envPath := "/app/.env"
	_ = store.Save(envPath, map[string]string{"K": "v"})

	if err := store.Delete(envPath); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	cp, err := store.Load(envPath)
	if err != nil {
		t.Fatalf("Load after delete: %v", err)
	}
	if cp != nil {
		t.Error("expected nil after delete")
	}
}

func TestCheckpointStore_Delete_MissingFile_NoError(t *testing.T) {
	store, _ := NewCheckpointStore(t.TempDir())
	if err := store.Delete("/ghost/.env"); err != nil {
		t.Errorf("expected no error for missing file, got %v", err)
	}
}

func TestCheckpoint_TimestampIsUTC(t *testing.T) {
	store, _ := NewCheckpointStore(t.TempDir())
	_ = store.Save("/app/.env", map[string]string{})
	cp, _ := store.Load("/app/.env")
	if cp.Timestamp.Location() != time.UTC {
		t.Errorf("expected UTC timestamp, got %v", cp.Timestamp.Location())
	}
}
