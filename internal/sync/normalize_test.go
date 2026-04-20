package sync

import (
	"testing"
)

func TestNormalizeStrategyFromString_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  NormalizeStrategy
	}{
		{"upper", NormalizeUpper},
		{"lower", NormalizeLower},
		{"snake", NormalizeSnake},
		{"kebab", NormalizeKebab},
		{"UPPER", NormalizeUpper},
	}
	for _, tc := range cases {
		got, err := NormalizeStrategyFromString(tc.input)
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("input %q: got %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestNormalizeStrategyFromString_Invalid(t *testing.T) {
	_, err := NormalizeStrategyFromString("title")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestNormalizer_Upper_ConvertsKeys(t *testing.T) {
	n, _ := NewNormalizer("upper")
	out, err := n.Apply(map[string]string{"db_host": "localhost", "api_key": "secret"})
	if err != nil {
		t.Fatal(err)
	}
	if out["DB_HOST"] != "localhost" || out["API_KEY"] != "secret" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestNormalizer_Lower_ConvertsKeys(t *testing.T) {
	n, _ := NewNormalizer("lower")
	out, err := n.Apply(map[string]string{"DB_HOST": "localhost"})
	if err != nil {
		t.Fatal(err)
	}
	if out["db_host"] != "localhost" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestNormalizer_Snake_ReplacesDashes(t *testing.T) {
	n, _ := NewNormalizer("snake")
	out, err := n.Apply(map[string]string{"api-key": "abc"})
	if err != nil {
		t.Fatal(err)
	}
	if out["API_KEY"] != "abc" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestNormalizer_Kebab_ReplacesUnderscores(t *testing.T) {
	n, _ := NewNormalizer("kebab")
	out, err := n.Apply(map[string]string{"DB_HOST": "localhost"})
	if err != nil {
		t.Fatal(err)
	}
	if out["db-host"] != "localhost" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestNormalizer_PreservesValues(t *testing.T) {
	n, _ := NewNormalizer("upper")
	out, err := n.Apply(map[string]string{"key": "MixedCaseValue_123"})
	if err != nil {
		t.Fatal(err)
	}
	if out["KEY"] != "MixedCaseValue_123" {
		t.Errorf("value should be unchanged, got %q", out["KEY"])
	}
}

func TestNewNormalizer_InvalidStrategy_ReturnsError(t *testing.T) {
	_, err := NewNormalizer("pascal")
	if err == nil {
		t.Fatal("expected error")
	}
}
