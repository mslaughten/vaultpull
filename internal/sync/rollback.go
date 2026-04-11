package sync

import (
	"fmt"
	"os"
	"path/filepath"
)

// Snapshot holds a pre-sync backup of an env file's contents.
type Snapshot struct {
	Path    string
	Content []byte
	exists  bool
}

// TakeSnapshot reads the current contents of path (if it exists) so it can
// be restored later via Restore.
func TakeSnapshot(path string) (*Snapshot, error) {
	s := &Snapshot{Path: path}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		s.exists = false
		return s, nil
	}
	if err != nil {
		return nil, fmt.Errorf("snapshot %q: %w", path, err)
	}
	s.exists = true
	s.Content = data
	return s, nil
}

// Restore writes the snapshot content back to disk, or removes the file if it
// did not exist before the snapshot was taken.
func (s *Snapshot) Restore() error {
	if !s.exists {
		err := os.Remove(s.Path)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("rollback remove %q: %w", s.Path, err)
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(s.Path), 0o755); err != nil {
		return fmt.Errorf("rollback mkdir %q: %w", s.Path, err)
	}
	if err := os.WriteFile(s.Path, s.Content, 0o600); err != nil {
		return fmt.Errorf("rollback write %q: %w", s.Path, err)
	}
	return nil
}

// RollbackAll restores every snapshot in the slice, collecting all errors.
func RollbackAll(snapshots []*Snapshot) []error {
	var errs []error
	for _, s := range snapshots {
		if err := s.Restore(); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
