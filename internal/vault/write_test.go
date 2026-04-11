package vault_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/your-org/vaultpull/internal/vault"
)

func newWriteServer(t *testing.T, wantStatus int, captured *map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, captured)
		w.WriteHeader(wantStatus)
	}))
}

func TestWriteSecret_KVv1_SendsPlainMap(t *testing.T) {
	var captured map[string]interface{}
	srv := newWriteServer(t, http.StatusNoContent, &captured)
	defer srv.Close()

	client, err := vault.NewClient(srv.URL, "tok", "v1")
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]string{"KEY": "value"}
	if err := client.WriteSecret(context.Background(), "secret/app", data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if captured["KEY"] != "value" {
		t.Errorf("expected KEY=value in payload, got %v", captured)
	}
}

func TestWriteSecret_KVv2_WrapsInData(t *testing.T) {
	var captured map[string]interface{}
	srv := newWriteServer(t, http.StatusNoContent, &captured)
	defer srv.Close()

	client, err := vault.NewClient(srv.URL, "tok", "v2")
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]string{"FOO": "bar"}
	if err := client.WriteSecret(context.Background(), "secret/app", data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	inner, ok := captured["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected 'data' wrapper, got %v", captured)
	}
	if inner["FOO"] != "bar" {
		t.Errorf("expected FOO=bar inside data, got %v", inner)
	}
}

func TestWriteSecret_NonSuccessStatus_ReturnsError(t *testing.T) {
	var captured map[string]interface{}
	srv := newWriteServer(t, http.StatusForbidden, &captured)
	defer srv.Close()

	client, err := vault.NewClient(srv.URL, "bad-token", "v1")
	if err != nil {
		t.Fatal(err)
	}

	err = client.WriteSecret(context.Background(), "secret/app", map[string]string{"K": "v"})
	if err == nil {
		t.Fatal("expected error for 403 response")
	}
}
