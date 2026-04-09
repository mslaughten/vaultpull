package vault

import (
	"context"
	"fmt"
	"strings"
)

// SecretMap is a flat map of key-value secret pairs.
type SecretMap map[string]string

// ListSecrets returns all secret keys under the given mount and path prefix.
func (c *Client) ListSecrets(ctx context.Context, mount, pathPrefix string) ([]string, error) {
	listPath := pathPrefix
	if c.kvVersion == 2 {
		listPath = fmt.Sprintf("%s/metadata/%s", mount, strings.TrimPrefix(pathPrefix, "/"))
	} else {
		listPath = fmt.Sprintf("%s/%s", mount, strings.TrimPrefix(pathPrefix, "/"))
	}

	secret, err := c.logical.ListWithContext(ctx, listPath)
	if err != nil {
		return nil, fmt.Errorf("listing secrets at %q: %w", listPath, err)
	}
	if secret == nil || secret.Data == nil {
		return []string{}, nil
	}

	raw, ok := secret.Data["keys"]
	if !ok {
		return []string{}, nil
	}

	ifaces, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for keys field")
	}

	keys := make([]string, 0, len(ifaces))
	for _, v := range ifaces {
		if s, ok := v.(string); ok {
			keys = append(keys, s)
		}
	}
	return keys, nil
}

// ReadSecretMap reads a single secret path and returns its key-value data.
func (c *Client) ReadSecretMap(ctx context.Context, mount, secretPath string) (SecretMap, error) {
	result := make(SecretMap)

	data, err := c.ReadSecret(ctx, mount, secretPath)
	if err != nil {
		return nil, err
	}

	for k, v := range data {
		if str, ok := v.(string); ok {
			result[k] = str
		} else {
			result[k] = fmt.Sprintf("%v", v)
		}
	}
	return result, nil
}
