package sync

import (
	"fmt"
	"strconv"
	"strings"
)

// SliceStrategy controls how a multi-value secret field is sliced into
// individual keys.
type SliceStrategy int

const (
	SliceStrategyIndex  SliceStrategy = iota // KEY_0, KEY_1, ...
	SliceStrategyFirst                       // keep only first element
	SliceStrategyLast                        // keep only last element
	SliceStrategyJoin                        // join all elements with separator
)

// SlicerOptions configures the Slicer.
type SlicerOptions struct {
	Delimiter string
	Separator string // used by SliceStrategyJoin
	Strategy  SliceStrategy
}

// Slicer splits secret values that contain a delimiter into multiple keys.
type Slicer struct {
	opts SlicerOptions
}

// SliceStrategyFromString parses a strategy name.
func SliceStrategyFromString(s string) (SliceStrategy, error) {
	switch strings.ToLower(s) {
	case "index":
		return SliceStrategyIndex, nil
	case "first":
		return SliceStrategyFirst, nil
	case "last":
		return SliceStrategyLast, nil
	case "join":
		return SliceStrategyJoin, nil
	default:
		return 0, fmt.Errorf("unknown slice strategy %q: want index|first|last|join", s)
	}
}

// NewSlicer returns a Slicer with the given options.
// Delimiter defaults to "," if empty.
func NewSlicer(opts SlicerOptions) (*Slicer, error) {
	if opts.Delimiter == "" {
		opts.Delimiter = ","
	}
	if opts.Separator == "" {
		opts.Separator = ","
	}
	return &Slicer{opts: opts}, nil
}

// Apply processes the map and expands values that contain the delimiter.
func (s *Slicer) Apply(m map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(m))
	for k, v := range m {
		parts := strings.Split(v, s.opts.Delimiter)
		if len(parts) <= 1 {
			out[k] = v
			continue
		}
		switch s.opts.Strategy {
		case SliceStrategyFirst:
			out[k] = strings.TrimSpace(parts[0])
		case SliceStrategyLast:
			out[k] = strings.TrimSpace(parts[len(parts)-1])
		case SliceStrategyJoin:
			for i, p := range parts {
				parts[i] = strings.TrimSpace(p)
			}
			out[k] = strings.Join(parts, s.opts.Separator)
		default: // SliceStrategyIndex
			for i, p := range parts {
				out[k+"_"+strconv.Itoa(i)] = strings.TrimSpace(p)
			}
		}
	}
	return out, nil
}
