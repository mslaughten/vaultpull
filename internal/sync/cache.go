package sync

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// CacheEntry holds a cached snapshot of secrets for a given Vault path.
type CacheEntry struct {
	Path      string            `json:"path"`
	Secrets   map[string]string `json:"secrets"`
	FetchedAt time.Time         `json:"fetched_at"`
	Checksum  string            `json:"checksum"`
}

// SecretCache persists fetched secrets to disk to avoid redundant Vault reads.
type SecretCache struct {
	dir string
	ttl time.Duration
}

// NewSecretCache creates a cache backed by the given directory with a TTL.
func NewSecretCache(dir string, ttl time.Duration) (*SecretCache, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("cache: mkdir %s: %w", dir, err)
	}
	return &SecretCache{dir: dir, ttl: ttl}, nil
}

// Get returns a cached entry for path if it exists and has not expired.
func (c *SecretCache) Get(path string) (*CacheEntry, bool) {
	p := c.filePath(path)
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, false
	}
	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false
	}
	if time.Since(entry.FetchedAt) > c.ttl {
		_ = os.Remove(p)
		return nil, false
	}
	return &entry, true
}

// Set writes a cache entry for path to disk.
func (c *SecretCache) Set(path string, secrets map[string]string) error {
	entry := CacheEntry{
		Path:      path,
		Secrets:   secrets,
		FetchedAt: time.Now().UTC(),
		Checksum:  checksumSecrets(secrets),
	}
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("cache: marshal: %w", err)
	}
	return os.WriteFile(c.filePath(path), data, 0o600)
}

// Invalidate removes the cached entry for path.
func (c *SecretCache) Invalidate(path string) error {
	err := os.Remove(c.filePath(path))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (c *SecretCache) filePath(path string) string {
	h := sha256.Sum256([]byte(path))
	name := hex.EncodeToString(h[:]) + ".json"
	return filepath.Join(c.dir, name)
}

func checksumSecrets(secrets map[string]string) string {
	b, _ := json.Marshal(secrets)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
