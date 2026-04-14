package sync

import (
	"testing"
)

func TestJoinStrategyFromString_Valid(t *testing.T) {
	cases := []string{"concat", "first", "last"}
	for _, c := range cases {
		_, err := JoinStrategyFromString(c)
		if err != nil {
			t.Errorf("expected no error for %q, got %v", c, err)
		}
	}
}

func TestJoinStrategyFromString_Invalid(t *testing.T) {
	_, err := JoinStrategyFromString("unknown")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestNewJoiner_DefaultSeparator(t *testing.T) {
	j, err := NewJoiner(JoinStrategyConcat, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if j.separator != "," {
		t.Errorf("expected default separator ',', got %q", j.separator)
	}
}

func TestJoiner_Concat_JoinsValues(t *testing.T) {
	j, _ := NewJoiner(JoinStrategyConcat, "|")
	dst := map[string]string{"A": "hello", "B": "only-dst"}
	src := map[string]string{"A": "world", "C": "only-src"}
	out := j.Apply(dst, src)
	if out["A"] != "hello|world" {
		t.Errorf("concat: got %q, want %q", out["A"], "hello|world")
	}
	if out["B"] != "only-dst" {
		t.Errorf("dst-only key: got %q", out["B"])
	}
	if out["C"] != "only-src" {
		t.Errorf("src-only key: got %q", out["C"])
	}
}

func TestJoiner_First_KeepsDst(t *testing.T) {
	j, _ := NewJoiner(JoinStrategyFirstOnly, ",")
	dst := map[string]string{"A": "original"}
	src := map[string]string{"A": "override", "B": "new"}
	out := j.Apply(dst, src)
	if out["A"] != "original" {
		t.Errorf("first: expected original value, got %q", out["A"])
	}
	if out["B"] != "new" {
		t.Errorf("first: missing src-only key, got %q", out["B"])
	}
}

func TestJoiner_Last_SrcWins(t *testing.T) {
	j, _ := NewJoiner(JoinStrategyLastOnly, ",")
	dst := map[string]string{"A": "original", "B": "keep"}
	src := map[string]string{"A": "override"}
	out := j.Apply(dst, src)
	if out["A"] != "override" {
		t.Errorf("last: expected override, got %q", out["A"])
	}
	if out["B"] != "keep" {
		t.Errorf("last: dst-only key should remain, got %q", out["B"])
	}
}

func TestJoiner_EmptyMaps_ReturnsEmpty(t *testing.T) {
	j, _ := NewJoiner(JoinStrategyConcat, ",")
	out := j.Apply(map[string]string{}, map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
