package sync

import (
	"testing"
)

func TestAggregateStrategyFromString_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  AggregateStrategy
	}{
		{"concat", AggregateConcat},
		{"count", AggregateCount},
		{"unique", AggregateUnique},
		{"CONCAT", AggregateConcat},
	}
	for _, tc := range cases {
		got, err := AggregateStrategyFromString(tc.input)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("got %q, want %q", got, tc.want)
		}
	}
}

func TestAggregateStrategyFromString_Invalid(t *testing.T) {
	_, err := AggregateStrategyFromString("average")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestAggregator_Concat_JoinsValues(t *testing.T) {
	a, err := NewAggregator("SHARD_", "SHARDS", ",", AggregateConcat)
	if err != nil {
		t.Fatalf("NewAggregator: %v", err)
	}
	input := map[string]string{
		"SHARD_A": "alpha",
		"SHARD_B": "beta",
		"OTHER":   "x",
	}
	out, err := a.Apply(input)
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if _, ok := out["SHARD_A"]; ok {
		t.Error("source key SHARD_A should be removed")
	}
	got := out["SHARDS"]
	if got != "alpha,beta" {
		t.Errorf("got %q, want %q", got, "alpha,beta")
	}
	if out["OTHER"] != "x" {
		t.Error("non-matching key should be preserved")
	}
}

func TestAggregator_Count_ReturnsLength(t *testing.T) {
	a, _ := NewAggregator("TAG_", "TAG_COUNT", ",", AggregateCount)
	input := map[string]string{
		"TAG_1": "a",
		"TAG_2": "b",
		"TAG_3": "c",
	}
	out, _ := a.Apply(input)
	if out["TAG_COUNT"] != "3" {
		t.Errorf("got %q, want \"3\"", out["TAG_COUNT"])
	}
}

func TestAggregator_Unique_DeduplicatesValues(t *testing.T) {
	a, _ := NewAggregator("ENV_", "ENVS", "|", AggregateUnique)
	input := map[string]string{
		"ENV_1": "prod",
		"ENV_2": "prod",
		"ENV_3": "staging",
	}
	out, _ := a.Apply(input)
	got := out["ENVS"]
	if got != "prod|staging" {
		t.Errorf("got %q, want \"prod|staging\"", got)
	}
}

func TestAggregator_NoMatch_ReturnsUnchanged(t *testing.T) {
	a, _ := NewAggregator("MISSING_", "OUT", ",", AggregateConcat)
	input := map[string]string{"FOO": "bar"}
	out, _ := a.Apply(input)
	if _, ok := out["OUT"]; ok {
		t.Error("output key should not be written when no prefix matches")
	}
}

func TestNewAggregator_EmptyPrefix_ReturnsError(t *testing.T) {
	_, err := NewAggregator("", "OUT", ",", AggregateConcat)
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestNewAggregator_EmptyOutKey_ReturnsError(t *testing.T) {
	_, err := NewAggregator("PRE_", "", ",", AggregateConcat)
	if err == nil {
		t.Fatal("expected error for empty outKey")
	}
}
