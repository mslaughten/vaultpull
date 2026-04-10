package sync_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultpull/internal/config"
	"github.com/yourusername/vaultpull/internal/sync"
	"github.com/yourusername/vaultpull/internal/vault"
)

func newMockServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	// LIST secrets/metadata/myapp/
	mux.HandleFunc("/v1/secret/metadata/myapp/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"keys":["db"]}}`))
	})

	// GET secret/data/myapp/db
	mux.HandleFunc("/v1/secret/data/myapp/db", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"data":{"DB_HOST":"localhost","DB_PORT":"5432"}}}`))
	})

	return httptest.NewServer(mux)
}

func TestRun_WritesEnvFile(t *testing.T) {
	srv := newMockServer(t)
	defer srv.Close()

	cfg := config.Config{
		VaultAddr:  srv.URL,
		VaultToken: "test-token",
		KVVersion:  2,
		MountPath:  "secret",
	}
	client, err := vault.NewClient(cfg)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	dir := t.TempDir()
	s := sync.New(client, sync.Options{
		MountPath: "secret",
		Namespace: "myapp",
		OutputDir: dir,
	})

	result, err := s.Run()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.HasErrors() {
		t.Fatalf("unexpected errors: %s", result.ErrorMessages())
	}
	if len(result.Written) != 1 {
		t.Fatalf("expected 1 written file, got %d", len(result.Written))
	}

	outPath := filepath.Join(dir, "db.env")
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		t.Errorf("expected %s to exist", outPath)
	}
}

func TestRun_DryRun_NoFiles(t *testing.T) {
	srv := newMockServer(t)
	defer srv.Close()

	cfg := config.Config{
		VaultAddr:  srv.URL,
		VaultToken: "test-token",
		KVVersion:  2,
		MountPath:  "secret",
	}
	client, _ := vault.NewClient(cfg)

	dir := t.TempDir()
	s := sync.New(client, sync.Options{
		MountPath: "secret",
		Namespace: "myapp",
		OutputDir: dir,
		DryRun:    true,
	})

	result, err := s.Run()
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 0 {
		t.Errorf("dry-run should not create files, found %d", len(entries))
	}
	if len(result.Written) != 1 {
		t.Errorf("Written should still record paths in dry-run mode")
	}
}
