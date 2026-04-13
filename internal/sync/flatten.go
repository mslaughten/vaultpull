package sync

import (
	"fmt"
	"strings"
)

// FlattenStrategy controls how nested map keys are joined.
type FlattenStrategy int

const (
	// FlattenUnderscore joins keys with underscores: parent_child_key
	FlattenUnderscore FlattenStrategy = iota
	// FlattenDot joins keys with dots: parent.child.key
	FlattenDot
)

// FlattenStrategyFromString parses a strategy name.
func FlattenStrategyFromString(s string) (FlattenStrategy, error) {
	switch strings.ToLower(s) {
	case "underscore", "":
		return FlattenUnderscore, nil
	case "dot":
		return FlattenDot, nil
	default:
		return 0, fmt.Errorf("unknown flatten strategy %q: want underscore or dot", s)
	}
}

// Flattener collapses nested map[string]any structures into a flat
// map[string]string, joining key segments with the chosen separator.
type Flattener struct {
	strategy FlattenStrategy
}

// NewFlattener constructs a Flattener with the given strategy.
func NewFlattener(strategy FlattenStrategy) *Flattener {
	return &Flattener{strategy: strategy}
}

// separator returns the key separator for the configured strategy.
func (f *Flattener) separator() string {
	if f.strategy == FlattenDot {
		return "."
	}
	return "_"
}

// Flatten converts a nested map into a flat map[string]string.
// Only leaf values are kept; intermediate maps are expanded.
func (f *Flattener) Flatten(input map[string]any) map[string]string {
	out := make(map[string]string)
	f.flatten(input, "", out)
	return out
}

func (f *Flattener) flatten(m map[string]any, prefix string, out map[string]string) {
	sep := f.separator()
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + sep + k
		}
		switch val := v.(type) {
		case map[string]any:
			f.flatten(val, key, out)
		case string:
			out[key] = val
		default:
			out[key] = fmt.Sprintf("%v", val)
		}
	}
}
