package sync

import (
	"bytes"
	"strings"
	"testing"
)

func TestConfirmPrompt_Yes(t *testing.T) {
	for _, input := range []string{"y", "Y", "yes", "YES", "Yes"} {
		r := strings.NewReader(input)
		var w bytes.Buffer
		ok, err := ConfirmPrompt(&w, r, "Continue?")
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", input, err)
		}
		if !ok {
			t.Errorf("input %q: expected true, got false", input)
		}
	}
}

func TestConfirmPrompt_No(t *testing.T) {
	for _, input := range []string{"n", "N", "no", "", "maybe"} {
		r := strings.NewReader(input)
		var w bytes.Buffer
		ok, err := ConfirmPrompt(&w, r, "Continue?")
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", input, err)
		}
		if ok {
			t.Errorf("input %q: expected false, got true", input)
		}
	}
}

func TestConfirmPrompt_EOF(t *testing.T) {
	r := strings.NewReader("")
	var w bytes.Buffer
	ok, err := ConfirmPrompt(&w, r, "Continue?")
	if err != nil {
		t.Fatalf("unexpected error on EOF: %v", err)
	}
	if ok {
		t.Error("expected false on EOF, got true")
	}
}

func TestConfirmDiff_NoChanges(t *testing.T) {
	d := DiffResult{}
	var w bytes.Buffer
	r := strings.NewReader("y")
	ok, err := ConfirmDiff(d, &w, r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected false when no changes, got true")
	}
	if !strings.Contains(w.String(), "No changes") {
		t.Errorf("expected 'No changes' in output, got: %s", w.String())
	}
}

func TestConfirmDiff_WithChanges_Confirmed(t *testing.T) {
	d := DiffResult{
		Added:   []string{"API_KEY"},
		Changed: []string{"DB_URL"},
	}
	var w bytes.Buffer
	r := strings.NewReader("y")
	ok, err := ConfirmDiff(d, &w, r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected true after confirming changes")
	}
	output := w.String()
	if !strings.Contains(output, "+ API_KEY") {
		t.Errorf("expected added key in output, got: %s", output)
	}
	if !strings.Contains(output, "~ DB_URL") {
		t.Errorf("expected changed key in output, got: %s", output)
	}
}
