package sync

// MergeStrategy controls how existing .env values are handled during sync.
type MergeStrategy int

const (
	// MergeStrategyOverwrite replaces all existing values with Vault values.
	MergeStrategyOverwrite MergeStrategy = iota
	// MergeStrategyKeepExisting preserves local values when a key already exists.
	MergeStrategyKeepExisting
	// MergeStrategyVaultWins merges both sets but Vault values take precedence.
	MergeStrategyVaultWins
)

// Merge combines existing local env values with incoming Vault values
// according to the given MergeStrategy. It returns the merged map.
func Merge(existing, incoming map[string]string, strategy MergeStrategy) map[string]string {
	result := make(map[string]string, len(existing)+len(incoming))

	switch strategy {
	case MergeStrategyKeepExisting:
		// Start with incoming, then overwrite with existing (existing wins).
		for k, v := range incoming {
			result[k] = v
		}
		for k, v := range existing {
			result[k] = v
		}

	case MergeStrategyVaultWins:
		// Start with existing, then overwrite with incoming (vault wins).
		for k, v := range existing {
			result[k] = v
		}
		for k, v := range incoming {
			result[k] = v
		}

	case MergeStrategyOverwrite:
		fallthrough
	default:
		// Only vault values are used.
		for k, v := range incoming {
			result[k] = v
		}
	}

	return result
}

// MergeStrategyFromString parses a strategy name into a MergeStrategy.
// Returns MergeStrategyOverwrite and false if the name is unrecognised.
func MergeStrategyFromString(s string) (MergeStrategy, bool) {
	switch s {
	case "keep-existing", "keep":
		return MergeStrategyKeepExisting, true
	case "vault-wins", "vault":
		return MergeStrategyVaultWins, true
	case "overwrite", "":
		return MergeStrategyOverwrite, true
	default:
		return MergeStrategyOverwrite, false
	}
}
