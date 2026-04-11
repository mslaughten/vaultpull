package sync

import (
	"fmt"
	"regexp"
	"strings"
)

// RenameRule maps a source key pattern to a target key name.
// Pattern may contain a single capture group used in replacement.
type RenameRule struct {
	Pattern     string
	Replacement string
	re          *regexp.Regexp
}

// Compile compiles the rule's pattern. Must be called before Apply.
func (r *RenameRule) Compile() error {
	re, err := regexp.Compile(r.Pattern)
	if err != nil {
		return fmt.Errorf("rename rule: invalid pattern %q: %w", r.Pattern, err)
	}
	r.re = re
	return nil
}

// Apply returns the renamed key if the pattern matches, along with a boolean
// indicating whether the rule matched.
func (r *RenameRule) Apply(key string) (string, bool) {
	if r.re == nil {
		return key, false
	}
	if !r.re.MatchString(key) {
		return key, false
	}
	result := r.re.ReplaceAllString(key, r.Replacement)
	return result, true
}

// Renamer applies a set of RenameRules to a secrets map, producing a new map
// with keys renamed according to the first matching rule.
type Renamer struct {
	rules []RenameRule
}

// NewRenamer constructs a Renamer from the provided rules, compiling each pattern.
func NewRenamer(rules []RenameRule) (*Renamer, error) {
	compiled := make([]RenameRule, len(rules))
	for i, rule := range rules {
		r := rule
		if err := r.Compile(); err != nil {
			return nil, err
		}
		compiled[i] = r
	}
	return &Renamer{rules: compiled}, nil
}

// Apply renames keys in the provided secrets map according to the rules.
// The original map is not modified; a new map is returned.
func (rn *Renamer) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		newKey := k
		for _, rule := range rn.rules {
			if renamed, matched := rule.Apply(k); matched {
				newKey = strings.ToUpper(renamed)
				break
			}
		}
		out[newKey] = v
	}
	return out
}
