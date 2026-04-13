package sync

import (
	"testing"
)

func TestNewPrefixer_EmptyPrefix_ReturnsError(t *testing.T) {
	_, err := NewPrefixer("", "add")
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestNewPrefixer_UnknownStrategy_ReturnsError(t *testing.T) {
	_, err := NewPrefixer("APP_", "replace")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestNewPrefixer_Valid(t *testing.T) {
	for _, s := range []string{"add", "strip", "ADD", "STRIP"} {
		_, err := NewPrefixer("APP_", s)
		if err != nil {
			t.Fatalf("strategy %q: unexpected error: %v", s, err)
		}
	}
}

func TestPrefixer_Add_PrependsPrefix(t *testing.T) {
	p, _ := NewPrefixer("APP_", "add")
	in := map[string]string{"HOST": "localhost", "PORT": "5432"}
	out, err := p.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, k := range []string{"APP_HOST", "APP_PORT"} {
		if _, ok := out[k]; !ok {
			t.Errorf("expected key %q in output", k)
		}
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}

func TestPrefixer_Strip_RemovesPrefix(t *testing.T) {
	p, _ := NewPrefixer("APP_", "strip")
	in := map[string]string{"APP_HOST": "localhost", "APP_PORT": "5432"}
	out, err := p.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, k := range []string{"HOST", "PORT"} {
		if _, ok := out[k]; !ok {
			t.Errorf("expected key %q in output", k)
		}
	}
}

func TestPrefixer_Strip_PassesThroughUnmatched(t *testing.T) {
	p, _ := NewPrefixer("APP_", "strip")
	in := map[string]string{"OTHER_KEY": "value"}
	out, err := p.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["OTHER_KEY"]; !ok {
		t.Error("expected unmatched key to pass through unchanged")
	}
}

func TestPrefixer_Name(t *testing.T) {
	add, _ := NewPrefixer("APP_", "add")
	if add.Name() != "prefix-add(APP_)" {
		t.Errorf("unexpected name: %s", add.Name())
	}
	strip, _ := NewPrefixer("APP_", "strip")
	if strip.Name() != "prefix-strip(APP_)" {
		t.Errorf("unexpected name: %s", strip.Name())
	}
}
