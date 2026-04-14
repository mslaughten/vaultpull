package sync

import "fmt"

// CoalesceStrategy determines how the first non-empty value is selected.
type CoalesceStrategy int

const (
	CoalesceFirst CoalesceStrategy = iota
	CoalesceLongest
	CoalesceNonZero
)

// CoalesceStrategyFromString parses a strategy name.
func CoalesceStrategyFromString(s string) (CoalesceStrategy, error) {
	switch s {
	case "first", "":
		return CoalesceFirst, nil
	case "longest":
		return CoalesceLongest, nil
	case "nonzero":
		return CoalesceNonZero, nil
	default:
		return CoalesceFirst, fmt.Errorf("unknown coalesce strategy %q: want first|longest|nonzero", s)
	}
}

// Coalescer picks the winning value for duplicate keys across multiple maps.
type Coalescer struct {
	strategy CoalesceStrategy
}

// NewCoalescer returns a Coalescer for the given strategy name.
func NewCoalescer(strategy string) (*Coalescer, error) {
	s, err := CoalesceStrategyFromString(strategy)
	if err != nil {
		return nil, err
	}
	return &Coalescer{strategy: s}, nil
}

// Apply merges the provided maps in order, resolving duplicate keys using the
// configured strategy. Later maps have lower priority than earlier ones.
func (c *Coalescer) Apply(maps ...map[string]string) map[string]string {
	result := make(map[string]string)

	for _, m := range maps {
		for k, v := range m {
			existing, exists := result[k]
			if !exists {
				result[k] = v
				continue
			}
			switch c.strategy {
			case CoalesceFirst:
				// keep existing — first wins
			case CoalesceLongest:
				if len(v) > len(existing) {
					result[k] = v
				}
			case CoalesceNonZero:
				if existing == "" && v != "" {
					result[k] = v
				}
			}
		}
	}
	return result
}
