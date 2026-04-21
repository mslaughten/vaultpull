package sync

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// MaskEnvMode controls which part of a value is masked when rendering.
type MaskEnvMode string

const (
	MaskEnvAll    MaskEnvMode = "all"
	MaskEnvPrefix MaskEnvMode = "prefix"
	MaskEnvSuffix MaskEnvMode = "suffix"
)

// MaskEnvOptions configures the MaskEnvRenderer.
type MaskEnvOptions struct {
	Mode        MaskEnvMode
	RevealChars int
	MaskSymbol  string
	Keys        []string // if empty, all keys are masked
}

// MaskEnvRenderer renders a secret map to a writer with values masked.
type MaskEnvRenderer struct {
	opts   MaskEnvOptions
	writer io.Writer
}

// NewMaskEnvRenderer creates a MaskEnvRenderer. If w is nil it defaults to stdout.
func NewMaskEnvRenderer(opts MaskEnvOptions, w io.Writer) (*MaskEnvRenderer, error) {
	if opts.MaskSymbol == "" {
		opts.MaskSymbol = "*"
	}
	if opts.Mode == "" {
		opts.Mode = MaskEnvAll
	}
	if opts.Mode != MaskEnvAll && opts.Mode != MaskEnvPrefix && opts.Mode != MaskEnvSuffix {
		return nil, fmt.Errorf("mask_env: unknown mode %q", opts.Mode)
	}
	if opts.RevealChars < 0 {
		return nil, fmt.Errorf("mask_env: reveal_chars must be >= 0")
	}
	if w == nil {
		w = os.Stdout
	}
	return &MaskEnvRenderer{opts: opts, writer: w}, nil
}

// Render writes masked key=value lines to the writer.
func (r *MaskEnvRenderer) Render(secrets map[string]string) error {
	keySet := make(map[string]struct{}, len(r.opts.Keys))
	for _, k := range r.opts.Keys {
		keySet[k] = struct{}{}
	}

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := secrets[k]
		if len(r.opts.Keys) == 0 {
			v = r.mask(v)
		} else if _, ok := keySet[k]; ok {
			v = r.mask(v)
		}
		if _, err := fmt.Fprintf(r.writer, "%s=%s\n", k, v); err != nil {
			return fmt.Errorf("mask_env: write error: %w", err)
		}
	}
	return nil
}

func (r *MaskEnvRenderer) mask(v string) string {
	sym := r.opts.MaskSymbol
	reveal := r.opts.RevealChars
	switch r.opts.Mode {
	case MaskEnvAll:
		return strings.Repeat(sym, len(v))
	case MaskEnvSuffix:
		if reveal >= len(v) {
			return v
		}
		return v[:reveal] + strings.Repeat(sym, len(v)-reveal)
	case MaskEnvPrefix:
		if reveal >= len(v) {
			return v
		}
		return strings.Repeat(sym, len(v)-reveal) + v[len(v)-reveal:]
	}
	return v
}
