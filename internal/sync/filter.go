package sync

import (
	"regexp"
	"strings"
)

// KeyFilter filters secret keys based on include/exclude patterns.
type KeyFilter struct {
	include []*regexp.Regexp
	exclude []*regexp.Regexp
}

// NewKeyFilter builds a KeyFilter from lists of glob-style patterns.
// Patterns support '*' as a wildcard.
func NewKeyFilter(include, exclude []string) (*KeyFilter, error) {
	inc, err := compilePatterns(include)
	if err != nil {
		return nil, err
	}
	exc, err := compilePatterns(exclude)
	if err != nil {
		return nil, err
	}
	return &KeyFilter{include: inc, exclude: exc}, nil
}

// Allow returns true when key passes the filter rules.
// A key is allowed when it matches at least one include pattern (or no
// include patterns are defined) and does not match any exclude pattern.
func (f *KeyFilter) Allow(key string) bool {
	if len(f.exclude) > 0 {
		for _, re := range f.exclude {
			if re.MatchString(key) {
				return false
			}
		}
	}
	if len(f.include) == 0 {
		return true
	}
	for _, re := range f.include {
		if re.MatchString(key) {
			return true
		}
	}
	return false
}

// Apply returns a filtered copy of m containing only allowed keys.
func (f *KeyFilter) Apply(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		if f.Allow(k) {
			out[k] = v
		}
	}
	return out
}

func compilePatterns(patterns []string) ([]*regexp.Regexp, error) {
	out := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		regexStr := "^" + strings.ReplaceAll(regexp.QuoteMeta(p), `\*`, `.*`) + "$"
		re, err := regexp.Compile(regexStr)
		if err != nil {
			return nil, err
		}
		out = append(out, re)
	}
	return out, nil
}
