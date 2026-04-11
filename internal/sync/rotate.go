package sync

import (
	"fmt"
	"time"

	"github.com/your-org/vaultpull/internal/vault"
)

// RotateResult holds the outcome of a single secret rotation.
type RotateResult struct {
	Path    string
	Key     string
	Success bool
	Err     error
}

// Rotator generates new values for secrets and writes them back to Vault.
type Rotator struct {
	client    *vault.Client
	generator func() string
}

// NewRotator creates a Rotator using the provided Vault client.
// If generator is nil, a default timestamp-based generator is used.
func NewRotator(client *vault.Client, generator func() string) *Rotator {
	if generator == nil {
		generator = func() string {
			return fmt.Sprintf("rotated-%d", time.Now().UnixNano())
		}
	}
	return &Rotator{client: client, generator: generator}
}

// Rotate replaces the value of the given key at the given Vault path
// with a newly generated value and returns the result.
func (r *Rotator) Rotate(path, key string) RotateResult {
	secrets, err := r.client.ReadSecret(path)
	if err != nil {
		return RotateResult{Path: path, Key: key, Success: false, Err: fmt.Errorf("read: %w", err)}
	}

	if _, ok := secrets[key]; !ok {
		return RotateResult{Path: path, Key: key, Success: false, Err: fmt.Errorf("key %q not found at %s", key, path)}
	}

	secrets[key] = r.generator()

	if err := r.client.WriteSecret(path, secrets); err != nil {
		return RotateResult{Path: path, Key: key, Success: false, Err: fmt.Errorf("write: %w", err)}
	}

	return RotateResult{Path: path, Key: key, Success: true}
}

// RotateAll rotates the specified key across multiple paths, collecting results.
func (r *Rotator) RotateAll(paths []string, key string) []RotateResult {
	results := make([]RotateResult, 0, len(paths))
	for _, p := range paths {
		results = append(results, r.Rotate(p, key))
	}
	return results
}
