package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newMockVaultServer(t *testing.T, kvVersion int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		if kvVersion == 2 {
			payload = map[string]interface{}{
				"data": map[string]interface{}{
					"data": map[string]interface{}{"DB_PASS": "s3cr3t"},
				},
			}
		} else {
			payload = map[string]interface{}{
				"data": map[string]interface{}{"DB_PASS": "s3cr3t"},
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(payload) //nolint:errcheck
	}))
}

func TestReadSecret_KVv1(t *testing.T) {
	srv := newMockVaultServer(t, 1)
	defer srv.Close()

	c, err := NewClient(srv.URL, "test-token", "", 1)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	data, err := c.ReadSecret("secret/myapp")
	if err != nil {
		t.Fatalf("ReadSecret: %v", err)
	}
	if data["DB_PASS"] != "s3cr3t" {
		t.Errorf("expected DB_PASS=s3cr3t, got %v", data["DB_PASS"])
	}
}

func TestNewClient_InvalidAddr(t *testing.T) {
	_, err := NewClient("://bad-url", "token", "", 1)
	if err == nil {
		t.Fatal("expected error for invalid address")
	}
}
