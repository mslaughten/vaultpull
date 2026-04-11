package sync

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// ConfirmPrompt asks the user to confirm an action interactively.
// It writes the prompt to w and reads from r. Returns true if the
// user confirms with "y" or "yes" (case-insensitive).
func ConfirmPrompt(w io.Writer, r io.Reader, message string) (bool, error) {
	fmt.Fprintf(w, "%s [y/N]: ", message)

	var response string
	_, err := fmt.Fscan(r, &response)
	if err != nil {
		// EOF or no input treated as "no"
		if err == io.EOF {
			fmt.Fprintln(w)
			return false, nil
		}
		return false, fmt.Errorf("reading confirmation: %w", err)
	}

	normalized := strings.ToLower(strings.TrimSpace(response))
	return normalized == "y" || normalized == "yes", nil
}

// ConfirmDiff prints a summary of pending changes and prompts the user
// to confirm before proceeding. Returns true if the user approves.
func ConfirmDiff(d DiffResult, w io.Writer, r io.Reader) (bool, error) {
	if w == nil {
		w = os.Stdout
	}
	if r == nil {
		r = os.Stdin
	}

	fmt.Fprintln(w, "Pending changes:")
	for _, e := range d.Added {
		fmt.Fprintf(w, "  + %s\n", e)
	}
	for _, e := range d.Removed {
		fmt.Fprintf(w, "  - %s\n", e)
	}
	for _, e := range d.Changed {
		fmt.Fprintf(w, "  ~ %s\n", e)
	}
	fmt.Fprintln(w, d.Summary())

	if !d.HasChanges() {
		fmt.Fprintln(w, "No changes to apply.")
		return false, nil
	}

	return ConfirmPrompt(w, r, "Apply these changes?")
}
