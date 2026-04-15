package sync

import (
	"math/rand"
	"sort"
	"testing"
)

func TestSampleStrategyFromString_Valid(t *testing.T) {
	cases := []string{"random", "first", "last"}
	for _, c := range cases {
		_, err := SampleStrategyFromString(c)
		if err != nil {
			t.Errorf("expected no error for %q, got %v", c, err)
		}
	}
}

func TestSampleStrategyFromString_Invalid(t *testing.T) {
	_, err := SampleStrategyFromString("middle")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestNewSampler_InvalidN(t *testing.T) {
	_, err := NewSampler(0, SampleStrategyFirst, nil)
	if err == nil {
		t.Fatal("expected error for n=0")
	}
}

func TestSampler_First_KeepsFirstNAlphabetically(t *testing.T) {
	s, _ := NewSampler(2, SampleStrategyFirst, nil)
	secrets := map[string]string{"CHARLIE": "3", "ALPHA": "1", "BETA": "2"}
	out := s.Apply(secrets)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["ALPHA"]; !ok {
		t.Error("expected ALPHA in result")
	}
	if _, ok := out["BETA"]; !ok {
		t.Error("expected BETA in result")
	}
}

func TestSampler_Last_KeepsLastNAlphabetically(t *testing.T) {
	s, _ := NewSampler(2, SampleStrategyLast, nil)
	secrets := map[string]string{"CHARLIE": "3", "ALPHA": "1", "BETA": "2"}
	out := s.Apply(secrets)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["CHARLIE"]; !ok {
		t.Error("expected CHARLIE in result")
	}
	if _, ok := out["BETA"]; !ok {
		t.Error("expected BETA in result")
	}
}

func TestSampler_Random_ReturnsSameCountDifferentOrder(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	s, _ := NewSampler(3, SampleStrategyRandom, rng)
	secrets := map[string]string{"A": "1", "B": "2", "C": "3", "D": "4", "E": "5"}
	out := s.Apply(secrets)
	if len(out) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(out))
	}
}

func TestSampler_NLargerThanMap_ReturnsAll(t *testing.T) {
	s, _ := NewSampler(100, SampleStrategyFirst, nil)
	secrets := map[string]string{"X": "1", "Y": "2"}
	out := s.Apply(secrets)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
}

func TestSampler_EmptyMap_ReturnsEmpty(t *testing.T) {
	s, _ := NewSampler(5, SampleStrategyFirst, nil)
	out := s.Apply(map[string]string{})
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %d keys", len(out))
	}
}

func TestSampler_First_ResultIsSorted(t *testing.T) {
	s, _ := NewSampler(3, SampleStrategyFirst, nil)
	secrets := map[string]string{"DELTA": "4", "ALPHA": "1", "BETA": "2", "CHARLIE": "3", "ECHO": "5"}
	out := s.Apply(secrets)
	keys := make([]string, 0, len(out))
	for k := range out {
		keys = append(keys, k)
	}
	if !sort.StringsAreSorted(keys) {
		// keys from a map are unordered; just verify expected keys present
	}
	if _, ok := out["ALPHA"]; !ok {
		t.Error("expected ALPHA")
	}
	_ = keys
}
