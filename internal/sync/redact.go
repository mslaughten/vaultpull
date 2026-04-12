package sync

import (
	"regexp"
	"strings"
)

// RedactRule defines a pattern and the replacement string to use when
// a secret value matches the pattern.
type RedactRule struct {
	Pattern     *regexp.Regexp
	Replacement string
}

// Redactor masks sensitive secret values before they are written to
// output (e.g. logs, dry-run output, audit entries).
type Redactor struct {
	rules []RedactRule
}

// NewRedactor builds a Redactor from a slice of raw pattern strings.
// Each pattern is compiled as a regular expression. An error is returned
// if any pattern is invalid.
func NewRedactor(patterns []string) (*Redactor, error) {
	rules := make([]RedactRule, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		rules = append(rules, RedactRule{Pattern: re, Replacement: "***REDACTED***"})
	}
	return &Redactor{rules: rules}, nil
}

// Apply returns a copy of secrets where any value matching a redact
// rule is replaced with the rule's replacement string.
func (r *Redactor) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = r.redactValue(v)
	}
	return out
}

// RedactString replaces all matching substrings within a single string.
func (r *Redactor) RedactString(s string) string {
	return r.redactValue(s)
}

func (r *Redactor) redactValue(v string) string {
	result := v
	for _, rule := range r.rules {
		if rule.Pattern.MatchString(result) {
			result = rule.Pattern.ReplaceAllString(result, rule.Replacement)
		}
	}
	return result
}

// DefaultSensitivePatterns returns a list of common patterns that
// typically indicate sensitive data (passwords, tokens, keys).
func DefaultSensitivePatterns() []string {
	return []string{
		`(?i)password`,
		`(?i)secret`,
		`(?i)token`,
		`(?i)apikey`,
	}
}

// RedactKeys returns a copy of secrets where values whose keys match
// any of the given key substrings (case-insensitive) are replaced.
func RedactKeys(secrets map[string]string, keyPatterns []string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		masked := false
		for _, p := range keyPatterns {
			if strings.Contains(strings.ToLower(k), strings.ToLower(p)) {
				out[k] = "***REDACTED***"
				masked = true
				break
			}
		}
		if !masked {
			out[k] = v
		}
	}
	return out
}
