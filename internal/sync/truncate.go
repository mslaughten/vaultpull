package sync

import (
	"fmt"
	"strings"
)

// TruncateMode controls how values are truncated.
type TruncateMode string

const (
	TruncateModeEnd    TruncateMode = "end"
	TruncateModeStart  TruncateMode = "start"
	TruncateModeMiddle TruncateMode = "middle"
)

// Truncator shortens secret values that exceed a maximum length.
type Truncator struct {
	maxLen  int
	mode    TruncateMode
	ellipsis string
}

// TruncatorOption configures a Truncator.
type TruncatorOption func(*Truncator)

// WithEllipsis sets the ellipsis string appended/prepended on truncation.
func WithEllipsis(e string) TruncatorOption {
	return func(t *Truncator) { t.ellipsis = e }
}

// NewTruncator creates a Truncator with the given max length and mode.
// mode must be one of "end", "start", or "middle".
func NewTruncator(maxLen int, mode string, opts ...TruncatorOption) (*Truncator, error) {
	if maxLen <= 0 {
		return nil, fmt.Errorf("truncate: maxLen must be positive, got %d", maxLen)
	}
	m := TruncateMode(strings.ToLower(mode))
	switch m {
	case TruncateModeEnd, TruncateModeStart, TruncateModeMiddle:
	default:
		return nil, fmt.Errorf("truncate: unknown mode %q (want end|start|middle)", mode)
	}
	t := &Truncator{maxLen: maxLen, mode: m, ellipsis: "..."}
	for _, o := range opts {
		o(t)
	}
	return t, nil
}

// Apply truncates all values in the map that exceed the configured max length.
func (t *Truncator) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = t.truncate(v)
	}
	return out
}

func (t *Truncator) truncate(v string) string {
	if len(v) <= t.maxLen {
		return v
	}
	e := t.ellipsis
	switch t.mode {
	case TruncateModeEnd:
		cutoff := t.maxLen - len(e)
		if cutoff < 0 {
			cutoff = 0
		}
		return v[:cutoff] + e
	case TruncateModeStart:
		cutoff := len(v) - (t.maxLen - len(e))
		if cutoff < 0 {
			cutoff = 0
		}
		return e + v[cutoff:]
	case TruncateModeMiddle:
		half := (t.maxLen - len(e)) / 2
		if half < 0 {
			half = 0
		}
		return v[:half] + e + v[len(v)-half:]
	}
	return v
}
