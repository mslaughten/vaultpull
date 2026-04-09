package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListSecrets_KVv1(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/secret/myapp":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"keys": []interface{}{"db_password", "api_key"},
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, "test-token", "secret", 1)
	require.NoError(t, err)

	keys, err := client.ListSecrets(context.Background(), "secret", "myapp")
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"db_password", "api_key"}, keys)
}

func TestListSecrets_EmptyPath(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{},
		})
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL, "test-token", "secret", 1)
	require.NoError(t, err)

	keys, err := client.ListSecrets(context.Background(), "secret", "empty")
	require.NoError(t, err)
	assert.Empty(t, keys)
}

func TestReadSecretMap_ReturnsStringMap(t *testing.T) {
	ts := newMockVaultServer(t)
	defer ts.Close()

	client, err := NewClient(ts.URL, "test-token", "secret", 1)
	require.NoError(t, err)

	result, err := client.ReadSecretMap(context.Background(), "secret", "myapp/config")
	require.NoError(t, err)
	assert.IsType(t, SecretMap{}, result)
	for k, v := range result {
		assert.IsType(t, "", k)
		assert.IsType(t, "", v)
	}
}
