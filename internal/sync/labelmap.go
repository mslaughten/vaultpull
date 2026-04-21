package sync

import (
	"fmt"
	"strings"
)

// LabelMapper rewrites secret keys by applying a static label→key mapping.
// Each rule has the form "label=key", meaning any secret whose key equals
// label is emitted under the new key name instead.
type LabelMapper struct {
	rules map[string]string
}

// NewLabelMapper constructs a LabelMapper from a slice of "label=newkey" rules.
func NewLabelMapper(rules []string) (*LabelMapper, error) {
	if len(rules) == 0 {
		return &LabelMapper{rules: map[string]string{}}, nil
	}
	m := make(map[string]string, len(rules))
	for _, r := range rules {
		parts := strings.SplitN(r, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("labelmap: invalid rule %q: expected label=newkey", r)
		}
		label := strings.TrimSpace(parts[0])
		newKey := strings.TrimSpace(parts[1])
		if label == "" {
			return nil, fmt.Errorf("labelmap: rule %q has empty label", r)
		}
		if newKey == "" {
			return nil, fmt.Errorf("labelmap: rule %q has empty target key", r)
		}
		m[label] = newKey
	}
	return &LabelMapper{rules: m}, nil
}

// Apply rewrites keys according to the label map.
// Keys not present in the map are passed through unchanged.
// If two source keys map to the same target, the last one wins.
func (lm *LabelMapper) Apply(secrets map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if target, ok := lm.rules[k]; ok {
			out[target] = v
		} else {
			out[k] = v
		}
	}
	return out, nil
}

// Len returns the number of mapping rules.
func (lm *LabelMapper) Len() int { return len(lm.rules) }
