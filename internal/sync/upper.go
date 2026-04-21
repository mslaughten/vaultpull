package sync

import (
	"fmt"
	"strings"
)

// UpperStrategyFromString parses a case transformation strategy name.
func UpperStrategyFromString(s string) (string, error) {
	switch strings.ToLower(s) {
	case "keys", "values", "both":
		return strings.ToLower(s), nil
	default:
		return "", fmt.Errorf("unknown upper strategy %q: must be keys, values, or both", s)
	}
}

// CaseTransformer applies upper-case transformation to keys, values, or both.
type CaseTransformer struct {
	strategy string
}

// NewCaseTransformer creates a CaseTransformer for the given strategy.
func NewCaseTransformer(strategy string) (*CaseTransformer, error) {
	s, err := UpperStrategyFromString(strategy)
	if err != nil {
		return nil, err
	}
	return &CaseTransformer{strategy: s}, nil
}

// Apply returns a new map with keys and/or values converted to upper case.
func (c *CaseTransformer) Apply(m map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(m))
	for k, v := range m {
		newKey := k
		newVal := v
		if c.strategy == "keys" || c.strategy == "both" {
			newKey = strings.ToUpper(k)
		}
		if c.strategy == "values" || c.strategy == "both" {
			newVal = strings.ToUpper(v)
		}
		out[newKey] = newVal
	}
	return out, nil
}
