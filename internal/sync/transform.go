package sync

import (
	"fmt"
	"strings"
)

// TransformFunc is a function that transforms a secret value.
type TransformFunc func(value string) (string, error)

// Transformer applies a named chain of transformations to secret values.
type Transformer struct {
	rules []transformRule
}

type transformRule struct {
	key  string
	fn   TransformFunc
}

var builtinTransforms = map[string]TransformFunc{
	"upper":   func(v string) (string, error) { return strings.ToUpper(v), nil },
	"lower":   func(v string) (string, error) { return strings.ToLower(v), nil },
	"trimspace": func(v string) (string, error) { return strings.TrimSpace(v), nil },
	"base64": func(v string) (string, error) {
		import64 := strings.NewReplacer("+", "-", "/", "_")
		return import64.Replace(v), nil
	},
}

// NewTransformer creates a Transformer from a slice of "KEY=transform" rule strings.
// Each rule maps a secret key to a named built-in transform.
func NewTransformer(rules []string) (*Transformer, error) {
	t := &Transformer{}
	for _, r := range rules {
		parts := strings.SplitN(r, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("transform: invalid rule %q: expected KEY=transform", r)
		}
		fn, ok := builtinTransforms[parts[1]]
		if !ok {
			return nil, fmt.Errorf("transform: unknown transform %q for key %q", parts[1], parts[0])
		}
		t.rules = append(t.rules, transformRule{key: parts[0], fn: fn})
	}
	return t, nil
}

// Apply runs all matching transform rules against the provided secrets map,
// returning a new map with transformed values.
func (t *Transformer) Apply(secrets map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for _, rule := range t.rules {
		v, ok := out[rule.key]
		if !ok {
			continue
		}
		transformed, err := rule.fn(v)
		if err != nil {
			return nil, fmt.Errorf("transform: key %q: %w", rule.key, err)
		}
		out[rule.key] = transformed
	}
	return out, nil
}
