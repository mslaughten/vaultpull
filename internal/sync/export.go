package sync

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
)

// ExportFormat controls the output format for exported secrets.
type ExportFormat string

const (
	ExportFormatJSON ExportFormat = "json"
	ExportFormatDotEnv ExportFormat = "dotenv"
)

// Exporter writes secrets to an output stream in the requested format.
type Exporter struct {
	format ExportFormat
	out    io.Writer
}

// NewExporter creates an Exporter writing to out in the given format.
func NewExporter(format ExportFormat, out io.Writer) (*Exporter, error) {
	switch format {
	case ExportFormatJSON, ExportFormatDotEnv:
		// valid
	default:
		return nil, fmt.Errorf("unsupported export format: %q", format)
	}
	if out == nil {
		out = os.Stdout
	}
	return &Exporter{format: format, out: out}, nil
}

// Write serialises secrets to the configured output.
func (e *Exporter) Write(secrets map[string]string) error {
	switch e.format {
	case ExportFormatJSON:
		return e.writeJSON(secrets)
	case ExportFormatDotEnv:
		return e.writeDotEnv(secrets)
	}
	return fmt.Errorf("unknown format: %q", e.format)
}

func (e *Exporter) writeJSON(secrets map[string]string) error {
	enc := json.NewEncoder(e.out)
	enc.SetIndent("", "  ")
	return enc.Encode(secrets)
}

func (e *Exporter) writeDotEnv(secrets map[string]string) error {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if _, err := fmt.Fprintf(e.out, "%s=%q\n", k, secrets[k]); err != nil {
			return err
		}
	}
	return nil
}
