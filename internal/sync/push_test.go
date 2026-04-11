package sync_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/sync"
	"github.com/your-org/vaultpull/internal/vault"
)

func newPushMockServer(t *testing.T, capturedBody *map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			_ = json.NewDecoder(r.Body).Decode(capturedBody)
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}

func writeEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestPush_WritesSecretsToVault(t *testing.T) {
	var captured map[string]interface{}
	srv := newPushMockServer(t, &captured)
	defer srv.Close()

	client, err := vault.NewClient(srv.URL, "test-token", "v1")
	if err != nil {
		t.Fatal(err)
	}

	envPath := writeEnvFile(t, "FOO=bar\nBAZ=qux\n")
	pusher := sync.NewPusher(client)
	res := pusher.Push(context.Background(), envPath, "secret/myapp")

	if res.Err != nil {
		t.Fatalf("unexpected error: %v", res.Err)
	}
	if res.Written != 2 {
		t.Errorf("expected 2 keys written, got %d", res.Written)
	}
}

func TestPush_EmptyFile_ReturnsError(t *testing.T) {
	var captured map[string]interface{}
	srv := newPushMockServer(t, &captured)
	defer srv.Close()

	client, err := vault.NewClient(srv.URL, "test-token", "v1")
	if err != nil {
		t.Fatal(err)
	}

	envPath := writeEnvFile(t, "")
	pusher := sync.NewPusher(client)
	res := pusher.Push(context.Background(), envPath, "secret/myapp")

	if res.Err == nil {
		t.Fatal("expected error for empty env file")
	}
}

func TestPush_MissingFile_ReturnsError(t *testing.T) {
	client, err := vault.NewClient("http://127.0.0.1:1", "tok", "v1")
	if err != nil {
		t.Fatal(err)
	}
	pusher := sync.NewPusher(client)
	res := pusher.Push(context.Background(), "/nonexistent/.env", "secret/x")
	if res.Err == nil {
		t.Fatal("expected error for missing file")
	}
}
