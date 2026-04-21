package sync

import (
	"fmt"
	"sort"
)

// AlignStrategy controls how keys are aligned between two maps.
type AlignStrategy int

const (
	// AlignIntersection keeps only keys present in both maps.
	AlignIntersection AlignStrategy = iota
	// AlignUnion keeps all keys from both maps.
	AlignUnion
	// AlignLeft keeps all keys from the left (base) map.
	AlignLeft
)

// AlignStrategyFromString parses a strategy name.
func AlignStrategyFromString(s string) (AlignStrategy, error) {
	switch s {
	case "intersection":
		return AlignIntersection, nil
	case "union":
		return AlignUnion, nil
	case "left":
		return AlignLeft, nil
	default:
		return 0, fmt.Errorf("unknown align strategy %q: want intersection|union|left", s)
	}
}

// Aligner aligns two env maps according to a chosen strategy.
type Aligner struct {
	strategy  AlignStrategy
	fillValue string
}

// NewAligner constructs an Aligner with the given strategy.
// fillValue is used for missing keys when strategy is union or left.
func NewAligner(strategy AlignStrategy, fillValue string) *Aligner {
	return &Aligner{strategy: strategy, fillValue: fillValue}
}

// Apply aligns base against ref and returns the resulting map.
func (a *Aligner) Apply(base, ref map[string]string) map[string]string {
	out := make(map[string]string)

	switch a.strategy {
	case AlignIntersection:
		for k, v := range base {
			if _, ok := ref[k]; ok {
				out[k] = v
			}
		}
	case AlignUnion:
		keys := unionKeys(base, ref)
		for _, k := range keys {
			if v, ok := base[k]; ok {
				out[k] = v
			} else {
				out[k] = a.fillValue
			}
		}
	case AlignLeft:
		for k, v := range base {
			out[k] = v
		}
		for k := range ref {
			if _, ok := out[k]; !ok {
				out[k] = a.fillValue
			}
		}
	}
	return out
}

func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
