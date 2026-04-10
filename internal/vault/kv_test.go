package vault

import (
	"testing"
)

func TestParseMount(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{"simple path", "secret/foo/bar", "secret", false},
		{"leading slash", "/secret/foo", "secret", false},
		{"mount only", "secret", "secret", false},
		{"empty path", "", "", true},
		{"slash only", "/", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMount(tt.path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseMount(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ParseMount(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestFullDataPath(t *testing.T) {
	tests := []struct {
		name    string
		mount   string
		sub     string
		version KVVersion
		want    string
	}{
		{"kv1 with sub", "secret", "app/config", KVv1, "secret/app/config"},
		{"kv1 no sub", "secret", "", KVv1, "secret"},
		{"kv2 with sub", "secret", "app/config", KVv2, "secret/data/app/config"},
		{"kv2 no sub", "secret", "", KVv2, "secret/data"},
		{"kv2 leading slash sub", "secret", "/app/config", KVv2, "secret/data/app/config"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FullDataPath(tt.mount, tt.sub, tt.version)
			if got != tt.want {
				t.Errorf("FullDataPath(%q, %q, %d) = %q, want %q", tt.mount, tt.sub, tt.version, got, tt.want)
			}
		})
	}
}

func TestFullMetadataPath(t *testing.T) {
	tests := []struct {
		name    string
		mount   string
		sub     string
		version KVVersion
		want    string
	}{
		{"kv1 with sub", "secret", "app", KVv1, "secret/app"},
		{"kv1 no sub", "secret", "", KVv1, "secret"},
		{"kv2 with sub", "secret", "app", KVv2, "secret/metadata/app"},
		{"kv2 no sub", "secret", "", KVv2, "secret/metadata"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FullMetadataPath(tt.mount, tt.sub, tt.version)
			if got != tt.want {
				t.Errorf("FullMetadataPath(%q, %q, %d) = %q, want %q", tt.mount, tt.sub, tt.version, got, tt.want)
			}
		})
	}
}
