// Package sync provides utilities for synchronising HashiCorp Vault secrets
// with local .env files.
//
// # Export
//
// The Exporter type serialises a flat map of secrets to an io.Writer in one
// of two formats:
//
//   - ExportFormatJSON   – pretty-printed JSON object
//   - ExportFormatDotEnv – KEY="value" lines, sorted alphabetically
//
// Example:
//
//	ex, _ := sync.NewExporter(sync.ExportFormatJSON, os.Stdout)
//	ex.Write(secrets)
package sync
