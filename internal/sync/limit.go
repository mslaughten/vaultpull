package sync

import (
	"fmt"
	"sort"
)

// LimitStrategy controls how keys are selected when limiting a map.
type LimitStrategy string

const (
	LimitStrategyFirst LimitStrategy = "first"
	LimitStrategyLast  LimitStrategy = "last"
	LimitStrategyAlpha LimitStrategy = "alpha"
)

// LimitStrategyFromString parses a strategy name, returning an error for unknown values.
func LimitStrategyFromString(s string) (LimitStrategy, error) {
	switch LimitStrategy(s) {
	case LimitStrategyFirst, LimitStrategyLast, LimitStrategyAlpha:
		return LimitStrategy(s), nil
	default:
		return "", fmt.Errorf("unknown limit strategy %q: want first|last|alpha", s)
	}
}

// Limiter reduces a secret map to at most N entries.
type Limiter struct {
	max      int
	strategy LimitStrategy
}

// NewLimiter constructs a Limiter. max must be >= 1.
func NewLimiter(max int, strategy LimitStrategy) (*Limiter, error) {
	if max < 1 {
		return nil, fmt.Errorf("limit: max must be >= 1, got %d", max)
	}
	return &Limiter{max: max, strategy: strategy}, nil
}

// Apply returns a new map containing at most max entries selected by strategy.
func (l *Limiter) Apply(m map[string]string) map[string]string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	switch l.strategy {
	case LimitStrategyAlpha:
		sort.Strings(keys)
	case LimitStrategyLast:
		sort.Strings(keys)
		// reverse
		for i, j := 0, len(keys)-1; i < j; i, j = i+1, j-1 {
			keys[i], keys[j] = keys[j], keys[i]
		}
	default: // first — stable insertion order is not guaranteed in Go maps;
		// sort alpha then take first N for determinism.
		sort.Strings(keys)
	}

	if len(keys) > l.max {
		keys = keys[:l.max]
	}

	out := make(map[string]string, len(keys))
	for _, k := range keys {
		out[k] = m[k]
	}
	return out
}
