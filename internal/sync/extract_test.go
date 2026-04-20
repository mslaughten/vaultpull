package sync

import (
	"testing"
)

func TestExtractStrategyFromString_Valid(t *testing.T) {
	for _, s := range []string{"prefix", "suffix", "regex", "PREFIX", "Suffix"} {
		_, err := ExtractStrategyFromString(s)
		if err != nil {
			t.Errorf("expected no error for %q, got %v", s, err)
		}
	}
}

func TestExtractStrategyFromString_Invalid(t *testing.T) {
	_, err := ExtractStrategyFromString("unknown")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestNewExtractor_EmptyPattern_ReturnsError(t *testing.T) {
	_, err := NewExtractor("prefix", "")
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNewExtractor_InvalidRegex_ReturnsError(t *testing.T) {
	_, err := NewExtractor("regex", "[invalid")
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestExtractor_Prefix_KeepsMatchingKeys(t *testing.T) {
	e, err := NewExtractor("prefix", "DB_")
	if err != nil {
		t.Fatal(err)
	}
	in := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_NAME": "vaultpull",
	}
	out := e.Apply(in)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["APP_NAME"]; ok {
		t.Error("APP_NAME should not be present")
	}
}

func TestExtractor_Suffix_KeepsMatchingKeys(t *testing.T) {
	e, err := NewExtractor("suffix", "_KEY")
	if err != nil {
		t.Fatal(err)
	}
	in := map[string]string{
		"API_KEY":    "abc",
		"SECRET_KEY": "xyz",
		"HOST":       "localhost",
	}
	out := e.Apply(in)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["HOST"]; ok {
		t.Error("HOST should not be present")
	}
}

func TestExtractor_Regex_KeepsMatchingKeys(t *testing.T) {
	e, err := NewExtractor("regex", `^[A-Z]+_\d+$`)
	if err != nil {
		t.Fatal(err)
	}
	in := map[string]string{
		"VAR_1":  "a",
		"VAR_2":  "b",
		"other":  "c",
		"NO_NUM": "d",
	}
	out := e.Apply(in)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d: %v", len(out), out)
	}
}

func TestExtractor_Apply_EmptyMap_ReturnsEmpty(t *testing.T) {
	e, _ := NewExtractor("prefix", "X_")
	out := e.Apply(map[string]string{})
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %v", out)
	}
}
