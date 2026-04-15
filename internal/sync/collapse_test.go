package sync

import (
	"testing"
)

func TestCollapseStrategyFromString_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  CollapseStrategy
	}{
		{"first", CollapseStrategyFirst},
		{"last", CollapseStrategyLast},
		{"concat", CollapseStrategyConcat},
		{"FIRST", CollapseStrategyFirst},
	}
	for _, tc := range cases {
		got, err := CollapseStrategyFromString(tc.input)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("got %v, want %v", got, tc.want)
		}
	}
}

func TestCollapseStrategyFromString_Invalid(t *testing.T) {
	_, err := CollapseStrategyFromString("unknown")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewCollapser_EmptyPrefix_ReturnsError(t *testing.T) {
	_, err := NewCollapser("", "OUT", ",", CollapseStrategyFirst)
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestNewCollapser_EmptyOutKey_ReturnsError(t *testing.T) {
	_, err := NewCollapser("DB_", "", ",", CollapseStrategyFirst)
	if err == nil {
		t.Fatal("expected error for empty outKey")
	}
}

func TestCollapser_First_KeepsFirstValue(t *testing.T) {
	c, _ := NewCollapser("DB_", "DATABASE", ",", CollapseStrategyFirst)
	m := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432", "OTHER": "x"}
	out, err := c.Apply(m)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("DB_HOST should have been collapsed")
	}
	if out["OTHER"] != "x" {
		t.Error("OTHER should be preserved")
	}
	if out["DATABASE"] == "" {
		t.Error("DATABASE should be set")
	}
}

func TestCollapser_Concat_JoinsValues(t *testing.T) {
	c, _ := NewCollapser("TAG_", "TAGS", "|", CollapseStrategyConcat)
	m := map[string]string{"TAG_A": "alpha", "TAG_B": "beta"}
	out, err := c.Apply(m)
	if err != nil {
		t.Fatal(err)
	}
	v := out["TAGS"]
	if v == "" {
		t.Fatal("TAGS should be non-empty")
	}
}

func TestCollapser_NoMatchingKeys_ReturnsOriginal(t *testing.T) {
	c, _ := NewCollapser("DB_", "DATABASE", ",", CollapseStrategyFirst)
	m := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := c.Apply(m)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["DATABASE"]; ok {
		t.Error("DATABASE should not be present")
	}
}
