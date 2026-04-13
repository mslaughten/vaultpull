package sync

import (
	"fmt"
	"sort"
	"strings"
)

// GroupStrategy determines how secrets are grouped into output files.
type GroupStrategy string

const (
	GroupByPrefix    GroupStrategy = "prefix"
	GroupByNamespace GroupStrategy = "namespace"
	GroupFlat        GroupStrategy = "flat"
)

// GroupEntry holds a named group of key-value pairs.
type GroupEntry struct {
	Name   string
	Values map[string]string
}

// Grouper splits a flat secret map into named groups.
type Grouper struct {
	strategy  GroupStrategy
	delimiter string
}

// NewGrouper creates a Grouper with the given strategy and delimiter.
// delimiter is used to split keys (e.g. "_" or "/").
func NewGrouper(strategy GroupStrategy, delimiter string) (*Grouper, error) {
	switch strategy {
	case GroupByPrefix, GroupByNamespace, GroupFlat:
		// valid
	default:
		return nil, fmt.Errorf("group: unknown strategy %q", strategy)
	}
	if delimiter == "" {
		delimiter = "_"
	}
	return &Grouper{strategy: strategy, delimiter: delimiter}, nil
}

// Apply partitions secrets into groups and returns them sorted by name.
func (g *Grouper) Apply(secrets map[string]string) []GroupEntry {
	if g.strategy == GroupFlat {
		return []GroupEntry{{Name: "default", Values: copyMap(secrets)}}
	}

	buckets := make(map[string]map[string]string)
	for k, v := range secrets {
		groupName, localKey := g.split(k)
		if _, ok := buckets[groupName]; !ok {
			buckets[groupName] = make(map[string]string)
		}
		buckets[groupName][localKey] = v
	}

	entries := make([]GroupEntry, 0, len(buckets))
	for name, vals := range buckets {
		entries = append(entries, GroupEntry{Name: name, Values: vals})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})
	return entries
}

func (g *Grouper) split(key string) (group, local string) {
	idx := strings.Index(key, g.delimiter)
	if idx < 0 {
		return "default", key
	}
	return key[:idx], key[idx+len(g.delimiter):]
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
