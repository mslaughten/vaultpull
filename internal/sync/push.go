package sync

import (
	"context"
	"fmt"

	"github.com/your-org/vaultpull/internal/envfile"
	"github.com/your-org/vaultpull/internal/vault"
)

// Pusher writes local .env file secrets back into Vault.
type Pusher struct {
	client *vault.Client
	reader *envfile.Reader
}

// NewPusher creates a Pusher using the provided Vault client.
func NewPusher(client *vault.Client) *Pusher {
	return &Pusher{
		client: client,
		reader: envfile.NewReader(),
	}
}

// PushResult holds the outcome of a push operation for a single path.
type PushResult struct {
	VaultPath string
	EnvFile   string
	Written   int
	Err       error
}

// Push reads the given env file and writes all key/value pairs to the
// specified Vault secret path. Returns a PushResult describing the outcome.
func (p *Pusher) Push(ctx context.Context, envFile, vaultPath string) PushResult {
	res := PushResult{VaultPath: vaultPath, EnvFile: envFile}

	data, err := p.reader.Read(envFile)
	if err != nil {
		res.Err = fmt.Errorf("reading %s: %w", envFile, err)
		return res
	}

	if len(data) == 0 {
		res.Err = fmt.Errorf("env file %s is empty, nothing to push", envFile)
		return res
	}

	if err := p.client.WriteSecret(ctx, vaultPath, data); err != nil {
		res.Err = fmt.Errorf("writing to vault path %s: %w", vaultPath, err)
		return res
	}

	res.Written = len(data)
	return res
}
