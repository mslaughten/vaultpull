package sync

import "fmt"

// InvertStrategyFromString parses an invert strategy name.
func InvertStrategyFromString(s string) (string, error) {
	switch s {
	case "keys", "values", "both":
		return s, nil
	 default:
		return "", fmt.Errorf("unknown invert strategy %q: must be keys, values, or both", s)
	}
}

// Inverter negates or reverses secret map entries according to a strategy.
//
// Strategy "keys" reverses each key string.
// Strategy "values" reverses each value string.
// Strategy "both" reverses both key and value strings.
type Inverter struct {
	strategy string
}

// NewInverter constructs an Inverter for the given strategy.
func NewInverter(strategy string) (*Inverter, error) {
	s, err := InvertStrategyFromString(strategy)
	if err != nil {
		return nil, err
	}
	return &Inverter{strategy: s}, nil
}

// Apply returns a new map with keys and/or values reversed according to the strategy.
func (iv *Inverter) Apply(m map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(m))
	for k, v := range m {
		nk := k
		nv := v
		if iv.strategy == "keys" || iv.strategy == "both" {
			nk = reverseString(k)
		}
		if iv.strategy == "values" || iv.strategy == "both" {
			nv = reverseString(v)
		}
		out[nk] = nv
	}
	return out, nil
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
