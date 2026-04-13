package sync

import (
	"testing"
)

func mockLookup(env map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		v, ok := env[key]
		return v, ok
	}
}

func TestEnvExpander_Apply_ExpandsKnownVar(t *testing.T) {
	e := NewEnvExpander(withLookup(mockLookup(map[string]string{"REGION": "us-east-1"})))
	out, err := e.Apply(map[string]string{"ENDPOINT": "https://s3.${REGION}.amazonaws.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["ENDPOINT"]; got != "https://s3.us-east-1.amazonaws.com" {
		t.Errorf("got %q, want %q", got, "https://s3.us-east-1.amazonaws.com")
	}
}

func TestEnvExpander_Apply_DollarSyntax(t *testing.T) {
	e := NewEnvExpander(withLookup(mockLookup(map[string]string{"HOME": "/home/user"})))
	out, err := e.Apply(map[string]string{"PATH_VAL": "$HOME/bin"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["PATH_VAL"]; got != "/home/user/bin" {
		t.Errorf("got %q, want %q", got, "/home/user/bin")
	}
}

func TestEnvExpander_Apply_MissingVar_Lenient_KeepsOriginal(t *testing.T) {
	e := NewEnvExpander(withLookup(mockLookup(map[string]string{})))
	out, err := e.Apply(map[string]string{"KEY": "$MISSING"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// lenient mode: value should be preserved as-is or with original reference
	if out["KEY"] == "" {
		t.Error("expected non-empty value in lenient mode")
	}
}

func TestEnvExpander_Apply_MissingVar_Strict_ReturnsError(t *testing.T) {
	e := NewEnvExpander(WithStrictExpand(), withLookup(mockLookup(map[string]string{})))
	_, err := e.Apply(map[string]string{"KEY": "${MISSING}"})
	if err == nil {
		t.Fatal("expected error in strict mode, got nil")
	}
}

func TestEnvExpander_Apply_NoReferences_Unchanged(t *testing.T) {
	e := NewEnvExpander(withLookup(mockLookup(map[string]string{})))
	out, err := e.Apply(map[string]string{"DB_PORT": "5432"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["DB_PORT"]; got != "5432" {
		t.Errorf("got %q, want %q", got, "5432")
	}
}

func TestEnvExpander_Apply_KeyNotModified(t *testing.T) {
	e := NewEnvExpander(withLookup(mockLookup(map[string]string{"X": "1"})))
	out, err := e.Apply(map[string]string{"$MY_KEY": "value"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["$MY_KEY"]; !ok {
		t.Error("key should not be modified by expander")
	}
}

func TestEnvExpander_Apply_MultipleVarsInValue(t *testing.T) {
	e := NewEnvExpander(withLookup(mockLookup(map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	})))
	out, err := e.Apply(map[string]string{"DSN": "postgres://${HOST}:${PORT}/db"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["DSN"]; got != "postgres://localhost:5432/db" {
		t.Errorf("got %q, want %q", got, "postgres://localhost:5432/db")
	}
}
