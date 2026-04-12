package sync

import (
	"fmt"
	"strings"
)

// MaskMode controls how secret values are masked in output.
type MaskMode int

const (
	// MaskFull replaces the entire value with asterisks.
	MaskFull MaskMode = iota
	// MaskPartial reveals the last N characters of the value.
	MaskPartial
	// MaskNone performs no masking.
	MaskNone
)

// Masker applies a masking strategy to secret values.
type Masker struct {
	mode    MaskMode
	reveal  int    // number of trailing chars to reveal in MaskPartial
	symbol  string // replacement symbol, default "*"
	length  int    // fixed mask length; 0 means match value length
}

// MaskerOption configures a Masker.
type MaskerOption func(*Masker)

// WithRevealChars sets how many trailing characters to reveal (MaskPartial).
func WithRevealChars(n int) MaskerOption {
	return func(m *Masker) { m.reveal = n }
}

// WithMaskSymbol sets the masking symbol (default "*").
func WithMaskSymbol(s string) MaskerOption {
	return func(m *Masker) { m.symbol = s }
}

// WithFixedLength sets a fixed mask length instead of matching value length.
func WithFixedLength(n int) MaskerOption {
	return func(m *Masker) { m.length = n }
}

// NewMasker creates a Masker with the given mode and options.
func NewMasker(mode MaskMode, opts ...MaskerOption) (*Masker, error) {
	if mode < MaskFull || mode > MaskNone {
		return nil, fmt.Errorf("mask: unknown mode %d", mode)
	}
	m := &Masker{mode: mode, reveal: 4, symbol: "*"}
	for _, o := range opts {
		o(m)
	}
	if m.reveal < 0 {
		return nil, fmt.Errorf("mask: reveal chars must be non-negative")
	}
	return m, nil
}

// Apply masks a single value according to the configured mode.
func (m *Masker) Apply(value string) string {
	switch m.mode {
	case MaskNone:
		return value
	case MaskFull:
		return m.fill(len(value))
	case MaskPartial:
		if len(value) <= m.reveal {
			return m.fill(len(value))
		}
		visible := value[len(value)-m.reveal:]
		return m.fill(len(value)-m.reveal) + visible
	}
	return value
}

// ApplyMap masks all values in a map, returning a new map.
func (m *Masker) ApplyMap(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = m.Apply(v)
	}
	return out
}

func (m *Masker) fill(n int) string {
	size := n
	if m.length > 0 {
		size = m.length
	}
	if size <= 0 {
		size = 1
	}
	return strings.Repeat(m.symbol, size)
}
