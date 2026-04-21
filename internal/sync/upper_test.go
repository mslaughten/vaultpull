package sync

import (
	"testing"
)

func TestUpperStrategyFromString_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"keys", "keys"},
		{"values", "values"},
		{"both", "both"},
		{"KEYS", "keys"},
		{"Both", "both"},
	}
	for _, tc := range cases {
		got, err := UpperStrategyFromString(tc.input)
		if err != nil {
			t.Errorf("UpperStrategyFromString(%q) unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("UpperStrategyFromString(%q) = %q; want %q", tc.input, got, tc.want)
		}
	}
}

func TestUpperStrategyFromString_Invalid(t *testing.T) {
	_, err := UpperStrategyFromString("none")
	if err == nil {
		t.Error("expected error for unknown strategy, got nil")
	}
}

func TestNewCaseTransformer_InvalidStrategy_ReturnsError(t *testing.T) {
	_, err := NewCaseTransformer("upper")
	if err == nil {
		t.Error("expected error for invalid strategy")
	}
}

func TestCaseTransformer_Keys_UppercasesKeys(t *testing.T) {
	ct, err := NewCaseTransformer("keys")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := map[string]string{"db_host": "localhost", "app_port": "8080"}
	out, err := ct.Apply(input)
	if err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected key DB_HOST in output")
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected value localhost, got %q", out["DB_HOST"])
	}
}

func TestCaseTransformer_Values_UppercasesValues(t *testing.T) {
	ct, err := NewCaseTransformer("values")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := map[string]string{"env": "production", "mode": "release"}
	out, err := ct.Apply(input)
	if err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	if out["env"] != "PRODUCTION" {
		t.Errorf("expected PRODUCTION, got %q", out["env"])
	}
}

func TestCaseTransformer_Both_UppercasesKeysAndValues(t *testing.T) {
	ct, err := NewCaseTransformer("both")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := map[string]string{"log_level": "debug"}
	out, err := ct.Apply(input)
	if err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	if _, ok := out["LOG_LEVEL"]; !ok {
		t.Error("expected key LOG_LEVEL")
	}
	if out["LOG_LEVEL"] != "DEBUG" {
		t.Errorf("expected DEBUG, got %q", out["LOG_LEVEL"])
	}
}

func TestCaseTransformer_EmptyMap_ReturnsEmpty(t *testing.T) {
	ct, _ := NewCaseTransformer("both")
	out, err := ct.Apply(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
