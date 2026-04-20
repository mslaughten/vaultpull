package sync

import (
	"fmt"
	"regexp"
	"strings"
)

// ExtractStrategyFromString parses a strategy name into an extract mode.
func ExtractStrategyFromString(s string) (string, error) {
	switch strings.ToLower(s) {
	case "prefix", "suffix", "regex":
		return strings.ToLower(s), nil
	default:
		return "", fmt.Errorf("unknown extract strategy %q: must be prefix, suffix, or regex", s)
	}
}

// Extractor selects a subset of keys from a map based on a matching strategy.
type Extractor struct {
	strategy string
	pattern  string
	re       *regexp.Regexp
}

// NewExtractor creates an Extractor with the given strategy and pattern.
// For "prefix" and "suffix" strategies the pattern is a plain string.
// For "regex" the pattern is compiled as a regular expression.
func NewExtractor(strategy, pattern string) (*Extractor, error) {
	strat, err := ExtractStrategyFromString(strategy)
	if err != nil {
		return nil, err
	}
	if pattern == "" {
		return nil, fmt.Errorf("extract pattern must not be empty")
	}
	e := &Extractor{strategy: strat, pattern: pattern}
	if strat == "regex" {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid regex pattern %q: %w", pattern, err)
		}
		e.re = re
	}
	return e, nil
}

// Apply returns a new map containing only the entries whose keys match the
// configured strategy and pattern.
func (e *Extractor) Apply(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		if e.matches(k) {
			out[k] = v
		}
	}
	return out
}

func (e *Extractor) matches(key string) bool {
	switch e.strategy {
	case "prefix":
		return strings.HasPrefix(key, e.pattern)
	case "suffix":
		return strings.HasSuffix(key, e.pattern)
	case "regex":
		return e.re.MatchString(key)
	}
	return false
}
