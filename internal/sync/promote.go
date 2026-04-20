package sync

import (
	"fmt"
	"sort"
)

// PromoteStrategy controls how keys are promoted between environments.
type PromoteStrategy int

const (
	PromoteStrategyMissing PromoteStrategy = iota // only promote keys absent in dst
	PromoteStrategyAll                            // promote all keys, overwriting dst
)

// PromoteSummary describes the result of a promotion.
type PromoteSummary struct {
	Promoted []string
	Skipped  []string
}

func (s PromoteSummary) String() string {
	return fmt.Sprintf("promoted=%d skipped=%d", len(s.Promoted), len(s.Skipped))
}

// PromoteStrategyFromString parses a strategy name.
func PromoteStrategyFromString(s string) (PromoteStrategy, error) {
	switch s {
	case "missing", "":
		return PromoteStrategyMissing, nil
	case "all":
		return PromoteStrategyAll, nil
	default:
		return 0, fmt.Errorf("unknown promote strategy %q: want missing|all", s)
	}
}

// Promoter copies keys from a source map into a destination map.
type Promoter struct {
	strategy PromoteStrategy
}

// NewPromoter constructs a Promoter with the given strategy.
func NewPromoter(strategy PromoteStrategy) *Promoter {
	return &Promoter{strategy: strategy}
}

// Apply promotes keys from src into dst, returning a summary and the merged map.
func (p *Promoter) Apply(src, dst map[string]string) (map[string]string, PromoteSummary, error) {
	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	var summary PromoteSummary

	keys := make([]string, 0, len(src))
	for k := range src {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		_, exists := out[k]
		if exists && p.strategy == PromoteStrategyMissing {
			summary.Skipped = append(summary.Skipped, k)
			continue
		}
		out[k] = src[k]
		summary.Promoted = append(summary.Promoted, k)
	}

	return out, summary, nil
}
