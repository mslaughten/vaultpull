package sync

import (
	"fmt"
	"path"
	"strings"
)

// Plan represents a set of sync operations to be performed.
type Plan struct {
	// Entries maps a local .env filename to the Vault secret path it should be populated from.
	Entries []PlanEntry
}

// PlanEntry describes a single path-to-file mapping.
type PlanEntry struct {
	VaultPath string
	EnvFile   string
}

// String returns a human-readable representation of the plan entry.
func (e PlanEntry) String() string {
	return fmt.Sprintf("%s -> %s", e.VaultPath, e.EnvFile)
}

// String returns a human-readable summary of the plan, listing all entries.
func (p Plan) String() string {
	if len(p.Entries) == 0 {
		return "(empty plan)"
	}
	var sb strings.Builder
	for _, e := range p.Entries {
		sb.WriteString(e.String())
		sb.WriteByte('\n')
	}
	return strings.TrimRight(sb.String(), "\n")
}

// BuildPlan constructs a Plan from a list of Vault secret paths and an optional
// namespace prefix to strip when deriving local file names.
//
// Each path is converted to a .env filename using the last path component.
// Duplicate filenames are disambiguated by appending a numeric suffix.
func BuildPlan(paths []string, namespace string) Plan {
	seen := make(map[string]int)
	entries := make([]PlanEntry, 0, len(paths))

	for _, p := range paths {
		stripped := stripNamespacePrefix(p, namespace)
		base := path.Base(stripped)
		if base == "." || base == "" {
			continue
		}

		fileName := envFileName(base)

		count := seen[fileName]
		seen[fileName]++

		if count > 0 {
			ext := path.Ext(fileName)
			name := strings.TrimSuffix(fileName, ext)
			fileName = fmt.Sprintf("%s_%d%s", name, count, ext)
		}

		entries = append(entries, PlanEntry{
			VaultPath: p,
			EnvFile:   fileName,
		})
	}

	return Plan{Entries: entries}
}

// stripNamespacePrefix removes a leading namespace segment from a vault path.
func stripNamespacePrefix(vaultPath, namespace string) string {
	if namespace == "" {
		return vaultPath
	}
	prefix := strings.TrimSuffix(namespace, "/") + "/"
	return strings.TrimPrefix(vaultPath, prefix)
}
