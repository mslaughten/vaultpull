package sync

import (
	"fmt"
	"sort"
	"strings"
)

// SortStrategy defines how keys in a secret map should be ordered.
type SortStrategy int

const (
	SortStrategyAlpha SortStrategy = iota // alphabetical ascending
	SortStrategyAlphaDesc                 // alphabetical descending
	SortStrategyLength                    // shortest key first
	SortStrategyLengthDesc                // longest key first
)

// Sorter reorders the keys of a secret map according to a strategy.
type Sorter struct {
	strategy SortStrategy
}

// SortStrategyFromString parses a strategy name into a SortStrategy.
func SortStrategyFromString(s string) (SortStrategy, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "alpha", "asc", "":
		return SortStrategyAlpha, nil
	case "alpha-desc", "desc":
		return SortStrategyAlphaDesc, nil
	case "length", "len":
		return SortStrategyLength, nil
	case "length-desc", "len-desc":
		return SortStrategyLengthDesc, nil
	default:
		return 0, fmt.Errorf("unknown sort strategy %q: want alpha, alpha-desc, length, length-desc", s)
	}
}

// NewSorter creates a Sorter for the given strategy.
func NewSorter(strategy SortStrategy) *Sorter {
	return &Sorter{strategy: strategy}
}

// Apply returns a new map with the same key/value pairs; it also returns the
// keys in sorted order so callers can iterate deterministically.
func (s *Sorter) Apply(secrets map[string]string) (map[string]string, []string) {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}

	switch s.strategy {
	case SortStrategyAlphaDesc:
		sort.Slice(keys, func(i, j int) bool { return keys[i] > keys[j] })
	case SortStrategyLength:
		sort.Slice(keys, func(i, j int) bool {
			if len(keys[i]) == len(keys[j]) {
				return keys[i] < keys[j]
			}
			return len(keys[i]) < len(keys[j])
		})
	case SortStrategyLengthDesc:
		sort.Slice(keys, func(i, j int) bool {
			if len(keys[i]) == len(keys[j]) {
				return keys[i] > keys[j]
			}
			return len(keys[i]) > len(keys[j])
		})
	default: // SortStrategyAlpha
		sort.Strings(keys)
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	return out, keys
}
