package sync

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/vault"
)

func newRotateMockServer(t *testing.T) *httptest.Server {
	t.Helper()
	store := map[string]map[string]interface{}{
		"/v1/secret/app": {"API_KEY": "old-value", "DB_PASS": "old-pass"},
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			data, ok := store[r.URL.Path]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": data})
		case http.MethodPut, http.MethodPost:
			var body map[string]interface{}
			_ = json.NewDecoder(r.Body).Decode(&body)
			store[r.URL.Path] = body
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
}

func newRotateClient(t *testing.T, addr string) *vault.Client {
	t.Helper()
	c, err := vault.NewClient(addr, "test-token", 1)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return c
}

func TestRotate_Success(t *testing.T) {
	srv := newRotateMockServer(t)
	defer srv.Close()

	client := newRotateClient(t, srv.URL)
	rotator := NewRotator(client, func() string { return "new-value" })

	res := rotator.Rotate("secret/app", "API_KEY")
	if !res.Success {
		t.Fatalf("expected success, got error: %v", res.Err)
	}
}

func TestRotate_MissingKey_ReturnsError(t *testing.T) {
	srv := newRotateMockServer(t)
	defer srv.Close()

	client := newRotateClient(t, srv.URL)
	rotator := NewRotator(client, nil)

	res := rotator.Rotate("secret/app", "NONEXISTENT")
	if res.Success {
		t.Fatal("expected failure for missing key")
	}
	if !strings.Contains(res.Err.Error(), "not found") {
		t.Errorf("unexpected error: %v", res.Err)
	}
}

func TestRotateAll_CollectsResults(t *testing.T) {
	srv := newRotateMockServer(t)
	defer srv.Close()

	client := newRotateClient(t, srv.URL)
	rotator := NewRotator(client, func() string { return "rotated" })

	results := rotator.RotateAll([]string{"secret/app", "secret/missing"}, "API_KEY")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if !results[0].Success {
		t.Errorf("first result should succeed: %v", results[0].Err)
	}
	if results[1].Success {
		t.Error("second result should fail for missing path")
	}
}

func TestNewRotator_DefaultGenerator(t *testing.T) {
	srv := newRotateMockServer(t)
	defer srv.Close()

	client := newRotateClient(t, srv.URL)
	rotator := NewRotator(client, nil)

	if rotator.generator == nil {
		t.Fatal("expected default generator to be set")
	}
	val := rotator.generator()
	if !strings.HasPrefix(val, "rotated-") {
		t.Errorf("unexpected generated value: %s", val)
	}
}
