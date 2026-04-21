package sync

import (
	"fmt"
	"sort"
	"strings"
)

// SquashStrategy controls how multiple values are combined into one.
type SquashStrategy string

const (
	SquashStrategyConcat SquashStrategy = "concat"
	SquashStrategyFirst  SquashStrategy = "first"
	SquashStrategyLast   SquashStrategy = "last"
)

// SquashStrategyFromString parses a strategy name, returning an error for unknown values.
func SquashStrategyFromString(s string) (SquashStrategy, error) {
	switch strings.ToLower(s) {
	case "concat":
		return SquashStrategyConcat, nil
	case "first":
		return SquashStrategyFirst, nil
	case "last":
		return SquashStrategyLast, nil
	default:
		return "", fmt.Errorf("unknown squash strategy %q: want concat|first|last", s)
	}
}

// Squasher merges keys that share a common prefix into a single key.
type Squasher struct {
	prefix    string
	outKey    string
	sep       string
	strategy  SquashStrategy
}

// NewSquasher creates a Squasher that collapses all keys beginning with prefix
// into outKey using the given strategy. sep is the delimiter used for concat.
func NewSquasher(prefix, outKey, sep string, strategy SquashStrategy) (*Squasher, error) {
	if prefix == "" {
		return nil, fmt.Errorf("squash: prefix must not be empty")
	}
	if outKey == "" {
		return nil, fmt.Errorf("squash: outKey must not be empty")
	}
	return &Squasher{prefix: prefix, outKey: outKey, sep: sep, strategy: strategy}, nil
}

// Apply returns a new map with all keys matching prefix squashed into outKey.
func (s *Squasher) Apply(m map[string]string) (map[string]string, error) {
	var matched []string
	for k := range m {
		if strings.HasPrefix(k, s.prefix) {
			matched = append(matched, k)
		}
	}
	if len(matched) == 0 {
		return m, nil
	}
	sort.Strings(matched)

	out := make(map[string]string, len(m))
	for k, v := range m {
		if !strings.HasPrefix(k, s.prefix) {
			out[k] = v
		}
	}

	switch s.strategy {
	case SquashStrategyFirst:
		out[s.outKey] = m[matched[0]]
	case SquashStrategyLast:
		out[s.outKey] = m[matched[len(matched)-1]]
	case SquashStrategyConcat:
		parts := make([]string, len(matched))
		for i, k := range matched {
			parts[i] = m[k]
		}
		out[s.outKey] = strings.Join(parts, s.sep)
	}
	return out, nil
}
