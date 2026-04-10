// Package envfile provides utilities for reading and writing .env files.
package envfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Reader reads key-value pairs from an existing .env file.
type Reader struct {
	path string
}

// NewReader creates a new Reader for the given file path.
func NewReader(path string) *Reader {
	return &Reader{path: path}
}

// Read parses the .env file and returns a map of key-value pairs.
// Lines beginning with '#' are treated as comments and skipped.
// Lines that do not contain '=' are skipped.
// Returns an empty map if the file does not exist.
func (r *Reader) Read() (map[string]string, error) {
	f, err := os.Open(r.path)
	if os.IsNotExist(err) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("envfile: open %q: %w", r.path, err)
	}
	defer f.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// skip blank lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			continue
		}

		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])

		// strip surrounding quotes if present
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') ||
				(val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}

		if key != "" {
			result[key] = val
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("envfile: scan %q: %w", r.path, err)
	}

	return result, nil
}
