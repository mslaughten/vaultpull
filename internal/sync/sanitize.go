package sync

import (
	"fmt"
	"regexp"
	"strings"
)

// SanitizeRule describes a single key sanitization rule.
type SanitizeRule struct {
	Pattern     *regexp.Regexp
	Replacement string
}

// Sanitizer replaces disallowed characters in secret keys.
type Sanitizer struct {
	rules []SanitizeRule
}

// NewSanitizer builds a Sanitizer from a slice of "pattern=replacement" strings.
// If rules is empty, a default rule that replaces any non-alphanumeric/underscore
// character with "_" is applied.
func NewSanitizer(rules []string) (*Sanitizer, error) {
	if len(rules) == 0 {
		defaultRe := regexp.MustCompile(`[^A-Za-z0-9_]`)
		return &Sanitizer{rules: []SanitizeRule{{Pattern: defaultRe, Replacement: "_"}}}, nil
	}

	parsed := make([]SanitizeRule, 0, len(rules))
	for _, r := range rules {
		idx := strings.Index(r, "=")
		if idx < 0 {
			return nil, fmt.Errorf("sanitize: rule %q missing '=' separator", r)
		}
		patStr := r[:idx]
		replacement := r[idx+1:]
		re, err := regexp.Compile(patStr)
		if err != nil {
			return nil, fmt.Errorf("sanitize: invalid pattern %q: %w", patStr, err)
		}
		parsed = append(parsed, SanitizeRule{Pattern: re, Replacement: replacement})
	}
	return &Sanitizer{rules: parsed}, nil
}

// Apply rewrites every key in secrets according to the sanitizer's rules.
// Values are left untouched. If two keys collide after sanitization the last
// one (in iteration order) wins.
func (s *Sanitizer) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		clean := k
		for _, rule := range s.rules {
			clean = rule.Pattern.ReplaceAllString(clean, rule.Replacement)
		}
		out[clean] = v
	}
	return out
}
