package sync

import (
	"fmt"
	"strings"
)

// SplitStrategy controls how a map is split into multiple maps.
type SplitStrategy string

const (
	SplitByPrefix    SplitStrategy = "prefix"
	SplitByDelimiter SplitStrategy = "delimiter"
)

// SplitStrategyFromString parses a SplitStrategy from a string.
func SplitStrategyFromString(s string) (SplitStrategy, error) {
	switch strings.ToLower(s) {
	case "prefix":
		return SplitByPrefix, nil
	case "delimiter":
		return SplitByDelimiter, nil
	default:
		return "", fmt.Errorf("unknown split strategy %q: want prefix or delimiter", s)
	}
}

// Splitter partitions a flat key/value map into named buckets.
type Splitter struct {
	strategy  SplitStrategy
	delimiter string
}

// NewSplitter creates a Splitter with the given strategy and delimiter.
// delimiter is used as the boundary character (e.g. "_" or ".").
func NewSplitter(strategy SplitStrategy, delimiter string) (*Splitter, error) {
	if delimiter == "" {
		delimiter = "_"
	}
	return &Splitter{strategy: strategy, delimiter: delimiter}, nil
}

// Apply partitions src into named groups.
// Each key in the returned map is a bucket name; each value is the
// subset of src that belongs to that bucket.
func (s *Splitter) Apply(src map[string]string) map[string]map[string]string {
	out := make(map[string]map[string]string)
	for k, v := range src {
		bucket := s.bucketFor(k)
		if out[bucket] == nil {
			out[bucket] = make(map[string]string)
		}
		out[bucket][k] = v
	}
	return out
}

func (s *Splitter) bucketFor(key string) string {
	switch s.strategy {
	case SplitByPrefix:
		if idx := strings.Index(key, s.delimiter); idx > 0 {
			return key[:idx]
		}
		return "default"
	case SplitByDelimiter:
		parts := strings.SplitN(key, s.delimiter, 2)
		if len(parts) == 2 && parts[0] != "" {
			return parts[0]
		}
		return "default"
	default:
		return "default"
	}
}
