package sync

import (
	"fmt"
	"strings"
)

// NamespaceMapper rewrites secret keys by applying a namespace-to-prefix
// mapping. Keys whose path begins with a known namespace are rewritten so
// the namespace segment is replaced by a configured prefix string.
type NamespaceMapper struct {
	rules map[string]string // namespace -> prefix
}

// NewNamespaceMapper creates a NamespaceMapper from a slice of "namespace=prefix"
// rule strings. An error is returned if any rule is malformed or has an empty
// namespace segment.
func NewNamespaceMapper(rules []string) (*NamespaceMapper, error) {
	m := &NamespaceMapper{rules: make(map[string]string, len(rules))}
	for _, r := range rules {
		parts := strings.SplitN(r, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("namespace mapper: invalid rule %q: expected namespace=prefix", r)
		}
		ns := strings.TrimSpace(parts[0])
		if ns == "" {
			return nil, fmt.Errorf("namespace mapper: rule %q has empty namespace", r)
		}
		m.rules[ns] = strings.TrimSpace(parts[1])
	}
	return m, nil
}

// Apply rewrites the keys in src according to the configured namespace rules.
// For each key, if any namespace matches the beginning of the key (using
// "namespace/" as a separator), the namespace segment is replaced by the
// mapped prefix. Keys that do not match any rule are left unchanged.
func (m *NamespaceMapper) Apply(src map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[m.rewrite(k)] = v
	}
	return out, nil
}

// rewrite returns the rewritten key for a single input key.
func (m *NamespaceMapper) rewrite(key string) string {
	for ns, prefix := range m.rules {
		segment := ns + "/"
		if strings.HasPrefix(key, segment) {
			remainder := strings.TrimPrefix(key, segment)
			if prefix == "" {
				return remainder
			}
			return prefix + "_" + remainder
		}
		// exact match with no trailing slash
		if key == ns {
			if prefix == "" {
				return key
			}
			return prefix
		}
	}
	return key
}
