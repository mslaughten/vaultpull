package sync

import (
	"bytes"
	"os"
	"testing"
)

func TestNewSnapshotter_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	sub := dir + "/snaps"
	_, err := NewSnapshotter(sub)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(sub); os.IsNotExist(err) {
		t.Fatal("expected directory to be created")
	}
}

func TestSnapshotter_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	ss, _ := NewSnapshotter(dir)

	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := ss.Save("test", secrets); err != nil {
		t.Fatalf("save: %v", err)
	}

	entry, err := ss.Load("test")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if entry == nil {
		t.Fatal("expected entry, got nil")
	}
	if entry.Label != "test" {
		t.Errorf("label: got %q, want %q", entry.Label, "test")
	}
	if entry.Secrets["FOO"] != "bar" {
		t.Errorf("FOO: got %q, want %q", entry.Secrets["FOO"], "bar")
	}
}

func TestSnapshotter_Load_Missing_ReturnsNil(t *testing.T) {
	dir := t.TempDir()
	ss, _ := NewSnapshotter(dir)

	entry, err := ss.Load("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry != nil {
		t.Fatal("expected nil for missing snapshot")
	}
}

func TestSnapshotter_Delete_RemovesFile(t *testing.T) {
	dir := t.TempDir()
	ss, _ := NewSnapshotter(dir)

	_ = ss.Save("todelete", map[string]string{"K": "V"})
	if err := ss.Delete("todelete"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	entry, _ := ss.Load("todelete")
	if entry != nil {
		t.Fatal("expected snapshot to be deleted")
	}
}

func TestSnapshotter_Delete_MissingFile_NoError(t *testing.T) {
	dir := t.TempDir()
	ss, _ := NewSnapshotter(dir)
	if err := ss.Delete("ghost"); err != nil {
		t.Fatalf("unexpected error on missing delete: %v", err)
	}
}

func TestSnapshotter_Print_Output(t *testing.T) {
	dir := t.TempDir()
	ss, _ := NewSnapshotter(dir)

	secrets := map[string]string{"ALPHA": "one"}
	_ = ss.Save("mysnap", secrets)
	entry, _ := ss.Load("mysnap")

	var buf bytes.Buffer
	ss.Print(entry, &buf)
	out := buf.String()
	if out == "" {
		t.Fatal("expected non-empty output")
	}
	if !bytes.Contains([]byte(out), []byte("mysnap")) {
		t.Errorf("expected label in output, got: %s", out)
	}
}
