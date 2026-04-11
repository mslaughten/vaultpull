package sync

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

type mockVersionLister struct {
	versions []SecretVersion
	err      error
}

func (m *mockVersionLister) ListVersions(_ context.Context, _, _ string) ([]SecretVersion, error) {
	return m.versions, m.err
}

func TestVersionPrinter_NilWriter_DefaultsToStdout(t *testing.T) {
	vp := NewVersionPrinter(nil)
	if vp.out == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestVersionPrinter_PrintsVersions(t *testing.T) {
	now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	lister := &mockVersionLister{
		versions: []SecretVersion{
			{Path: "secret/db", Version: 1, CreatedAt: now, Deleted: false},
			{Path: "secret/db", Version: 2, CreatedAt: now.Add(time.Hour), Deleted: true},
		},
	}
	var buf bytes.Buffer
	vp := NewVersionPrinter(&buf)
	err := vp.Print(context.Background(), lister, "secret", "db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "v1") {
		t.Errorf("expected v1 in output, got: %s", out)
	}
	if !strings.Contains(out, "deleted") {
		t.Errorf("expected 'deleted' in output, got: %s", out)
	}
	if !strings.Contains(out, "active") {
		t.Errorf("expected 'active' in output, got: %s", out)
	}
}

func TestVersionPrinter_EmptyVersions(t *testing.T) {
	lister := &mockVersionLister{versions: []SecretVersion{}}
	var buf bytes.Buffer
	vp := NewVersionPrinter(&buf)
	err := vp.Print(context.Background(), lister, "secret", "db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no versions found") {
		t.Errorf("expected 'no versions found', got: %s", buf.String())
	}
}

func TestVersionPrinter_ListerError(t *testing.T) {
	lister := &mockVersionLister{err: errors.New("vault unavailable")}
	var buf bytes.Buffer
	vp := NewVersionPrinter(&buf)
	err := vp.Print(context.Background(), lister, "secret", "db")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "vault unavailable") {
		t.Errorf("unexpected error message: %v", err)
	}
}
