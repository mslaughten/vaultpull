package sync

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Archiver bundles one or more .env files into a timestamped ZIP archive.
// It is useful for creating point-in-time backups before a bulk sync or
// rotation operation.
type Archiver struct {
	outDir  string
	prefix  string
	clock   func() time.Time
}

// ArchiverOption configures an Archiver.
type ArchiverOption func(*Archiver)

// WithArchivePrefix sets the filename prefix for the generated archive.
// Defaults to "vaultpull".
func WithArchivePrefix(prefix string) ArchiverOption {
	return func(a *Archiver) {
		if prefix != "" {
			a.prefix = prefix
		}
	}
}

// withArchiveClock overrides the time source (test helper).
func withArchiveClock(fn func() time.Time) ArchiverOption {
	return func(a *Archiver) { a.clock = fn }
}

// NewArchiver creates an Archiver that writes archives to outDir.
// outDir is created if it does not exist.
func NewArchiver(outDir string, opts ...ArchiverOption) (*Archiver, error) {
	if outDir == "" {
		return nil, fmt.Errorf("archive: outDir must not be empty")
	}
	if err := os.MkdirAll(outDir, 0o700); err != nil {
		return nil, fmt.Errorf("archive: create outDir: %w", err)
	}
	a := &Archiver{
		outDir: outDir,
		prefix: "vaultpull",
		clock:  time.Now,
	}
	for _, o := range opts {
		o(a)
	}
	return a, nil
}

// Archive writes the given files into a ZIP archive named
// "<prefix>-<timestamp>.zip" inside outDir and returns the full path.
// files must be a non-empty slice of existing file paths.
func (a *Archiver) Archive(files []string) (string, error) {
	if len(files) == 0 {
		return "", fmt.Errorf("archive: no files provided")
	}

	ts := a.clock().UTC().Format("20060102-150405")
	archiveName := fmt.Sprintf("%s-%s.zip", a.prefix, ts)
	archivePath := filepath.Join(a.outDir, archiveName)

	zf, err := os.Create(archivePath)
	if err != nil {
		return "", fmt.Errorf("archive: create zip file: %w", err)
	}
	defer zf.Close()

	zw := zip.NewWriter(zf)
	defer zw.Close()

	// Sort for deterministic archive order.
	sorted := make([]string, len(files))
	copy(sorted, files)
	sort.Strings(sorted)

	for _, src := range sorted {
		if err := addFileToZip(zw, src); err != nil {
			return "", err
		}
	}

	return archivePath, nil
}

// addFileToZip adds a single file to the ZIP writer using only the base name
// as the entry name so archives remain portable.
func addFileToZip(zw *zip.Writer, src string) error {
	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("archive: open %s: %w", src, err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return fmt.Errorf("archive: stat %s: %w", src, err)
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("archive: build header for %s: %w", src, err)
	}
	header.Name = filepath.Base(src)
	header.Method = zip.Deflate

	w, err := zw.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("archive: create zip entry for %s: %w", src, err)
	}

	if _, err := io.Copy(w, f); err != nil {
		return fmt.Errorf("archive: write %s: %w", src, err)
	}
	return nil
}
