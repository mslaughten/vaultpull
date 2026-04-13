package sync

import (
	"fmt"
	"regexp"
	"strings"
)

// interpolatePattern matches ${VAR_NAME} or $VAR_NAME style references.
var interpolatePattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Z_][A-Z0-9_]*)`)

// Interpolator replaces variable references inside secret values with
// values drawn from a lookup map (typically the same secret map or a
// set of ambient environment variables).
type Interpolator struct {
	lookup map[string]string
	strict bool
}

// InterpolatorOption configures an Interpolator.
type InterpolatorOption func(*Interpolator)

// WithStrictInterpolation causes Apply to return an error when a
// referenced variable is not found in the lookup map.
func WithStrictInterpolation() InterpolatorOption {
	return func(i *Interpolator) { i.strict = true }
}

// NewInterpolator creates an Interpolator that resolves variable
// references using the provided lookup map.
func NewInterpolator(lookup map[string]string, opts ...InterpolatorOption) *Interpolator {
	ip := &Interpolator{lookup: lookup}
	for _, o := range opts {
		o(ip)
	}
	return ip
}

// Apply performs variable interpolation on every value in secrets.
// Keys are left unchanged. A copy of the map is returned.
func (ip *Interpolator) Apply(secrets map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		resolved, err := ip.resolve(v)
		if err != nil {
			return nil, fmt.Errorf("interpolate key %q: %w", k, err)
		}
		out[k] = resolved
	}
	return out, nil
}

func (ip *Interpolator) resolve(s string) (string, error) {
	var resolveErr error
	result := interpolatePattern.ReplaceAllStringFunc(s, func(match string) string {
		if resolveErr != nil {
			return match
		}
		name := strings.TrimPrefix(strings.Trim(match, "${}"), "$")
		if val, ok := ip.lookup[name]; ok {
			return val
		}
		if ip.strict {
			resolveErr = fmt.Errorf("undefined variable %q", name)
			return match
		}
		return match
	})
	if resolveErr != nil {
		return "", resolveErr
	}
	return result, nil
}
