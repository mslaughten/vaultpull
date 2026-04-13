package sync

import (
	"fmt"
	"strconv"
	"strings"
)

// TypeCastRule defines how a specific key's value should be cast.
type TypeCastRule struct {
	Key  string
	Type string // "string", "int", "float", "bool"
}

// TypeCaster applies type coercion to env map values, normalising them
// to their canonical string representation (e.g. "true" for booleans,
// "42" for integers). This is useful before writing to .env files or
// pushing to Vault to ensure consistent formatting.
type TypeCaster struct {
	rules []TypeCastRule
}

// NewTypeCaster builds a TypeCaster from a slice of "KEY=type" rule strings.
// Supported types: string, int, float, bool.
func NewTypeCaster(rules []string) (*TypeCaster, error) {
	parsed := make([]TypeCastRule, 0, len(rules))
	for _, r := range rules {
		parts := strings.SplitN(r, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("typecast: invalid rule %q: expected KEY=type", r)
		}
		key := strings.TrimSpace(parts[0])
		typ := strings.TrimSpace(parts[1])
		if key == "" {
			return nil, fmt.Errorf("typecast: rule %q has empty key", r)
		}
		switch typ {
		case "string", "int", "float", "bool":
		default:
			return nil, fmt.Errorf("typecast: unknown type %q in rule %q", typ, r)
		}
		parsed = append(parsed, TypeCastRule{Key: key, Type: typ})
	}
	return &TypeCaster{rules: parsed}, nil
}

// Apply coerces values in env according to the registered rules.
// Returns an error if a value cannot be parsed as the target type.
func (tc *TypeCaster) Apply(env map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}
	for _, rule := range tc.rules {
		v, ok := out[rule.Key]
		if !ok {
			continue
		}
		coerced, err := coerce(v, rule.Type)
		if err != nil {
			return nil, fmt.Errorf("typecast: key %q: %w", rule.Key, err)
		}
		out[rule.Key] = coerced
	}
	return out, nil
}

func coerce(value, typ string) (string, error) {
	switch typ {
	case "string":
		return value, nil
	case "int":
		_, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return "", fmt.Errorf("cannot cast %q to int", value)
		}
		return value, nil
	case "float":
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return "", fmt.Errorf("cannot cast %q to float", value)
		}
		return strconv.FormatFloat(f, 'f', -1, 64), nil
	case "bool":
		b, err := strconv.ParseBool(value)
		if err != nil {
			return "", fmt.Errorf("cannot cast %q to bool", value)
		}
		return strconv.FormatBool(b), nil
	}
	return value, nil
}
