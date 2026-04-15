package sync

import (
	"fmt"
	"strings"
)

// CollapseStrategy determines how duplicate-prefix keys are collapsed.
type CollapseStrategy int

const (
	CollapseStrategyFirst CollapseStrategy = iota
	CollapseStrategyLast
	CollapseStrategyConcat
)

// CollapseStrategyFromString parses a strategy name.
func CollapseStrategyFromString(s string) (CollapseStrategy, error) {
	switch strings.ToLower(s) {
	case "first":
		return CollapseStrategyFirst, nil
	case "last":
		return CollapseStrategyLast, nil
	case "concat":
		return CollapseStrategyConcat, nil
	default:
		return 0, fmt.Errorf("unknown collapse strategy %q: want first|last|concat", s)
	}
}

// Collapser merges keys that share a common prefix into a single key.
type Collapser struct {
	prefix    string
	outKey    string
	sep       string
	strategy  CollapseStrategy
}

// NewCollapser creates a Collapser that collapses all keys starting with
// prefix into outKey using the given strategy. sep is used for concat.
func NewCollapser(prefix, outKey, sep string, strategy CollapseStrategy) (*Collapser, error) {
	if prefix == "" {
		return nil, fmt.Errorf("collapse: prefix must not be empty")
	}
	if outKey == "" {
		return nil, fmt.Errorf("collapse: outKey must not be empty")
	}
	return &Collapser{prefix: prefix, outKey: outKey, sep: sep, strategy: strategy}, nil
}

// Apply collapses matching keys in m and returns a new map.
func (c *Collapser) Apply(m map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(m))
	var parts []string

	for k, v := range m {
		if !strings.HasPrefix(k, c.prefix) {
			out[k] = v
			continue
		}
		parts = append(parts, v)
	}

	if len(parts) == 0 {
		return out, nil
	}

	switch c.strategy {
	case CollapseStrategyFirst:
		out[c.outKey] = parts[0]
	case CollapseStrategyLast:
		out[c.outKey] = parts[len(parts)-1]
	case CollapseStrategyConcat:
		out[c.outKey] = strings.Join(parts, c.sep)
	}

	return out, nil
}
