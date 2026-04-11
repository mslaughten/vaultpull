package sync

import (
	"testing"
)

func TestNewKeyFilter_Valid(t *testing.T) {
	f, err := NewKeyFilter([]string{"DB_*"}, []string{"DB_PASS*"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestKeyFilter_Allow_NoPatterns_AllowsAll(t *testing.T) {
	f, _ := NewKeyFilter(nil, nil)
	for _, key := range []string{"FOO", "BAR", "BAZ_123"} {
		if !f.Allow(key) {
			t.Errorf("expected %q to be allowed", key)
		}
	}
}

func TestKeyFilter_Allow_IncludePattern(t *testing.T) {
	f, _ := NewKeyFilter([]string{"DB_*"}, nil)
	if !f.Allow("DB_HOST") {
		t.Error("expected DB_HOST to be allowed")
	}
	if f.Allow("APP_KEY") {
		t.Error("expected APP_KEY to be blocked")
	}
}

func TestKeyFilter_Allow_ExcludePattern(t *testing.T) {
	f, _ := NewKeyFilter(nil, []string{"*_SECRET", "*_PASS"})
	if f.Allow("DB_SECRET") {
		t.Error("expected DB_SECRET to be excluded")
	}
	if f.Allow("APP_PASS") {
		t.Error("expected APP_PASS to be excluded")
	}
	if !f.Allow("DB_HOST") {
		t.Error("expected DB_HOST to be allowed")
	}
}

func TestKeyFilter_Allow_IncludeAndExclude(t *testing.T) {
	f, _ := NewKeyFilter([]string{"DB_*"}, []string{"DB_PASS*"})
	if !f.Allow("DB_HOST") {
		t.Error("expected DB_HOST to be allowed")
	}
	if f.Allow("DB_PASSWORD") {
		t.Error("expected DB_PASSWORD to be excluded")
	}
	if f.Allow("APP_KEY") {
		t.Error("expected APP_KEY to be blocked by include filter")
	}
}

func TestKeyFilter_Apply_ReturnsFilteredMap(t *testing.T) {
	f, _ := NewKeyFilter([]string{"DB_*"}, []string{"DB_PASS*"})
	input := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "secret",
		"APP_KEY":     "value",
	}
	out := f.Apply(input)
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", out["DB_HOST"])
	}
}

func TestKeyFilter_Apply_EmptyInput(t *testing.T) {
	f, _ := NewKeyFilter([]string{"DB_*"}, nil)
	out := f.Apply(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d entries", len(out))
	}
}
