package sync

import (
	"fmt"
	"regexp"
	"strings"
)

// LintRule describes a single naming convention check applied to secret keys.
type LintRule struct {
	Name    string
	Pattern *regexp.Regexp
	Message string
}

// LintViolation holds a key that failed a lint rule along with the rule name and message.
type LintViolation struct {
	Key     string
	Rule    string
	Message string
}

func (v LintViolation) Error() string {
	return fmt.Sprintf("key %q violates rule %q: %s", v.Key, v.Rule, v.Message)
}

// Linter checks secret keys against a set of naming-convention rules.
type Linter struct {
	rules []LintRule
}

// DefaultLintRules returns a sensible set of built-in rules.
var DefaultLintRules = []LintRule{
	{
		Name:    "uppercase",
		Pattern: regexp.MustCompile(`^[A-Z0-9_]+$`),
		Message: "key must be uppercase with only letters, digits, and underscores",
	},
	{
		Name:    "no-leading-digit",
		Pattern: regexp.MustCompile(`^[^0-9]`),
		Message: "key must not start with a digit",
	},
	{
		Name:    "no-double-underscore",
		Pattern: regexp.MustCompile(`^(?!.*__)`),
		Message: "key must not contain consecutive underscores",
	},
}

// NewLinter constructs a Linter. Pass nil to use DefaultLintRules.
func NewLinter(rules []LintRule) *Linter {
	if rules == nil {
		rules = DefaultLintRules
	}
	return &Linter{rules: rules}
}

// Check evaluates all keys in secrets against every rule and returns all violations.
func (l *Linter) Check(secrets map[string]string) []LintViolation {
	var violations []LintViolation
	for key := range secrets {
		for _, rule := range l.rules {
			if !rule.Pattern.MatchString(key) {
				violations = append(violations, LintViolation{
					Key:     key,
					Rule:    rule.Name,
					Message: rule.Message,
				})
			}
		}
	}
	return violations
}

// Summary returns a human-readable summary of all violations.
func (l *Linter) Summary(violations []LintViolation) string {
	if len(violations) == 0 {
		return "lint: all keys passed"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "lint: %d violation(s)\n", len(violations))
	for _, v := range violations {
		fmt.Fprintf(&sb, "  - %s\n", v.Error())
	}
	return strings.TrimRight(sb.String(), "\n")
}
