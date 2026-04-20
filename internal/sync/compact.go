package sync

import "fmt"

// CompactStrategy controls how empty or blank values are handled during compaction.
type CompactStrategy int

const (
	// CompactEmpty removes keys whose values are the empty string.
	CompactEmpty CompactStrategy = iota
	// CompactBlank removes keys whose values are empty or contain only whitespace.
	CompactBlank
)

// CompactStrategyFromString parses a strategy name into a CompactStrategy.
func CompactStrategyFromString(s string) (CompactStrategy, error) {
	switch s {
	case "empty":
		return CompactEmpty, nil
	case "blank":
		return CompactBlank, nil
	default:
		return 0, fmt.Errorf("compact: unknown strategy %q (want empty|blank)", s)
	}
}

// Compacter removes keys from a secret map according to a CompactStrategy.
type Compacter struct {
	strategy CompactStrategy
}

// NewCompacter creates a Compacter for the given strategy name.
func NewCompacter(strategy string) (*Compacter, error) {
	s, err := CompactStrategyFromString(strategy)
	if err != nil {
		return nil, err
	}
	return &Compacter{strategy: s}, nil
}

// Apply returns a copy of m with keys removed according to the strategy.
func (c *Compacter) Apply(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		if c.shouldRemove(v) {
			continue
		}
		out[k] = v
	}
	return out
}

func (c *Compacter) shouldRemove(v string) bool {
	switch c.strategy {
	case CompactBlank:
		for _, ch := range v {
			if ch != ' ' && ch != '\t' && ch != '\n' && ch != '\r' {
				return false
			}
		}
		return true
	default: // CompactEmpty
		return v == ""
	}
}
