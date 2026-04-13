package sync

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Checkpoint records the last successful sync state for a given env file,
// allowing incremental or resumed sync operations.
type Checkpoint struct {
	Path      string            `json:"path"`
	Timestamp time.Time         `json:"timestamp"`
	Hashes    map[string]string `json:"hashes"`
}

// CheckpointStore persists and retrieves Checkpoint records on disk.
type CheckpointStore struct {
	dir string
}

// NewCheckpointStore creates a CheckpointStore rooted at dir.
func NewCheckpointStore(dir string) (*CheckpointStore, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("checkpoint: create dir: %w", err)
	}
	return &CheckpointStore{dir: dir}, nil
}

func (s *CheckpointStore) checkpointPath(envPath string) string {
	base := filepath.Base(envPath)
	return filepath.Join(s.dir, base+".checkpoint.json")
}

// Save writes a checkpoint for the given env file path.
func (s *CheckpointStore) Save(envPath string, hashes map[string]string) error {
	cp := Checkpoint{
		Path:      envPath,
		Timestamp: time.Now().UTC(),
		Hashes:    hashes,
	}
	data, err := json.MarshalIndent(cp, "", "  ")
	if err != nil {
		return fmt.Errorf("checkpoint: marshal: %w", err)
	}
	return os.WriteFile(s.checkpointPath(envPath), data, 0o600)
}

// Load retrieves the last checkpoint for the given env file path.
// Returns nil, nil if no checkpoint exists.
func (s *CheckpointStore) Load(envPath string) (*Checkpoint, error) {
	data, err := os.ReadFile(s.checkpointPath(envPath))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("checkpoint: read: %w", err)
	}
	var cp Checkpoint
	if err := json.Unmarshal(data, &cp); err != nil {
		return nil, fmt.Errorf("checkpoint: unmarshal: %w", err)
	}
	return &cp, nil
}

// Delete removes the checkpoint file for the given env file path.
func (s *CheckpointStore) Delete(envPath string) error {
	err := os.Remove(s.checkpointPath(envPath))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
