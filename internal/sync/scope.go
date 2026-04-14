package sync

import (
	"fmt"
	"strings"
)

// ScopeStrategy controls how keys are scoped when mapping secrets.
type ScopeStrategy string

const (
	ScopeStrategyPrefix ScopeStrategy = "prefix"
	ScopeStrategyExact  ScopeStrategy = "exact"
	ScopeStrategyGlob   ScopeStrategy = "glob"
)

// ScopeRule defines a single scoping rule: a pattern and the scope name.
type ScopeRule struct {
	Pattern  string
	Scope    string
	strategy ScopeStrategy
}

// Scoper maps secret keys to named scopes based on configurable rules.
type Scoper struct {
	rules []ScopeRule
}

// NewScoper creates a Scoper from a slice of "pattern=scope" rule strings.
// Each rule must contain exactly one "=" separator.
// Strategy is inferred: "*" in pattern implies glob, trailing ":" implies prefix, else exact.
func NewScoper(rules []string) (*Scoper, error) {
	parsed := make([]ScopeRule, 0, len(rules))
	for _, r := range rules {
		parts := strings.SplitN(r, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("scope: invalid rule %q: missing '=' separator", r)
		}
		pattern, scope := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		if pattern == "" {
			return nil, fmt.Errorf("scope: invalid rule %q: empty pattern", r)
		}
		if scope == "" {
			return nil, fmt.Errorf("scope: invalid rule %q: empty scope", r)
		}
		var strategy ScopeStrategy
		switch {
		case strings.Contains(pattern, "*"):
			strategy = ScopeStrategyGlob
		case strings.HasSuffix(pattern, ":"):
			strategy = ScopeStrategyPrefix
			pattern = strings.TrimSuffix(pattern, ":")
		default:
			strategy = ScopeStrategyExact
		}
		parsed = append(parsed, ScopeRule{Pattern: pattern, Scope: scope, strategy: strategy})
	}
	return &Scoper{rules: parsed}, nil
}

// Apply returns a map of scope name -> keys that belong to that scope.
// Keys that match no rule are placed under the "default" scope.
func (s *Scoper) Apply(secrets map[string]string) map[string][]string {
	result := make(map[string][]string)
	for key := range secrets {
		scope := s.matchScope(key)
		result[scope] = append(result[scope], key)
	}
	return result
}

func (s *Scoper) matchScope(key string) string {
	for _, rule := range s.rules {
		switch rule.strategy {
		case ScopeStrategyExact:
			if key == rule.Pattern {
				return rule.Scope
			}
		case ScopeStrategyPrefix:
			if strings.HasPrefix(key, rule.Pattern) {
				return rule.Scope
			}
		case ScopeStrategyGlob:
			if globMatch(rule.Pattern, key) {
				return rule.Scope
			}
		}
	}
	return "default"
}

// globMatch performs simple single-wildcard glob matching.
func globMatch(pattern, s string) bool {
	parts := strings.SplitN(pattern, "*", 2)
	if len(parts) == 1 {
		return pattern == s
	}
	return strings.HasPrefix(s, parts[0]) && strings.HasSuffix(s, parts[1])
}
