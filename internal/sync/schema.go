package sync

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
)

// FieldRule defines a validation rule for a single secret key.
type FieldRule struct {
	Key      string `json:"key"`
	Required bool   `json:"required"`
	Pattern  string `json:"pattern,omitempty"`
	regexp   *regexp.Regexp
}

// SchemaValidator validates a secret map against a set of field rules.
type SchemaValidator struct {
	rules []FieldRule
}

// NewSchemaValidator parses rules from JSON (e.g. [{"key":"DB_URL","required":true,"pattern":"^postgres://"}]).
func NewSchemaValidator(src string) (*SchemaValidator, error) {
	var rules []FieldRule
	if err := json.Unmarshal([]byte(src), &rules); err != nil {
		return nil, fmt.Errorf("schema: invalid JSON: %w", err)
	}
	for i, r := range rules {
		if r.Key == "" {
			return nil, fmt.Errorf("schema: rule %d missing key", i)
		}
		if r.Pattern != "" {
			re, err := regexp.Compile(r.Pattern)
			if err != nil {
				return nil, fmt.Errorf("schema: rule %q invalid pattern: %w", r.Key, err)
			}
			rules[i].regexp = re
		}
	}
	return &SchemaValidator{rules: rules}, nil
}

// Validate checks secrets against all rules and returns a list of violations.
func (v *SchemaValidator) Validate(secrets map[string]string) []string {
	var violations []string
	for _, r := range v.rules {
		val, ok := secrets[r.Key]
		if r.Required && !ok {
			violations = append(violations, fmt.Sprintf("required key %q is missing", r.Key))
			continue
		}
		if ok && r.regexp != nil && !r.regexp.MatchString(val) {
			violations = append(violations, fmt.Sprintf("key %q value does not match pattern %q", r.Key, r.Pattern))
		}
	}
	return violations
}

// WriteReport writes validation violations to w, returning whether validation passed.
func (v *SchemaValidator) WriteReport(w io.Writer, violations []string) bool {
	if len(violations) == 0 {
		fmt.Fprintln(w, "schema validation passed")
		return true
	}
	fmt.Fprintf(w, "schema validation failed (%d violation(s)):\n", len(violations))
	for _, msg := range violations {
		fmt.Fprintf(w, "  - %s\n", msg)
	}
	return false
}

// DefaultSchemaPath is the conventional schema file name.
const DefaultSchemaPath = ".vaultschema.json"

// LoadSchemaFile reads a schema file from disk and returns a validator.
func LoadSchemaFile(path string) (*SchemaValidator, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("schema: cannot read file %q: %w", path, err)
	}
	return NewSchemaValidator(string(data))
}
