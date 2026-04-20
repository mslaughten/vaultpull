package sync

import (
	"fmt"
	"strings"
)

// NormalizeStrategy controls how key normalization is applied.
type NormalizeStrategy string

const (
	NormalizeUpper  NormalizeStrategy = "upper"
	NormalizeLower  NormalizeStrategy = "lower"
	NormalizeSnake  NormalizeStrategy = "snake"
	NormalizeKebab  NormalizeStrategy = "kebab"
)

// NormalizeStrategyFromString parses a strategy name, returning an error for
// unknown values.
func NormalizeStrategyFromString(s string) (NormalizeStrategy, error) {
	switch NormalizeStrategy(strings.ToLower(s)) {
	case NormalizeUpper, NormalizeLower, NormalizeSnake, NormalizeKebab:
		return NormalizeStrategy(strings.ToLower(s)), nil
	}
	return "", fmt.Errorf("unknown normalize strategy %q: want upper|lower|snake|kebab", s)
}

// Normalizer rewrites map keys according to a chosen strategy.
type Normalizer struct {
	strategy NormalizeStrategy
}

// NewNormalizer constructs a Normalizer for the given strategy string.
func NewNormalizer(strategy string) (*Normalizer, error) {
	s, err := NormalizeStrategyFromString(strategy)
	if err != nil {
		return nil, err
	}
	return &Normalizer{strategy: s}, nil
}

// Apply returns a new map whose keys have been rewritten according to the
// configured strategy. Values are preserved unchanged.
func (n *Normalizer) Apply(m map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[n.rewrite(k)] = v
	}
	return out, nil
}

func (n *Normalizer) rewrite(key string) string {
	switch n.strategy {
	case NormalizeUpper:
		return strings.ToUpper(key)
	case NormalizeLower:
		return strings.ToLower(key)
	case NormalizeSnake:
		return strings.ReplaceAll(strings.ToUpper(key), "-", "_")
	case NormalizeKebab:
		return strings.ReplaceAll(strings.ToLower(key), "_", "-")
	}
	return key
}
