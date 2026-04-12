package sync

import "sort"

// DedupeStrategy controls how duplicate keys are resolved when merging
// secrets from multiple Vault paths into a single env file.
type DedupeStrategy int

const (
	// DedupeStrategyFirst keeps the first occurrence of a duplicate key.
	DedupeStrategyFirst DedupeStrategy = iota
	// DedupeStrategyLast keeps the last occurrence of a duplicate key.
	DedupeStrategyLast
	// DedupeStrategyError returns an error if duplicate keys are found.
	DedupeStrategyError
)

// DedupeResult holds the deduplicated map and any keys that were duplicated.
type DedupeResult struct {
	Secrets    map[string]string
	Duplicates []string
}

// Dedupe merges a slice of secret maps according to the given strategy.
// Maps are applied in order; earlier maps take priority under DedupeStrategyFirst.
func Dedupe(maps []map[string]string, strategy DedupeStrategy) (DedupeResult, error) {
	seen := make(map[string]bool)
	dupes := make(map[string]bool)
	result := make(map[string]string)

	for _, m := range maps {
		for k, v := range m {
			if _, exists := result[k]; exists {
				dupes[k] = true
				if strategy == DedupeStrategyError {
					return DedupeResult{}, &DuplicateKeyError{Key: k}
				}
				if strategy == DedupeStrategyLast {
					result[k] = v
				}
				// DedupeStrategyFirst: keep existing value, do nothing
			} else {
				result[k] = v
				seen[k] = true
			}
		}
	}

	dupeKeys := make([]string, 0, len(dupes))
	for k := range dupes {
		dupeKeys = append(dupeKeys, k)
	}
	sort.Strings(dupeKeys)

	return DedupeResult{Secrets: result, Duplicates: dupeKeys}, nil
}

// DuplicateKeyError is returned when DedupeStrategyError encounters a duplicate.
type DuplicateKeyError struct {
	Key string
}

func (e *DuplicateKeyError) Error() string {
	return "duplicate key: " + e.Key
}

// DedupeStrategyFromString parses a strategy name into a DedupeStrategy.
func DedupeStrategyFromString(s string) (DedupeStrategy, bool) {
	switch s {
	case "first":
		return DedupeStrategyFirst, true
	case "last":
		return DedupeStrategyLast, true
	case "error":
		return DedupeStrategyError, true
	default:
		return DedupeStrategyFirst, false
	}
}
