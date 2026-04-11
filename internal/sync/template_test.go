package sync

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewTemplateRenderer_EmptySource_ReturnsError(t *testing.T) {
	_, err := NewTemplateRenderer("", nil)
	if err == nil {
		t.Fatal("expected error for empty template source")
	}
}

func TestNewTemplateRenderer_InvalidTemplate_ReturnsError(t *testing.T) {
	_, err := NewTemplateRenderer("{{ .Unclosed", nil)
	if err == nil {
		t.Fatal("expected parse error for invalid template")
	}
}

func TestNewTemplateRenderer_NilWriter_DefaultsToStdout(t *testing.T) {
	r, err := NewTemplateRenderer("hello", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.writer == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestRender_InterpolatesSecrets(t *testing.T) {
	const src = `DB_HOST={{ index . "DB_HOST" }}
DB_PORT={{ index . "DB_PORT" }}`
	var buf bytes.Buffer
	r, err := NewTemplateRenderer(src, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	if err := r.Render(secrets); err != nil {
		t.Fatalf("Render error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Errorf("expected DB_HOST=localhost in output, got: %s", out)
	}
	if !strings.Contains(out, "DB_PORT=5432") {
		t.Errorf("expected DB_PORT=5432 in output, got: %s", out)
	}
}

func TestRender_MissingKey_ReturnsError(t *testing.T) {
	const src = `{{ index . "MISSING_KEY" }}`
	var buf bytes.Buffer
	r, err := NewTemplateRenderer(src, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = r.Render(map[string]string{})
	if err == nil {
		t.Fatal("expected error for missing key with missingkey=error")
	}
}

func TestRenderToString_ReturnsRenderedString(t *testing.T) {
	const src = `TOKEN={{ index . "TOKEN" }}`
	r, err := NewTemplateRenderer(src, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := r.RenderToString(map[string]string{"TOKEN": "abc123"})
	if err != nil {
		t.Fatalf("RenderToString error: %v", err)
	}
	if out != "TOKEN=abc123" {
		t.Errorf("expected 'TOKEN=abc123', got %q", out)
	}
}
