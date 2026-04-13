package sync

import (
	"fmt"
	"strings"
)

// Cloner copies secrets from one Vault path prefix to another, optionally
// transforming key names with a find-and-replace on the destination prefix.
type Cloner struct {
	src    string
	dst    string
	dryRun bool
}

// CloneResult holds the outcome of a single key clone operation.
type CloneResult struct {
	SrcKey string
	DstKey string
	Err    error
}

// NewCloner creates a Cloner that copies keys whose names start with src,
// rewriting the prefix to dst. Returns an error if either prefix is empty.
func NewCloner(src, dst string, dryRun bool) (*Cloner, error) {
	src = strings.TrimSpace(src)
	dst = strings.TrimSpace(dst)
	if src == "" {
		return nil, fmt.Errorf("clone: source prefix must not be empty")
	}
	if dst == "" {
		return nil, fmt.Errorf("clone: destination prefix must not be empty")
	}
	return &Cloner{src: src, dst: dst, dryRun: dryRun}, nil
}

// Apply iterates over secrets, clones every key that starts with src into a
// new key with src replaced by dst, and returns the combined map plus results.
// In dry-run mode the returned map is a copy and originals are not modified.
func (c *Cloner) Apply(secrets map[string]string) (map[string]string, []CloneResult) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	var results []CloneResult
	for k, v := range secrets {
		if !strings.HasPrefix(k, c.src) {
			continue
		}
		newKey := c.dst + strings.TrimPrefix(k, c.src)
		results = append(results, CloneResult{SrcKey: k, DstKey: newKey})
		if !c.dryRun {
			out[newKey] = v
		}
	}
	return out, results
}

// Summary returns a human-readable description of the clone results.
func CloneSummary(results []CloneResult) string {
	if len(results) == 0 {
		return "clone: no keys matched"
	}
	var sb strings.Builder
	for _, r := range results {
		if r.Err != nil {
			fmt.Fprintf(&sb, "  ERROR %s -> %s: %v\n", r.SrcKey, r.DstKey, r.Err)
		} else {
			fmt.Fprintf(&sb, "  %s -> %s\n", r.SrcKey, r.DstKey)
		}
	}
	return sb.String()
}
