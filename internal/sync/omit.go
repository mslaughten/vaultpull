package sync

import (
	"fmt"
	"regexp"
)

// Omitter removes keys from a secret map whose values match one or more
// regular-expression patterns. It is the logical complement of KeyFilter's
// include mode: rather than keeping matches it discards them.
type Omitter struct {
	patterns []*regexp.Regexp
}

// NewOmitter compiles each pattern and returns an Omitter. An empty pattern
// slice is valid and produces a no-op Omitter.
func NewOmitter(patterns []string) (*Omitter, error) {
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		if p == "" {
			return nil, fmt.Errorf("omit: pattern must not be empty")
		}
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("omit: invalid pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}
	return &Omitter{patterns: compiled}, nil
}

// Apply returns a copy of m with any key whose value matches at least one
// pattern removed. Keys whose values do not match any pattern are kept.
func (o *Omitter) Apply(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		if o.matches(v) {
			continue
		}
		out[k] = v
	}
	return out
}

// matches reports whether value is matched by any compiled pattern.
func (o *Omitter) matches(value string) bool {
	for _, re := range o.patterns {
		if re.MatchString(value) {
			return true
		}
	}
	return false
}
