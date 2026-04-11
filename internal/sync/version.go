package sync

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

// SecretVersion represents a point-in-time snapshot of a secret's metadata.
type SecretVersion struct {
	Path      string
	Version   int
	CreatedAt time.Time
	Deleted   bool
}

// VersionLister can retrieve version history for a secret path.
type VersionLister interface {
	ListVersions(ctx context.Context, mount, path string) ([]SecretVersion, error)
}

// VersionPrinter writes version history to an output stream.
type VersionPrinter struct {
	out io.Writer
}

// NewVersionPrinter creates a VersionPrinter writing to w.
// If w is nil it defaults to os.Stdout.
func NewVersionPrinter(w io.Writer) *VersionPrinter {
	if w == nil {
		w = os.Stdout
	}
	return &VersionPrinter{out: w}
}

// Print fetches and renders version history for the given mount + path.
func (vp *VersionPrinter) Print(ctx context.Context, lister VersionLister, mount, path string) error {
	versions, err := lister.ListVersions(ctx, mount, path)
	if err != nil {
		return fmt.Errorf("list versions %s/%s: %w", mount, path, err)
	}
	if len(versions) == 0 {
		fmt.Fprintf(vp.out, "no versions found for %s/%s\n", mount, path)
		return nil
	}
	fmt.Fprintf(vp.out, "versions for %s/%s:\n", mount, path)
	for _, v := range versions {
		status := "active"
		if v.Deleted {
			status = "deleted"
		}
		fmt.Fprintf(vp.out, "  v%-4d  %-8s  %s\n", v.Version, status, v.CreatedAt.Format(time.RFC3339))
	}
	return nil
}
