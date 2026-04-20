package sync

import (
	"fmt"
	"strings"
)

// ReorderStrategy controls how keys are reordered within a map.
type ReorderStrategy string

const (
	ReorderStrategyExplicit ReorderStrategy = "explicit" // keys listed first, rest follow
	ReorderStrategyReverse  ReorderStrategy = "reverse"  // reverse current order
)

// ReorderStrategyFromString parses a strategy name.
func ReorderStrategyFromString(s string) (ReorderStrategy, error) {
	switch strings.ToLower(s) {
	case "explicit":
		return ReorderStrategyExplicit, nil
	case "reverse":
		return ReorderStrategyReverse, nil
	default:
		return "", fmt.Errorf("unknown reorder strategy %q: want explicit|reverse", s)
	}
}

// Reorderer reorders keys in a map according to a strategy.
type Reorderer struct {
	strategy ReorderStrategy
	keys     []string // used by explicit strategy
}

// NewReorderer constructs a Reorderer.
// For the explicit strategy, keys is the ordered list of keys to place first.
func NewReorderer(strategy ReorderStrategy, keys []string) (*Reorderer, error) {
	if strategy == ReorderStrategyExplicit && len(keys) == 0 {
		return nil, fmt.Errorf("reorder: explicit strategy requires at least one key")
	}
	return &Reorderer{strategy: strategy, keys: keys}, nil
}

// Apply returns a new map preserving all entries but with an ordered key slice.
// The returned map is identical in content; callers that need ordering should
// use OrderedKeys alongside the map.
func (r *Reorderer) Apply(m map[string]string) (map[string]string, []string, error) {
	switch r.strategy {
	case ReorderStrategyExplicit:
		return r.applyExplicit(m)
	case ReorderStrategyReverse:
		return r.applyReverse(m)
	default:
		return nil, nil, fmt.Errorf("reorder: unknown strategy %q", r.strategy)
	}
}

func (r *Reorderer) applyExplicit(m map[string]string) (map[string]string, []string, error) {
	seen := make(map[string]bool, len(r.keys))
	ordered := make([]string, 0, len(m))
	for _, k := range r.keys {
		if _, ok := m[k]; ok {
			ordered = append(ordered, k)
			seen[k] = true
		}
	}
	for k := range m {
		if !seen[k] {
			ordered = append(ordered, k)
		}
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out, ordered, nil
}

func (r *Reorderer) applyReverse(m map[string]string) (map[string]string, []string, error) {
	sorted := sortedKeys(m)
	for i, j := 0, len(sorted)-1; i < j; i, j = i+1, j-1 {
		sorted[i], sorted[j] = sorted[j], sorted[i]
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out, sorted, nil
}

// sortedKeys returns alphabetically sorted keys of m.
func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sortStrings(keys)
	return keys
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
