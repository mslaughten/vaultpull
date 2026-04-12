package sync

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewSecretCache_CreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "cache")
	c, err := NewSecretCache(dir, time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil cache")
	}
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("dir not created: %v", err)
	}
}

func TestSecretCache_SetAndGet(t *testing.T) {
	c, _ := NewSecretCache(t.TempDir(), time.Minute)
	secrets := map[string]string{"KEY": "value", "OTHER": "data"}

	if err := c.Set("secret/app", secrets); err != nil {
		t.Fatalf("Set: %v", err)
	}

	entry, ok := c.Get("secret/app")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if entry.Secrets["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %s", entry.Secrets["KEY"])
	}
	if entry.Checksum == "" {
		t.Error("expected non-empty checksum")
	}
}

func TestSecretCache_Get_MissingEntry_ReturnsFalse(t *testing.T) {
	c, _ := NewSecretCache(t.TempDir(), time.Minute)
	_, ok := c.Get("secret/missing")
	if ok {
		t.Error("expected cache miss")
	}
}

func TestSecretCache_Get_ExpiredEntry_ReturnsFalse(t *testing.T) {
	c, _ := NewSecretCache(t.TempDir(), -time.Second) // already expired
	_ = c.Set("secret/app", map[string]string{"A": "1"})
	_, ok := c.Get("secret/app")
	if ok {
		t.Error("expected expired cache to return false")
	}
}

func TestSecretCache_Invalidate_RemovesEntry(t *testing.T) {
	c, _ := NewSecretCache(t.TempDir(), time.Minute)
	_ = c.Set("secret/app", map[string]string{"X": "y"})

	if err := c.Invalidate("secret/app"); err != nil {
		t.Fatalf("Invalidate: %v", err)
	}
	_, ok := c.Get("secret/app")
	if ok {
		t.Error("expected cache miss after invalidation")
	}
}

func TestSecretCache_Invalidate_MissingEntry_NoError(t *testing.T) {
	c, _ := NewSecretCache(t.TempDir(), time.Minute)
	if err := c.Invalidate("secret/nonexistent"); err != nil {
		t.Errorf("expected no error for missing entry, got %v", err)
	}
}

func TestChecksumSecrets_Deterministic(t *testing.T) {
	s := map[string]string{"A": "1", "B": "2"}
	a := checksumSecrets(s)
	b := checksumSecrets(s)
	if a != b {
		t.Errorf("checksum not deterministic: %s vs %s", a, b)
	}
}
