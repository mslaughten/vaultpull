package sync

import (
	"testing"
)

func TestRenameRule_Apply_Matches(t *testing.T) {
	rule := RenameRule{Pattern: `^DB_(.+)$`, Replacement: `DATABASE_$1`}
	if err := rule.Compile(); err != nil {
		t.Fatalf("unexpected compile error: %v", err)
	}
	got, matched := rule.Apply("DB_PASSWORD")
	if !matched {
		t.Fatal("expected match")
	}
	if got != "DATABASE_PASSWORD" {
		t.Errorf("got %q, want DATABASE_PASSWORD", got)
	}
}

func TestRenameRule_Apply_NoMatch(t *testing.T) {
	rule := RenameRule{Pattern: `^DB_(.+)$`, Replacement: `DATABASE_$1`}
	_ = rule.Compile()
	got, matched := rule.Apply("APP_SECRET")
	if matched {
		t.Fatal("expected no match")
	}
	if got != "APP_SECRET" {
		t.Errorf("key should be unchanged, got %q", got)
	}
}

func TestRenameRule_Compile_InvalidPattern(t *testing.T) {
	rule := RenameRule{Pattern: `[invalid`, Replacement: `X`}
	if err := rule.Compile(); err == nil {
		t.Fatal("expected compile error for invalid pattern")
	}
}

func TestNewRenamer_InvalidRule_ReturnsError(t *testing.T) {
	_, err := NewRenamer([]RenameRule{
		{Pattern: `[bad`, Replacement: `X`},
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRenamer_Apply_RenamesMatchingKeys(t *testing.T) {
	rn, err := NewRenamer([]RenameRule{
		{Pattern: `^DB_(.+)$`, Replacement: `DATABASE_$1`},
		{Pattern: `^CACHE_(.+)$`, Replacement: `REDIS_$1`},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	input := map[string]string{
		"DB_HOST":    "localhost",
		"CACHE_PORT": "6379",
		"APP_NAME":   "myapp",
	}
	out := rn.Apply(input)

	if v, ok := out["DATABASE_HOST"]; !ok || v != "localhost" {
		t.Errorf("expected DATABASE_HOST=localhost, got %v", out)
	}
	if v, ok := out["REDIS_PORT"]; !ok || v != "6379" {
		t.Errorf("expected REDIS_PORT=6379, got %v", out)
	}
	if v, ok := out["APP_NAME"]; !ok || v != "myapp" {
		t.Errorf("expected APP_NAME=myapp unchanged, got %v", out)
	}
}

func TestRenamer_Apply_FirstRuleWins(t *testing.T) {
	rn, err := NewRenamer([]RenameRule{
		{Pattern: `^SECRET$`, Replacement: `FIRST`},
		{Pattern: `^SECRET$`, Replacement: `SECOND`},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := rn.Apply(map[string]string{"SECRET": "val"})
	if _, ok := out["FIRST"]; !ok {
		t.Errorf("expected first rule to win, got keys: %v", out)
	}
}

func TestRenamer_Apply_EmptyRules_Passthrough(t *testing.T) {
	rn, err := NewRenamer(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := map[string]string{"FOO": "bar"}
	out := rn.Apply(input)
	if out["FOO"] != "bar" {
		t.Errorf("expected passthrough, got %v", out)
	}
}
