package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
	"time"

	syncp "github.com/vaultpull/internal/sync"
)

func TestCacheCmd_RegisteredOnRoot(t *testing.T) {
	found := false
	for _, c := range rootCmd.Commands() {
		if c.Name() == "cache" {
			found = true
			break
		}
	}
	if !found {
		t.Error("cache command not registered on root")
	}
}

func TestCacheGetCmd_RequiresOneArg(t *testing.T) {
	rootCmd.SetArgs([]string{"cache", "get"})
	if err := rootCmd.Execute(); err == nil {
		t.Error("expected error for missing argument")
	}
}

func TestCacheGetCmd_CacheMiss(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer
	cacheGetCmd.SetOut(&buf)
	cacheGetCmd.SetArgs([]string{"--cache-dir", dir, "--ttl", "1m"})

	args := []string{"secret/missing"}
	if err := cacheGetCmd.RunE(cacheGetCmd, args); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "cache miss") {
		t.Errorf("expected 'cache miss', got %q", buf.String())
	}
}

func TestCacheGetCmd_CacheHit_PrintsSecrets(t *testing.T) {
	dir := t.TempDir()
	c, _ := syncp.NewSecretCache(dir, time.Minute)
	_ = c.Set("secret/app", map[string]string{"TOKEN": "abc123"})

	var buf bytes.Buffer
	cacheGetCmd.SetOut(&buf)

	if err := cacheGetCmd.RunE(cacheGetCmd, []string{"secret/app"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "TOKEN=abc123") {
		t.Errorf("expected TOKEN=abc123 in output, got %q", buf.String())
	}
}

func TestCacheInvalidateCmd_Invalidates(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "cache")
	c, _ := syncp.NewSecretCache(dir, time.Minute)
	_ = c.Set("secret/app", map[string]string{"K": "v"})

	var buf bytes.Buffer
	cacheInvalidateCmd.SetOut(&buf)

	if err := cacheInvalidateCmd.RunE(cacheInvalidateCmd, []string{"secret/app"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "invalidated") {
		t.Errorf("expected 'invalidated', got %q", buf.String())
	}
	_, ok := c.Get("secret/app")
	if ok {
		t.Error("expected cache entry to be removed")
	}
}
