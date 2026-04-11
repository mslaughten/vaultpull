package sync

import (
	"fmt"
	"strings"
)

// TagFilter holds a set of required key=value tags used to filter Vault secret
// paths before syncing. Only paths whose metadata contains ALL specified tags
// are included in the sync plan.
type TagFilter struct {
	tags map[string]string
}

// NewTagFilter parses a slice of "key=value" strings into a TagFilter.
// Returns an error if any entry is not in the expected format.
func NewTagFilter(pairs []string) (*TagFilter, error) {
	tags := make(map[string]string, len(pairs))
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			return nil, fmt.Errorf("invalid tag %q: must be key=value", p)
		}
		tags[parts[0]] = parts[1]
	}
	return &TagFilter{tags: tags}, nil
}

// Match reports whether the provided metadata map satisfies all tags in the
// filter. An empty filter matches everything.
func (f *TagFilter) Match(meta map[string]string) bool {
	for k, v := range f.tags {
		if meta[k] != v {
			return false
		}
	}
	return true
}

// FilterPaths returns only those paths whose metadata satisfies the filter.
// The metadata lookup function is called once per path.
func (f *TagFilter) FilterPaths(paths []string, metaFn func(string) (map[string]string, error)) ([]string, error) {
	if len(f.tags) == 0 {
		return paths, nil
	}
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		meta, err := metaFn(p)
		if err != nil {
			return nil, fmt.Errorf("tag filter: metadata for %q: %w", p, err)
		}
		if f.Match(meta) {
			out = append(out, p)
		}
	}
	return out, nil
}
