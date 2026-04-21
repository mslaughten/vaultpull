package sync

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// ObserveStrategy controls what the observer reports.
type ObserveStrategy string

const (
	ObserveAll     ObserveStrategy = "all"
	ObserveChanged ObserveStrategy = "changed"
	ObserveMissing ObserveStrategy = "missing"
)

// ObserveStrategyFromString parses a strategy string.
func ObserveStrategyFromString(s string) (ObserveStrategy, error) {
	switch strings.ToLower(s) {
	case "all":
		return ObserveAll, nil
	case "changed":
		return ObserveChanged, nil
	case "missing":
		return ObserveMissing, nil
	default:
		return "", fmt.Errorf("unknown observe strategy %q: want all|changed|missing", s)
	}
}

// ObserveResult holds the output of an observation pass.
type ObserveResult struct {
	Key      string
	Status   string
	Current  string
	Previous string
}

// Observer compares a current env map against a reference map and reports
// keys according to the chosen strategy.
type Observer struct {
	strategy  ObserveStrategy
	reference map[string]string
	w         io.Writer
}

// NewObserver creates an Observer. reference is the baseline map to compare
// against. w defaults to os.Stdout when nil.
func NewObserver(reference map[string]string, strategy ObserveStrategy, w io.Writer) (*Observer, error) {
	if reference == nil {
		return nil, fmt.Errorf("observe: reference map must not be nil")
	}
	if w == nil {
		w = os.Stdout
	}
	return &Observer{strategy: strategy, reference: reference, w: w}, nil
}

// Observe inspects current against the reference map and writes a report.
func (o *Observer) Observe(current map[string]string) ([]ObserveResult, error) {
	keys := make(map[string]struct{})
	for k := range o.reference {
		keys[k] = struct{}{}
	}
	for k := range current {
		keys[k] = struct{}{}
	}

	sorted := make([]string, 0, len(keys))
	for k := range keys {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)

	var results []ObserveResult
	for _, k := range sorted {
		prev, hasPrev := o.reference[k]
		curr, hasCurr := current[k]

		var status string
		switch {
		case hasPrev && !hasCurr:
			status = "missing"
		case !hasPrev && hasCurr:
			status = "added"
		case prev != curr:
			status = "changed"
		default:
			status = "unchanged"
		}

		if o.strategy == ObserveChanged && status != "changed" {
			continue
		}
		if o.strategy == ObserveMissing && status != "missing" {
			continue
		}

		results = append(results, ObserveResult{Key: k, Status: status, Current: curr, Previous: prev})
	}

	for _, r := range results {
		fmt.Fprintf(o.w, "%-30s %-10s\n", r.Key, r.Status)
	}
	return results, nil
}
