package sync

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewMaskEnvRenderer_Defaults(t *testing.T) {
	r, err := NewMaskEnvRenderer(MaskEnvOptions{}, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.opts.MaskSymbol != "*" {
		t.Errorf("expected default symbol *, got %q", r.opts.MaskSymbol)
	}
	if r.opts.Mode != MaskEnvAll {
		t.Errorf("expected default mode all, got %q", r.opts.Mode)
	}
}

func TestNewMaskEnvRenderer_InvalidMode(t *testing.T) {
	_, err := NewMaskEnvRenderer(MaskEnvOptions{Mode: "middle"}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestNewMaskEnvRenderer_NegativeReveal(t *testing.T) {
	_, err := NewMaskEnvRenderer(MaskEnvOptions{RevealChars: -1}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for negative reveal_chars")
	}
}

func TestMaskEnvRenderer_All_MasksAll(t *testing.T) {
	var buf bytes.Buffer
	r, _ := NewMaskEnvRenderer(MaskEnvOptions{Mode: MaskEnvAll}, &buf)
	err := r.Render(map[string]string{"SECRET": "abc123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "SECRET=******") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestMaskEnvRenderer_Suffix_RevealsPrefix(t *testing.T) {
	var buf bytes.Buffer
	r, _ := NewMaskEnvRenderer(MaskEnvOptions{Mode: MaskEnvSuffix, RevealChars: 2}, &buf)
	err := r.Render(map[string]string{"KEY": "abcdef"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "KEY=ab****") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestMaskEnvRenderer_Prefix_RevealsSuffix(t *testing.T) {
	var buf bytes.Buffer
	r, _ := NewMaskEnvRenderer(MaskEnvOptions{Mode: MaskEnvPrefix, RevealChars: 2}, &buf)
	err := r.Render(map[string]string{"KEY": "abcdef"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "KEY=****ef") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestMaskEnvRenderer_SelectedKeysOnly(t *testing.T) {
	var buf bytes.Buffer
	r, _ := NewMaskEnvRenderer(MaskEnvOptions{Mode: MaskEnvAll, Keys: []string{"SECRET"}}, &buf)
	err := r.Render(map[string]string{"SECRET": "abc", "PLAIN": "xyz"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "SECRET=***") {
		t.Errorf("SECRET should be masked, got: %q", out)
	}
	if !strings.Contains(out, "PLAIN=xyz") {
		t.Errorf("PLAIN should not be masked, got: %q", out)
	}
}

func TestMaskEnvRenderer_SortedOutput(t *testing.T) {
	var buf bytes.Buffer
	r, _ := NewMaskEnvRenderer(MaskEnvOptions{}, &buf)
	err := r.Render(map[string]string{"Z_KEY": "1", "A_KEY": "2", "M_KEY": "3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 || !strings.HasPrefix(lines[0], "A_KEY") {
		t.Errorf("expected sorted output, got: %v", lines)
	}
}
