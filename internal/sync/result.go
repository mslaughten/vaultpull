package sync

import (
	"fmt"
	"strings"
)

// Result summarises the outcome of a sync run.
type Result struct {
	Total   int
	Written []string
	Errors  []error
}

// HasErrors returns true when at least one path failed.
func (r Result) HasErrors() bool {
	return len(r.Errors) > 0
}

// Summary returns a human-readable one-line summary.
func (r Result) Summary() string {
	return fmt.Sprintf(
		"synced %d/%d secret(s), %d error(s)",
		len(r.Written), r.Total, len(r.Errors),
	)
}

// ErrorMessages returns all error messages joined by newlines.
func (r Result) ErrorMessages() string {
	msgs := make([]string, 0, len(r.Errors))
	for _, e := range r.Errors {
		msgs = append(msgs, e.Error())
	}
	return strings.Join(msgs, "\n")
}
