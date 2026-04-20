package sync

import (
	"fmt"
	"strings"
)

// PivotStrategy controls how keys are pivoted into a new structure.
type PivotStrategy int

const (
	PivotStrategyKeyToValue PivotStrategy = iota
	PivotStrategyValueToKey
)

// PivotStrategyFromString parses a strategy name.
func PivotStrategyFromString(s string) (PivotStrategy, error) {
	switch strings.ToLower(s) {
	case "key_to_value", "ktv":
		return PivotStrategyKeyToValue, nil
	case "value_to_key", "vtk":
		return PivotStrategyValueToKey, nil
	default:
		return 0, fmt.Errorf("unknown pivot strategy %q: want key_to_value or value_to_key", s)
	}
}

// Pivoter swaps or restructures keys and values in a secret map.
type Pivotter struct {
	strategy  PivotStrategy
	prefix    string
	overwrite bool
}

// NewPivoter creates a Pivotter with the given strategy.
// prefix is prepended to generated keys when strategy is ValueToKey.
func NewPivoter(strategy PivotStrategy, prefix string, overwrite bool) *Pivotter {
	return &Pivotter{strategy: strategy, prefix: prefix, overwrite: overwrite}
}

// Apply transforms src according to the pivot strategy and merges the result
// into dst. When overwrite is false, existing keys in dst are preserved.
func (p *Pivotter) Apply(dst, src map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	switch p.strategy {
	case PivotStrategyKeyToValue:
		// Each key in src becomes a value; the original value becomes the key.
		for k, v := range src {
			newKey := p.prefix + v
			if _, exists := out[newKey]; exists && !p.overwrite {
				continue
			}
			out[newKey] = k
		}
	case PivotStrategyValueToKey:
		// Each value in src becomes a key; the original key becomes the value.
		for k, v := range src {
			newKey := p.prefix + k
			if _, exists := out[newKey]; exists && !p.overwrite {
				continue
			}
			out[newKey] = v
		}
	default:
		return nil, fmt.Errorf("pivot: unsupported strategy %d", p.strategy)
	}

	return out, nil
}
