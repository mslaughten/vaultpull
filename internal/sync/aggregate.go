package sync

import (
	"fmt"
	"sort"
	"strings"
)

// AggregateStrategy controls how multiple values are combined into one.
type AggregateStrategy string

const (
	AggregateConcat AggregateStrategy = "concat"
	AggregateCount  AggregateStrategy = "count"
	AggregateUnique AggregateStrategy = "unique"
)

// AggregateStrategyFromString parses a strategy name, returning an error for
// unknown values.
func AggregateStrategyFromString(s string) (AggregateStrategy, error) {
	switch AggregateStrategy(strings.ToLower(s)) {
	case AggregateConcat, AggregateCount, AggregateUnique:
		return AggregateStrategy(strings.ToLower(s)), nil
	}
	return "", fmt.Errorf("unknown aggregate strategy %q: want concat|count|unique", s)
}

// Aggregator combines keys sharing a common prefix into a single output key.
type Aggregator struct {
	prefix    string
	outKey    string
	separator string
	strategy  AggregateStrategy
}

// NewAggregator constructs an Aggregator. prefix is the key prefix to match,
// outKey is the destination key written to the map, sep is the joining
// separator used by the concat strategy.
func NewAggregator(prefix, outKey, sep string, strategy AggregateStrategy) (*Aggregator, error) {
	if prefix == "" {
		return nil, fmt.Errorf("aggregate: prefix must not be empty")
	}
	if outKey == "" {
		return nil, fmt.Errorf("aggregate: outKey must not be empty")
	}
	if sep == "" {
		sep = ","
	}
	return &Aggregator{prefix: prefix, outKey: outKey, separator: sep, strategy: strategy}, nil
}

// Apply scans m for keys beginning with prefix, aggregates their values
// according to the strategy, writes the result under outKey, and removes the
// matched source keys.
func (a *Aggregator) Apply(m map[string]string) (map[string]string, error) {
	var matched []string
	for k := range m {
		if strings.HasPrefix(k, a.prefix) {
			matched = append(matched, k)
		}
	}
	if len(matched) == 0 {
		return m, nil
	}
	sort.Strings(matched)

	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}

	switch a.strategy {
	case AggregateCount:
		out[a.outKey] = fmt.Sprintf("%d", len(matched))
	case AggregateUnique:
		seen := map[string]struct{}{}
		var vals []string
		for _, k := range matched {
			v := m[k]
			if _, ok := seen[v]; !ok {
				seen[v] = struct{}{}
				vals = append(vals, v)
			}
		}
		out[a.outKey] = strings.Join(vals, a.separator)
	default: // concat
		var vals []string
		for _, k := range matched {
			vals = append(vals, m[k])
		}
		out[a.outKey] = strings.Join(vals, a.separator)
	}

	for _, k := range matched {
		delete(out, k)
	}
	return out, nil
}
