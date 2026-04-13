package sync

import (
	"testing"
)

func TestFlattenStrategyFromString_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  FlattenStrategy
	}{
		{"underscore", FlattenUnderscore},
		{"", FlattenUnderscore},
		{"dot", FlattenDot},
		{"DOT", FlattenDot},
	}
	for _, tc := range cases {
		got, err := FlattenStrategyFromString(tc.input)
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("input %q: got %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestFlattenStrategyFromString_Invalid(t *testing.T) {
	_, err := FlattenStrategyFromString("csv")
	if err == nil {
		t.Fatal("expected error for unknown strategy, got nil")
	}
}

func TestFlattener_Flat_Underscore(t *testing.T) {
	f := NewFlattener(FlattenUnderscore)
	input := map[string]any{
		"db": map[string]any{
			"host": "localhost",
			"port": "5432",
		},
		"app": "myapp",
	}
	out := f.Flatten(input)
	expect := map[string]string{
		"db_host": "localhost",
		"db_port": "5432",
		"app":     "myapp",
	}
	for k, v := range expect {
		if out[k] != v {
			t.Errorf("key %q: got %q, want %q", k, out[k], v)
		}
	}
}

func TestFlattener_Flat_Dot(t *testing.T) {
	f := NewFlattener(FlattenDot)
	input := map[string]any{
		"cache": map[string]any{
			"ttl": "300",
		},
	}
	out := f.Flatten(input)
	if out["cache.ttl"] != "300" {
		t.Errorf("got %q, want %q", out["cache.ttl"], "300")
	}
}

func TestFlattener_NonStringLeaf(t *testing.T) {
	f := NewFlattener(FlattenUnderscore)
	input := map[string]any{
		"timeout": 30,
	}
	out := f.Flatten(input)
	if out["timeout"] != "30" {
		t.Errorf("got %q, want \"30\"", out["timeout"])
	}
}

func TestFlattener_EmptyMap(t *testing.T) {
	f := NewFlattener(FlattenUnderscore)
	out := f.Flatten(map[string]any{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestFlattener_DeeplyNested(t *testing.T) {
	f := NewFlattener(FlattenUnderscore)
	input := map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": "deep",
			},
		},
	}
	out := f.Flatten(input)
	if out["a_b_c"] != "deep" {
		t.Errorf("got %q, want \"deep\"", out["a_b_c"])
	}
}
