package sync

import (
	"fmt"
	"os"
	"strings"
)

// EnvExpander replaces ${ENV_VAR} or $ENV_VAR references in secret values
// with values from the host environment. It is useful when a .env file
// contains placeholders that should be resolved at sync time.
type EnvExpander struct {
	strict bool
	lookup func(string) (string, bool)
}

// EnvExpanderOption configures an EnvExpander.
type EnvExpanderOption func(*EnvExpander)

// WithStrictExpand causes Apply to return an error when a referenced
// environment variable is not set on the host.
func WithStrictExpand() EnvExpanderOption {
	return func(e *EnvExpander) { e.strict = true }
}

// withLookup replaces the default os.LookupEnv for testing.
func withLookup(fn func(string) (string, bool)) EnvExpanderOption {
	return func(e *EnvExpander) { e.lookup = fn }
}

// NewEnvExpander creates an EnvExpander with the provided options.
func NewEnvExpander(opts ...EnvExpanderOption) *EnvExpander {
	e := &EnvExpander{lookup: os.LookupEnv}
	for _, o := range opts {
		o(e)
	}
	return e
}

// Apply iterates over secrets and expands any environment variable references
// found in the values. The keys are never modified.
func (e *EnvExpander) Apply(secrets map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		expanded, err := e.expand(v)
		if err != nil {
			return nil, fmt.Errorf("env_expand: key %q: %w", k, err)
		}
		out[k] = expanded
	}
	return out, nil
}

// expand resolves all ${VAR} and $VAR references within a single string.
func (e *EnvExpander) expand(s string) (string, error) {
	var expandErr error
	result := os.Expand(s, func(key string) string {
		if expandErr != nil {
			return ""
		}
		val, ok := e.lookup(key)
		if !ok {
			if e.strict {
				expandErr = fmt.Errorf("environment variable %q is not set", key)
				return ""
			}
			// lenient: keep the original reference
			if strings.HasPrefix(s, "${") {
				return "${" + key + "}"
			}
			return "$" + key
		}
		return val
	})
	if expandErr != nil {
		return "", expandErr
	}
	return result, nil
}
