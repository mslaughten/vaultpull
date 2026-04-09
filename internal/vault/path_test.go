package vault

import (
	"testing"
)

func TestInsertAfterMount(t *testing.T) {
	cases := []struct {
		path    string
		seg     string
		want    string
	}{
		{"secret/myapp/db", "data", "secret/data/myapp/db"},
		{"secret/myapp", "data", "secret/data/myapp"},
		{"secret", "data", "secret/data"},
		{"kv/prod/api", "metadata", "kv/metadata/prod/api"},
	}
	for _, tc := range cases {
		got := insertAfterMount(tc.path, tc.seg)
		if got != tc.want {
			t.Errorf("insertAfterMount(%q, %q) = %q; want %q", tc.path, tc.seg, got, tc.want)
		}
	}
}

func TestFilterByNamespace(t *testing.T) {
	paths := []string{
		"secret/prod/db",
		"secret/prod/api",
		"secret/staging/db",
		"secret/other",
	}

	got := FilterByNamespace(paths, "secret/prod")
	if len(got) != 2 {
		t.Fatalf("expected 2 paths, got %d: %v", len(got), got)
	}

	empty := FilterByNamespace(paths, "")
	if len(empty) != len(paths) {
		t.Fatalf("empty namespace should return all paths")
	}
}

func TestStripNamespace(t *testing.T) {
	cases := []struct {
		path      string
		namespace string
		want      string
	}{
		{"secret/prod/db", "secret/prod", "db"},
		{"secret/prod/api/key", "secret/prod", "api/key"},
		{"secret/other", "secret/prod", "secret/other"},
	}
	for _, tc := range cases {
		got := StripNamespace(tc.path, tc.namespace)
		if got != tc.want {
			t.Errorf("StripNamespace(%q, %q) = %q; want %q", tc.path, tc.namespace, got, tc.want)
		}
	}
}
