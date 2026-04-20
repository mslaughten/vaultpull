package sync

import (
	"fmt"
	"sort"
)

// MergeEnvStrategyFromString parses a merge-env strategy name.
func MergeEnvStrategyFromString(s string) (MergeStrategy, error) {
	switch s {
	case "overwrite":
		return MergeOverwrite, nil
	case "keep":
		return MergeKeepExisting, nil
	case "vault":
		return MergeVaultWins, nil
	default:
		return 0, fmt.Errorf("unknown merge-env strategy %q: want overwrite|keep|vault", s)
	}
}

// EnvMerger merges two env maps using a named strategy and optional key filter.
type EnvMerger struct {
	strategy MergeStrategy
	filter   []string // if non-empty, only these keys are merged from src
}

// NewEnvMerger constructs an EnvMerger.
func NewEnvMerger(strategy MergeStrategy, keys []string) *EnvMerger {
	return &EnvMerger{strategy: strategy, filter: keys}
}

// Apply merges src into dst according to the configured strategy.
// It returns a new map; dst and src are not mutated.
func (m *EnvMerger) Apply(dst, src map[string]string) map[string]string {
	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	allowed := make(map[string]bool, len(m.filter))
	for _, k := range m.filter {
		allowed[k] = true
	}

	for _, k := range sortedKeys(src) {
		if len(allowed) > 0 && !allowed[k] {
			continue
		}
		v := src[k]
		switch m.strategy {
		case MergeOverwrite:
			out[k] = v
		case MergeKeepExisting:
			if _, exists := out[k]; !exists {
				out[k] = v
			}
		case MergeVaultWins:
			if v != "" {
				out[k] = v
			} else if _, exists := out[k]; !exists {
				out[k] = v
			}
		}
	}
	return out
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
