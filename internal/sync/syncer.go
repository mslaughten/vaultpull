// Package sync orchestrates pulling secrets from Vault and writing them to .env files.
package sync

import (
	"fmt"
	"path/filepath"

	"github.com/yourusername/vaultpull/internal/envfile"
	"github.com/yourusername/vaultpull/internal/vault"
)

// Options controls the behaviour of a sync run.
type Options struct {
	MountPath string
	Namespace string
	OutputDir string
	DryRun    bool
}

// Syncer pulls secrets from Vault and writes them to .env files.
type Syncer struct {
	client *vault.Client
	opts   Options
}

// New creates a Syncer with the provided Vault client and options.
func New(client *vault.Client, opts Options) *Syncer {
	return &Syncer{client: client, opts: opts}
}

// Run executes the full sync: list → filter → read → write.
func (s *Syncer) Run() (Result, error) {
	paths, err := s.client.ListSecrets(s.opts.MountPath)
	if err != nil {
		return Result{}, fmt.Errorf("listing secrets: %w", err)
	}

	if s.opts.Namespace != "" {
		paths = vault.FilterByNamespace(paths, s.opts.Namespace)
	}

	result := Result{Total: len(paths)}

	for _, p := range paths {
		secrets, err := s.client.ReadSecretMap(p)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("read %s: %w", p, err))
			continue
		}

		stripped := vault.StripNamespace(p, s.opts.Namespace)
		outPath := filepath.Join(s.opts.OutputDir, envFileName(stripped))

		if !s.opts.DryRun {
			w := envfile.NewWriter(outPath)
			if err := w.Write(secrets); err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("write %s: %w", outPath, err))
				continue
			}
		}

		result.Written = append(result.Written, outPath)
	}

	return result, nil
}

// envFileName converts a secret path like "app/prod" to "app_prod.env".
func envFileName(path string) string {
	safe := ""
	for _, ch := range path {
		if ch == '/' {
			safe += "_"
		} else {
			safe += string(ch)
		}
	}
	if safe == "" {
		safe = "secrets"
	}
	return safe + ".env"
}
