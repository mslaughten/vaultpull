package sync

import (
	"errors"
	"testing"
)

func TestNewTagFilter_Valid(t *testing.T) {
	f, err := NewTagFilter([]string{"env=prod", "team=platform"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(f.tags))
	}
}

func TestNewTagFilter_InvalidPair(t *testing.T) {
	_, err := NewTagFilter([]string{"notavalidtag"})
	if err == nil {
		t.Fatal("expected error for invalid tag pair")
	}
}

func TestNewTagFilter_EmptyKey(t *testing.T) {
	_, err := NewTagFilter([]string{"=value"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestTagFilter_Match_AllPresent(t *testing.T) {
	f, _ := NewTagFilter([]string{"env=prod"})
	if !f.Match(map[string]string{"env": "prod", "team": "x"}) {
		t.Fatal("expected match")
	}
}

func TestTagFilter_Match_Missing(t *testing.T) {
	f, _ := NewTagFilter([]string{"env=prod"})
	if f.Match(map[string]string{"env": "staging"}) {
		t.Fatal("expected no match")
	}
}

func TestTagFilter_Match_Empty(t *testing.T) {
	f, _ := NewTagFilter(nil)
	if !f.Match(map[string]string{}) {
		t.Fatal("empty filter should match everything")
	}
}

func TestFilterPaths_NoTags_ReturnsAll(t *testing.T) {
	f, _ := NewTagFilter(nil)
	paths := []string{"a", "b", "c"}
	got, err := f.FilterPaths(paths, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 paths, got %d", len(got))
	}
}

func TestFilterPaths_FiltersCorrectly(t *testing.T) {
	f, _ := NewTagFilter([]string{"env=prod"})
	paths := []string{"secret/a", "secret/b"}
	metas := map[string]map[string]string{
		"secret/a": {"env": "prod"},
		"secret/b": {"env": "staging"},
	}
	got, err := f.FilterPaths(paths, func(p string) (map[string]string, error) {
		return metas[p], nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0] != "secret/a" {
		t.Fatalf("expected [secret/a], got %v", got)
	}
}

func TestFilterPaths_MetaError_ReturnsError(t *testing.T) {
	f, _ := NewTagFilter([]string{"env=prod"})
	_, err := f.FilterPaths([]string{"secret/x"}, func(p string) (map[string]string, error) {
		return nil, errors.New("vault unavailable")
	})
	if err == nil {
		t.Fatal("expected error from metaFn")
	}
}
