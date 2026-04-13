package sync

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationRule defines a single validation rule for a secret key/value pair.
type ValidationRule struct {
	Key     string
	Pattern *regexp.Regexp
	Message string
}

// Validator checks secret maps against a set of rules.
type Validator struct {
	rules []ValidationRule
}

// ValidationError holds all violations found during validation.
type ValidationError struct {
	Violations []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %s", strings.Join(e.Violations, "; "))
}

// NewValidator builds a Validator from rule strings of the form "KEY=pattern".
// An empty pattern means the key must be present and non-empty.
func NewValidator(rules []string) (*Validator, error) {
	parsed := make([]ValidationRule, 0, len(rules))
	for _, r := range rules {
		idx := strings.Index(r, "=")
		if idx < 0 {
			return nil, fmt.Errorf("invalid rule %q: expected KEY=pattern", r)
		}
		key := strings.TrimSpace(r[:idx])
		if key == "" {
			return nil, fmt.Errorf("invalid rule %q: key must not be empty", r)
		}
		patStr := r[idx+1:]
		var pat *regexp.Regexp
		if patStr != "" {
			var err error
			pat, err = regexp.Compile(patStr)
			if err != nil {
				return nil, fmt.Errorf("invalid rule %q: %w", r, err)
			}
		}
		parsed = append(parsed, ValidationRule{Key: key, Pattern: pat, Message: r})
	}
	return &Validator{rules: parsed}, nil
}

// Validate checks secrets against all rules. Returns *ValidationError if any
// violations are found, nil otherwise.
func (v *Validator) Validate(secrets map[string]string) error {
	var violations []string
	for _, rule := range v.rules {
		val, ok := secrets[rule.Key]
		if !ok || val == "" {
			violations = append(violations, fmt.Sprintf("key %q is required but missing or empty", rule.Key))
			continue
		}
		if rule.Pattern != nil && !rule.Pattern.MatchString(val) {
			violations = append(violations, fmt.Sprintf("key %q value does not match pattern %q", rule.Key, rule.Pattern.String()))
		}
	}
	if len(violations) > 0 {
		return &ValidationError{Violations: violations}
	}
	return nil
}
