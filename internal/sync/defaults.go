package sync

import "fmt"

// DefaultsStrategy controls how missing keys are filled.
type DefaultsStrategy int

const (
	// DefaultsSkip leaves missing keys absent (no-op).
	DefaultsSkip DefaultsStrategy = iota
	// DefaultsApply fills in missing keys from the defaults map.
	DefaultsApply
)

// DefaultsStrategyFromString parses a strategy name.
func DefaultsStrategyFromString(s string) (DefaultsStrategy, error) {
	switch s {
	case "", "skip":
		return DefaultsSkip, nil
	case "apply":
		return DefaultsApply, nil
	default:
		return DefaultsSkip, fmt.Errorf("unknown defaults strategy %q: want skip|apply", s)
	}
}

// Defaulter fills missing keys in a secret map from a set of default values.
type Defaulter struct {
	defaults map[string]string
	strategy DefaultsStrategy
}

// NewDefaulter constructs a Defaulter. defaults must be non-nil.
func NewDefaulter(defaults map[string]string, strategy DefaultsStrategy) (*Defaulter, error) {
	if defaults == nil {
		return nil, fmt.Errorf("defaults map must not be nil")
	}
	return &Defaulter{defaults: defaults, strategy: strategy}, nil
}

// Apply returns a new map that contains all keys from secrets, with any
// missing keys filled from the defaults map when strategy == DefaultsApply.
func (d *Defaulter) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	if d.strategy == DefaultsApply {
		for k, v := range d.defaults {
			if _, exists := out[k]; !exists {
				out[k] = v
			}
		}
	}
	return out
}

// Keys returns the list of default key names.
func (d *Defaulter) Keys() []string {
	keys := make([]string, 0, len(d.defaults))
	for k := range d.defaults {
		keys = append(keys, k)
	}
	return keys
}
