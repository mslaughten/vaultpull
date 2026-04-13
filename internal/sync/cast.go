package sync

import (
	"fmt"
	"strings"
)

// CastRule describes a key-to-format mapping used by CastFormatter.
type CastRule struct {
	Key    string
	Format string // "upper", "lower", "title"
}

// CastFormatter applies string-case formatting rules to env map values.
type CastFormatter struct {
	rules []CastRule
}

// NewCastFormatter parses rules of the form "KEY=format" and returns a CastFormatter.
// Valid formats: upper, lower, title.
func NewCastFormatter(raw []string) (*CastFormatter, error) {
	rules := make([]CastRule, 0, len(raw))
	for _, r := range raw {
		parts := strings.SplitN(r, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("cast: invalid rule %q: expected KEY=format", r)
		}
		key := strings.TrimSpace(parts[0])
		format := strings.TrimSpace(strings.ToLower(parts[1]))
		if key == "" {
			return nil, fmt.Errorf("cast: empty key in rule %q", r)
		}
		switch format {
		case "upper", "lower", "title":
		default:
			return nil, fmt.Errorf("cast: unknown format %q in rule %q", format, r)
		}
		rules = append(rules, CastRule{Key: key, Format: format})
	}
	return &CastFormatter{rules: rules}, nil
}

// Apply returns a new map with formatting rules applied to matching keys.
func (c *CastFormatter) Apply(env map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}
	for _, rule := range c.rules {
		v, ok := out[rule.Key]
		if !ok {
			continue
		}
		switch rule.Format {
		case "upper":
			out[rule.Key] = strings.ToUpper(v)
		case "lower":
			out[rule.Key] = strings.ToLower(v)
		case "title":
			out[rule.Key] = strings.Title(v) //nolint:staticcheck
		}
	}
	return out, nil
}
