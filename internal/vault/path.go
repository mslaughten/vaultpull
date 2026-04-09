package vault

import "strings"

// insertAfterMount inserts a segment (e.g. "data" or "metadata") after the
// first path component (the KV mount point).
// e.g. "secret/myapp/db" -> "secret/data/myapp/db"
func insertAfterMount(path, segment string) string {
	parts := strings.SplitN(path, "/", 2)
	if len(parts) == 1 {
		return parts[0] + "/" + segment
	}
	return parts[0] + "/" + segment + "/" + parts[1]
}

// FilterByNamespace returns only those paths that begin with the given
// namespace prefix. If namespace is empty all paths are returned.
func FilterByNamespace(paths []string, namespace string) []string {
	if namespace == "" {
		return paths
	}
	prefix := strings.TrimSuffix(namespace, "/") + "/"
	filtered := make([]string, 0, len(paths))
	for _, p := range paths {
		if strings.HasPrefix(p, prefix) {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

// StripNamespace removes the namespace prefix from a path so the remainder
// can be used as an env-var key segment.
func StripNamespace(path, namespace string) string {
	prefix := strings.TrimSuffix(namespace, "/") + "/"
	return strings.TrimPrefix(path, prefix)
}
