package vault

import (
	"fmt"
	"net/http"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with project-specific configuration.
type Client struct {
	api       *vaultapi.Client
	kvVersion int
	namespace string
}

// NewClient creates and configures a new Vault client.
func NewClient(addr, token, namespace string, kvVersion int) (*Client, error) {
	cfg := vaultapi.DefaultConfig()
	cfg.Address = addr
	cfg.HttpClient = &http.Client{Timeout: 10 * time.Second}

	api, err := vaultapi.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault api client: %w", err)
	}

	api.SetToken(token)

	if namespace != "" {
		api.SetNamespace(namespace)
	}

	return &Client{
		api:       api,
		kvVersion: kvVersion,
		namespace: namespace,
	}, nil
}

// ReadSecret reads a secret at the given path and returns its data map.
func (c *Client) ReadSecret(path string) (map[string]interface{}, error) {
	var secret *vaultapi.Secret
	var err error

	switch c.kvVersion {
	case 2:
		secret, err = c.api.Logical().Read(kvV2DataPath(path))
	default:
		secret, err = c.api.Logical().Read(path)
	}

	if err != nil {
		return nil, fmt.Errorf("reading secret at %q: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secret found at %q", path)
	}

	if c.kvVersion == 2 {
		data, ok := secret.Data["data"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected kv v2 response shape at %q", path)
		}
		return data, nil
	}

	return secret.Data, nil
}

// kvV2DataPath converts a mount/key path to the kv v2 data path.
func kvV2DataPath(path string) string {
	return insertAfterMount(path, "data")
}
