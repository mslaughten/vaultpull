package sync

import (
	"fmt"
	"sort"
	"strings"
)

// IndexStrategy controls how keys are indexed.
type IndexStrategy int

const (
	IndexStrategyAlpha IndexStrategy = iota
	IndexStrategyInsertion
)

// IndexStrategyFromString parses a strategy name.
func IndexStrategyFromString(s string) (IndexStrategy, error) {
	switch strings.ToLower(s) {
	case "alpha":
		return IndexStrategyAlpha, nil
	case "insertion":
		return IndexStrategyInsertion, nil
	default:
		return 0, fmt.Errorf("unknown index strategy %q: want alpha or insertion", s)
	}
}

// IndexEntry holds a key and its position in the index.
type IndexEntry struct {
	Key      string
	Position int
}

// Indexer builds a positional index over a set of env keys.
type Indexer struct {
	strategy IndexStrategy
}

// NewIndexer creates an Indexer with the given strategy.
func NewIndexer(strategy IndexStrategy) *Indexer {
	return &Indexer{strategy: strategy}
}

// Build returns an ordered slice of IndexEntry for the provided map.
func (ix *Indexer) Build(m map[string]string) []IndexEntry {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	if ix.strategy == IndexStrategyAlpha {
		sort.Strings(keys)
	}

	entries := make([]IndexEntry, len(keys))
	for i, k := range keys {
		entries[i] = IndexEntry{Key: k, Position: i}
	}
	return entries
}

// Lookup returns the position of key in the index, or -1 if absent.
func Lookup(entries []IndexEntry, key string) int {
	for _, e := range entries {
		if e.Key == key {
			return e.Position
		}
	}
	return -1
}
