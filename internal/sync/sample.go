package sync

import (
	"fmt"
	"math/rand"
	"sort"
)

// SampleStrategy controls how keys are selected during sampling.
type SampleStrategy string

const (
	SampleStrategyRandom SampleStrategy = "random"
	SampleStrategyFirst  SampleStrategy = "first"
	SampleStrategyLast   SampleStrategy = "last"
)

// SampleStrategyFromString parses a strategy name.
func SampleStrategyFromString(s string) (SampleStrategy, error) {
	switch SampleStrategy(s) {
	case SampleStrategyRandom, SampleStrategyFirst, SampleStrategyLast:
		return SampleStrategy(s), nil
	default:
		return "", fmt.Errorf("unknown sample strategy %q: want random|first|last", s)
	}
}

// Sampler selects a subset of keys from a secret map.
type Sampler struct {
	n        int
	strategy SampleStrategy
	rng      *rand.Rand
}

// NewSampler creates a Sampler that keeps at most n keys using the given strategy.
func NewSampler(n int, strategy SampleStrategy, rng *rand.Rand) (*Sampler, error) {
	if n <= 0 {
		return nil, fmt.Errorf("sample: n must be > 0, got %d", n)
	}
	if rng == nil {
		rng = rand.New(rand.NewSource(42))
	}
	return &Sampler{n: n, strategy: strategy, rng: rng}, nil
}

// Apply returns a new map containing at most n keys selected by the strategy.
func (s *Sampler) Apply(secrets map[string]string) map[string]string {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	switch s.strategy {
	case SampleStrategyRandom:
		s.rng.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })
	case SampleStrategyLast:
		for i, j := 0, len(keys)-1; i < j; i, j = i+1, j-1 {
			keys[i], keys[j] = keys[j], keys[i]
		}
	}

	if s.n < len(keys) {
		keys = keys[:s.n]
	}

	out := make(map[string]string, len(keys))
	for _, k := range keys {
		out[k] = secrets[k]
	}
	return out
}
