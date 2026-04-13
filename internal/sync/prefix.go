package sync

import (
	"fmt"
	"strings"
)

// PrefixStrategy controls how key prefixes are applied or stripped.
type PrefixStrategy int

const (
	// PrefixAdd prepends a string to every key.
	PrefixAdd PrefixStrategy = iota
	// PrefixStrip removes a leading string from every key.
	PrefixStrip
)

// Prefixer adds or strips a prefix from secret keys.
type Prefixer struct {
	prefix   string
	strategy PrefixStrategy
}

// NewPrefixer creates a Prefixer with the given prefix and strategy string.
// strategy must be "add" or "strip".
func NewPrefixer(prefix, strategy string) (*Prefixer, error) {
	if prefix == "" {
		return nil, fmt.Errorf("prefix: prefix must not be empty")
	}
	var s PrefixStrategy
	switch strings.ToLower(strategy) {
	case "add":
		s = PrefixAdd
	case "strip":
		s = PrefixStrip
	default:
		return nil, fmt.Errorf("prefix: unknown strategy %q (want add|strip)", strategy)
	}
	return &Prefixer{prefix: prefix, strategy: s}, nil
}

// Apply transforms the provided map by adding or stripping the prefix from
// every key. Keys that do not carry the prefix are passed through unchanged
// when stripping.
func (p *Prefixer) Apply(secrets map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		switch p.strategy {
		case PrefixAdd:
			out[p.prefix+k] = v
		case PrefixStrip:
			out[strings.TrimPrefix(k, p.prefix)] = v
		}
	}
	return out, nil
}

// Name returns a human-readable stage name used by the pipeline.
func (p *Prefixer) Name() string {
	switch p.strategy {
	case PrefixAdd:
		return fmt.Sprintf("prefix-add(%s)", p.prefix)
	default:
		return fmt.Sprintf("prefix-strip(%s)", p.prefix)
	}
}
