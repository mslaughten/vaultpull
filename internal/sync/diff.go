package sync

import "sort"

// DiffResult holds the comparison between existing and new secrets.
type DiffResult struct {
	Added   []string
	Removed []string
	Changed []string
	Unchanged []string
}

// HasChanges returns true if there are any additions, removals, or modifications.
func (d *DiffResult) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}

// Summary returns a human-readable summary of the diff.
func (d *DiffResult) Summary() string {
	if !d.HasChanges() {
		return "no changes detected"
	}
	return fmt.Sprintf("+%d added, -%d removed, ~%d changed",
		len(d.Added), len(d.Removed), len(d.Changed))
}

// Diff compares an existing env map with an incoming secrets map and returns
// a DiffResult describing what has changed.
func Diff(existing, incoming map[string]string) DiffResult {
	result := DiffResult{}

	existingKeys := make(map[string]struct{}, len(existing))
	for k := range existing {
		existingKeys[k] = struct{}{}
	}

	for k, newVal := range incoming {
		if oldVal, ok := existing[k]; ok {
			if oldVal != newVal {
				result.Changed = append(result.Changed, k)
			} else {
				result.Unchanged = append(result.Unchanged, k)
			}
			delete(existingKeys, k)
		} else {
			result.Added = append(result.Added, k)
		}
	}

	for k := range existingKeys {
		result.Removed = append(result.Removed, k)
	}

	sort.Strings(result.Added)
	sort.Strings(result.Removed)
	sort.Strings(result.Changed)
	sort.Strings(result.Unchanged)

	return result
}
