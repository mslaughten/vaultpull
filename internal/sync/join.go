package sync

import (
	"fmt"
	"strings"
)

// JoinStrategy controls how values from multiple maps are combined.
type JoinStrategy string

const (
	JoinStrategyConcat    JoinStrategy = "concat"
	JoinStrategyFirstOnly JoinStrategy = "first"
	JoinStrategyLastOnly  JoinStrategy = "last"
)

// JoinStrategyFromString parses a strategy name, returning an error for unknown values.
func JoinStrategyFromString(s string) (JoinStrategy, error) {
	switch JoinStrategy(s) {
	case JoinStrategyConcat, JoinStrategyFirstOnly, JoinStrategyLastOnly:
		return JoinStrategy(s), nil
	default:
		return "", fmt.Errorf("unknown join strategy %q: want concat|first|last", s)
	}
}

// Joiner merges two secret maps according to a chosen strategy.
type Joiner struct {
	strategy  JoinStrategy
	separator string
}

// NewJoiner creates a Joiner with the given strategy and separator (used by concat).
func NewJoiner(strategy JoinStrategy, separator string) (*Joiner, error) {
	if _, err := JoinStrategyFromString(string(strategy)); err != nil {
		return nil, err
	}
	if separator == "" {
		separator = ","
	}
	return &Joiner{strategy: strategy, separator: separator}, nil
}

// Apply merges src into dst according to the configured strategy.
// For concat, values present in both maps are joined with the separator.
// For first, existing dst values are preserved; missing keys are filled from src.
// For last, src values always overwrite dst values.
func (j *Joiner) Apply(dst, src map[string]string) map[string]string {
	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}
	for k, sv := range src {
		dv, exists := out[k]
		switch j.strategy {
		case JoinStrategyConcat:
			if exists {
				out[k] = strings.Join([]string{dv, sv}, j.separator)
			} else {
				out[k] = sv
			}
		case JoinStrategyFirstOnly:
			if !exists {
				out[k] = sv
			}
		case JoinStrategyLastOnly:
			out[k] = sv
		}
	}
	return out
}
