package sync

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// SnapshotEntry holds a captured key-value map with metadata.
type SnapshotEntry struct {
	Label     string            `json:"label"`
	CreatedAt time.Time         `json:"created_at"`
	Secrets   map[string]string `json:"secrets"`
}

// Snapshotter captures and restores named snapshots of secret maps.
type Snapshotter struct {
	dir string
}

// NewSnapshotter creates a Snapshotter that persists snapshots under dir.
func NewSnapshotter(dir string) (*Snapshotter, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("snapshot: create dir: %w", err)
	}
	return &Snapshotter{dir: dir}, nil
}

// Save writes a labeled snapshot of secrets to disk.
func (s *Snapshotter) Save(label string, secrets map[string]string) error {
	entry := SnapshotEntry{
		Label:     label,
		CreatedAt: time.Now().UTC(),
		Secrets:   secrets,
	}
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	path := filepath.Join(s.dir, label+".json")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("snapshot: write %s: %w", path, err)
	}
	return nil
}

// Load reads a named snapshot from disk. Returns nil, nil if not found.
func (s *Snapshotter) Load(label string) (*SnapshotEntry, error) {
	path := filepath.Join(s.dir, label+".json")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("snapshot: read %s: %w", path, err)
	}
	var entry SnapshotEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &entry, nil
}

// Delete removes a named snapshot. Missing snapshots are silently ignored.
func (s *Snapshotter) Delete(label string) error {
	path := filepath.Join(s.dir, label+".json")
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("snapshot: delete %s: %w", path, err)
	}
	return nil
}

// Print writes a human-readable summary of a snapshot to w.
func (s *Snapshotter) Print(entry *SnapshotEntry, w io.Writer) {
	fmt.Fprintf(w, "snapshot: %s (captured %s)\n", entry.Label, entry.CreatedAt.Format(time.RFC3339))
	for k, v := range entry.Secrets {
		fmt.Fprintf(w, "  %s=%s\n", k, v)
	}
}
