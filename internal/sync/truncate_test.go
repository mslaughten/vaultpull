package sync

import (
	"testing"
)

func TestNewTruncator_InvalidMaxLen(t *testing.T) {
	_, err := NewTruncator(0, "end")
	if err == nil {
		t.Fatal("expected error for maxLen=0")
	}
}

func TestNewTruncator_UnknownMode(t *testing.T) {
	_, err := NewTruncator(10, "center")
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestNewTruncator_Valid(t *testing.T) {
	tr, err := NewTruncator(10, "end")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil Truncator")
	}
}

func TestTruncator_Apply_ShortValueUnchanged(t *testing.T) {
	tr, _ := NewTruncator(20, "end")
	in := map[string]string{"KEY": "short"}
	out := tr.Apply(in)
	if out["KEY"] != "short" {
		t.Errorf("expected %q, got %q", "short", out["KEY"])
	}
}

func TestTruncator_ModeEnd(t *testing.T) {
	tr, _ := NewTruncator(8, "end")
	out := tr.Apply(map[string]string{"K": "abcdefghij"})
	want := "abcde..."
	if out["K"] != want {
		t.Errorf("end: want %q got %q", want, out["K"])
	}
}

func TestTruncator_ModeStart(t *testing.T) {
	tr, _ := NewTruncator(8, "start")
	out := tr.Apply(map[string]string{"K": "abcdefghij"})
	want := "...fghij"
	if out["K"] != want {
		t.Errorf("start: want %q got %q", want, out["K"])
	}
}

func TestTruncator_ModeMiddle(t *testing.T) {
	tr, _ := NewTruncator(9, "middle")
	out := tr.Apply(map[string]string{"K": "abcdefghij"})
	// half = (9-3)/2 = 3 => "abc" + "..." + "hij"
	want := "abc...hij"
	if out["K"] != want {
		t.Errorf("middle: want %q got %q", want, out["K"])
	}
}

func TestTruncator_CustomEllipsis(t *testing.T) {
	tr, _ := NewTruncator(6, "end", WithEllipsis("~"))
	out := tr.Apply(map[string]string{"K": "abcdefgh"})
	want := "abcde~"
	if out["K"] != want {
		t.Errorf("custom ellipsis: want %q got %q", want, out["K"])
	}
}

func TestTruncator_Apply_MultipleKeys(t *testing.T) {
	tr, _ := NewTruncator(5, "end")
	in := map[string]string{
		"A": "hi",
		"B": "toolongvalue",
	}
	out := tr.Apply(in)
	if out["A"] != "hi" {
		t.Errorf("A: want %q got %q", "hi", out["A"])
	}
	if out["B"] != "to..." {
		t.Errorf("B: want %q got %q", "to...", out["B"])
	}
}
