package sync

import (
	"testing"
)

func TestSplitStrategyFromString_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  SplitStrategy
	}{
		{"prefix", SplitByPrefix},
		{"delimiter", SplitByDelimiter},
		{"PREFIX", SplitByPrefix},
	}
	for _, tc := range cases {
		got, err := SplitStrategyFromString(tc.input)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("got %q, want %q", got, tc.want)
		}
	}
}

func TestSplitStrategyFromString_Invalid(t *testing.T) {
	_, err := SplitStrategyFromString("unknown")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSplitter_Prefix_SplitsOnFirstDelimiter(t *testing.T) {
	s, _ := NewSplitter(SplitByPrefix, "_")
	src := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_HOST":  "db.local",
		"PLAIN":    "value",
	}
	out := s.Apply(src)
	if len(out["APP"]) != 2 {
		t.Errorf("APP bucket: got %d keys, want 2", len(out["APP"]))
	}
	if len(out["DB"]) != 1 {
		t.Errorf("DB bucket: got %d keys, want 1", len(out["DB"]))
	}
	if out["default"]["PLAIN"] != "value" {
		t.Errorf("expected PLAIN in default bucket")
	}
}

func TestSplitter_Delimiter_SplitsOnFirstOccurrence(t *testing.T) {
	s, _ := NewSplitter(SplitByDelimiter, ".")
	src := map[string]string{
		"prod.DB_URL":  "postgres://",
		"prod.API_KEY": "secret",
		"dev.DB_URL":   "sqlite://",
		"standalone":   "yes",
	}
	out := s.Apply(src)
	if len(out["prod"]) != 2 {
		t.Errorf("prod bucket: got %d, want 2", len(out["prod"]))
	}
	if len(out["dev"]) != 1 {
		t.Errorf("dev bucket: got %d, want 1", len(out["dev"]))
	}
	if out["default"]["standalone"] != "yes" {
		t.Errorf("expected standalone in default bucket")
	}
}

func TestSplitter_EmptyMap_ReturnsEmpty(t *testing.T) {
	s, _ := NewSplitter(SplitByPrefix, "_")
	out := s.Apply(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty result, got %d buckets", len(out))
	}
}

func TestSplitter_DefaultDelimiter_UsesUnderscore(t *testing.T) {
	s, err := NewSplitter(SplitByPrefix, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	src := map[string]string{"FOO_BAR": "baz"}
	out := s.Apply(src)
	if out["FOO"]["FOO_BAR"] != "baz" {
		t.Errorf("expected FOO_BAR in FOO bucket")
	}
}
