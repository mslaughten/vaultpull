package sync

import (
	"testing"
)

func TestLimitStrategyFromString_Valid(t *testing.T) {
	cases := []string{"first", "last", "alpha"}
	for _, c := range cases {
		_, err := LimitStrategyFromString(c)
		if err != nil {
			t.Errorf("expected no error for %q, got %v", c, err)
		}
	}
}

func TestLimitStrategyFromString_Invalid(t *testing.T) {
	_, err := LimitStrategyFromString("random")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestNewLimiter_InvalidMax(t *testing.T) {
	_, err := NewLimiter(0, LimitStrategyAlpha)
	if err == nil {
		t.Fatal("expected error for max=0")
	}
}

func TestLimiter_Alpha_KeepsFirstNAlphabetically(t *testing.T) {
	l, _ := NewLimiter(2, LimitStrategyAlpha)
	m := map[string]string{"ZEBRA": "z", "APPLE": "a", "MANGO": "m"}
	out := l.Apply(m)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["APPLE"]; !ok {
		t.Error("expected APPLE in output")
	}
	if _, ok := out["MANGO"]; !ok {
		t.Error("expected MANGO in output")
	}
}

func TestLimiter_Last_KeepsLastNAlphabetically(t *testing.T) {
	l, _ := NewLimiter(2, LimitStrategyLast)
	m := map[string]string{"ZEBRA": "z", "APPLE": "a", "MANGO": "m"}
	out := l.Apply(m)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["ZEBRA"]; !ok {
		t.Error("expected ZEBRA in output")
	}
	if _, ok := out["MANGO"]; !ok {
		t.Error("expected MANGO in output")
	}
}

func TestLimiter_MaxGreaterThanMap_ReturnsAll(t *testing.T) {
	l, _ := NewLimiter(100, LimitStrategyAlpha)
	m := map[string]string{"A": "1", "B": "2"}
	out := l.Apply(m)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestLimiter_First_IsDeterministic(t *testing.T) {
	l, _ := NewLimiter(1, LimitStrategyFirst)
	m := map[string]string{"B": "2", "A": "1", "C": "3"}
	out := l.Apply(m)
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
	if _, ok := out["A"]; !ok {
		t.Error("expected A (first alphabetically) in output")
	}
}
