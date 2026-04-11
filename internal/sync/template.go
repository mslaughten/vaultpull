package sync

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"text/template"
)

// TemplateRenderer renders secret maps using Go text/template syntax.
// The template receives a map[string]string of key/value pairs.
type TemplateRenderer struct {
	tmpl   *template.Template
	writer io.Writer
}

// NewTemplateRenderer parses templateSrc and returns a TemplateRenderer that
// writes rendered output to w. If w is nil, os.Stdout is used.
func NewTemplateRenderer(templateSrc string, w io.Writer) (*TemplateRenderer, error) {
	if templateSrc == "" {
		return nil, fmt.Errorf("template source must not be empty")
	}
	tmpl, err := template.New("secrets").Option("missingkey=error").Parse(templateSrc)
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}
	if w == nil {
		w = os.Stdout
	}
	return &TemplateRenderer{tmpl: tmpl, writer: w}, nil
}

// Render executes the template with the provided secrets map and writes the
// result to the renderer's writer.
func (r *TemplateRenderer) Render(secrets map[string]string) error {
	var buf bytes.Buffer
	if err := r.tmpl.Execute(&buf, secrets); err != nil {
		return fmt.Errorf("render template: %w", err)
	}
	_, err := r.writer.Write(buf.Bytes())
	return err
}

// RenderToString executes the template and returns the result as a string.
func (r *TemplateRenderer) RenderToString(secrets map[string]string) (string, error) {
	var buf bytes.Buffer
	if err := r.tmpl.Execute(&buf, buf); err != nil {
		// re-execute against secrets
	}
	buf.Reset()
	if err := r.tmpl.Execute(&buf, secrets); err != nil {
		return "", fmt.Errorf("render template: %w", err)
	}
	return buf.String(), nil
}
