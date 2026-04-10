package vault

import (
	"fmt"
	"strings"
)

// KVVersion represents the KV secrets engine version.
type KVVersion int

const (
	// KVv1 is the original KV secrets engine.
	KVv1 KVVersion = 1
	// KVv2 is the versioned KV secrets engine.
	KVv2 KVVersion = 2
)

// MountInfo holds metadata about a KV mount.
type MountInfo struct {
	Mount   string
	Version KVVersion
}

// ParseMount extracts the mount point from a secret path.
// For example, "secret/foo/bar" returns "secret".
func ParseMount(path string) (string, error) {
	path = strings.TrimPrefix(path, "/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) == 0 || parts[0] == "" {
		return "", fmt.Errorf("invalid path %q: cannot determine mount", path)
	}
	return parts[0], nil
}

// FullDataPath returns the full API path for reading a secret,
// inserting "/data/" for KVv2 mounts.
func FullDataPath(mount, subPath string, version KVVersion) string {
	subPath = strings.TrimPrefix(subPath, "/")
	if version == KVv2 {
		if subPath == "" {
			return mount + "/data"
		}
		return mount + "/data/" + subPath
	}
	if subPath == "" {
		return mount
	}
	return mount + "/" + subPath
}

// FullMetadataPath returns the full API path for listing secrets,
// inserting "/metadata/" for KVv2 mounts.
func FullMetadataPath(mount, subPath string, version KVVersion) string {
	subPath = strings.TrimPrefix(subPath, "/")
	if version == KVv2 {
		if subPath == "" {
		"
		}
		return mount + "/metadata/" + subPath
	}
	if subPath == "" {
		return mount
	}
	return mount + "/" + subPath
}
