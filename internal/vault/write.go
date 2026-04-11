package vault

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// WriteSecret writes the provided key/value map to the given Vault secret path.
// For KV v2 the payload is wrapped under a "data" key as required by the API.
func (c *Client) WriteSecret(ctx context.Context, path string, data map[string]string) error {
	payload, err := c.buildWritePayload(data)
	if err != nil {
		return fmt.Errorf("building payload: %w", err)
	}

	url := fmt.Sprintf("%s/v1/%s", c.baseURL, c.writeDataPath(path))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("vault returned status %d for path %s", resp.StatusCode, path)
	}
	return nil
}

func (c *Client) buildWritePayload(data map[string]string) ([]byte, error) {
	var body interface{}
	if c.kvVersion == "2" {
		body = map[string]interface{}{"data": data}
	} else {
		body = data
	}
	return json.Marshal(body)
}

func (c *Client) writeDataPath(path string) string {
	if c.kvVersion == "2" {
		return kvV2DataPath(path)
	}
	return path
}
