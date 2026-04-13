package sync

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// CompareResult holds the comparison outcome between a local .env file and Vault secrets.
type CompareResult struct {
	Path      string
	OnlyLocal []string
	OnlyVault []string
	Diverged  []string
	InSync    []string
}

// HasDrift returns true if there are any differences between local and Vault.
func (r *CompareResult) HasDrift() bool {
	return len(r.OnlyLocal) > 0 || len(r.OnlyVault) > 0 || len(r.Diverged) > 0
}

// Summary returns a human-readable summary of the comparison.
func (r *CompareResult) Summary() string {
	if !r.HasDrift() {
		return fmt.Sprintf("%s: in sync (%d keys)", r.Path, len(r.InSync))
	}
	parts := []string{fmt.Sprintf("%s: drift detected", r.Path)}
	if len(r.OnlyLocal) > 0 {
		parts = append(parts, fmt.Sprintf("  only-local: %s", strings.Join(r.OnlyLocal, ", ")))
	}
	if len(r.OnlyVault) > 0 {
		parts = append(parts, fmt.Sprintf("  only-vault: %s", strings.Join(r.OnlyVault, ", ")))
	}
	if len(r.Diverged) > 0 {
		parts = append(parts, fmt.Sprintf("  diverged:   %s", strings.Join(r.Diverged, ", ")))
	}
	return strings.Join(parts, "\n")
}

// Comparer compares a local env map against a Vault secret map.
type Comparer struct {
	w io.Writer
}

// NewComparer creates a Comparer that writes output to w.
// If w is nil, os.Stdout is used.
func NewComparer(w io.Writer) *Comparer {
	if w == nil {
		w = os.Stdout
	}
	return &Comparer{w: w}
}

// Compare produces a CompareResult for the given path and key maps.
func (c *Comparer) Compare(path string, local, vault map[string]string) CompareResult {
	res := CompareResult{Path: path}
	allKeys := map[string]struct{}{}
	for k := range local {
		allKeys[k] = struct{}{}
	}
	for k := range vault {
		allKeys[k] = struct{}{}
	}
	keys := make([]string, 0, len(allKeys))
	for k := range allKeys {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		lv, lok := local[k]
		vv, vok := vault[k]
		switch {
		case lok && !vok:
			res.OnlyLocal = append(res.OnlyLocal, k)
		case !lok && vok:
			res.OnlyVault = append(res.OnlyVault, k)
		case lv != vv:
			res.Diverged = append(res.Diverged, k)
		default:
			res.InSync = append(res.InSync, k)
		}
	}
	return res
}

// Print writes the summary of the result to the configured writer.
func (c *Comparer) Print(r CompareResult) {
	fmt.Fprintln(c.w, r.Summary())
}
