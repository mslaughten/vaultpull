package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDedupe_NoDuplicates(t *testing.T) {
	maps := []map[string]string{
		{"A": "1", "B": "2"},
		{"C": "3"},
	}
	res, err := Dedupe(maps, DedupeStrategyFirst)
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"A": "1", "B": "2", "C": "3"}, res.Secrets)
	assert.Empty(t, res.Duplicates)
}

func TestDedupe_StrategyFirst_KeepsFirstValue(t *testing.T) {
	maps := []map[string]string{
		{"KEY": "first"},
		{"KEY": "second"},
	}
	res, err := Dedupe(maps, DedupeStrategyFirst)
	require.NoError(t, err)
	assert.Equal(t, "first", res.Secrets["KEY"])
	assert.Equal(t, []string{"KEY"}, res.Duplicates)
}

func TestDedupe_StrategyLast_KeepsLastValue(t *testing.T) {
	maps := []map[string]string{
		{"KEY": "first"},
		{"KEY": "second"},
	}
	res, err := Dedupe(maps, DedupeStrategyLast)
	require.NoError(t, err)
	assert.Equal(t, "second", res.Secrets["KEY"])
	assert.Equal(t, []string{"KEY"}, res.Duplicates)
}

func TestDedupe_StrategyError_ReturnsDuplicateKeyError(t *testing.T) {
	maps := []map[string]string{
		{"KEY": "first"},
		{"KEY": "second"},
	}
	_, err := Dedupe(maps, DedupeStrategyError)
	require.Error(t, err)
	var dke *DuplicateKeyError
	require.ErrorAs(t, err, &dke)
	assert.Equal(t, "KEY", dke.Key)
}

func TestDedupe_DuplicateKeyError_Message(t *testing.T) {
	err := &DuplicateKeyError{Key: "MY_KEY"}
	assert.Equal(t, "duplicate key: MY_KEY", err.Error())
}

func TestDedupe_DuplicatesAreSorted(t *testing.T) {
	maps := []map[string]string{
		{"Z": "1", "A": "1", "M": "1"},
		{"Z": "2", "A": "2", "M": "2"},
	}
	res, err := Dedupe(maps, DedupeStrategyFirst)
	require.NoError(t, err)
	assert.Equal(t, []string{"A", "M", "Z"}, res.Duplicates)
}

func TestDedupeStrategyFromString(t *testing.T) {
	cases := []struct {
		input    string
		want     DedupeStrategy
		ok       bool
	}{
		{"first", DedupeStrategyFirst, true},
		{"last", DedupeStrategyLast, true},
		{"error", DedupeStrategyError, true},
		{"unknown", DedupeStrategyFirst, false},
		{"", DedupeStrategyFirst, false},
	}
	for _, tc := range cases {
		got, ok := DedupeStrategyFromString(tc.input)
		assert.Equal(t, tc.ok, ok, "input=%q", tc.input)
		if tc.ok {
			assert.Equal(t, tc.want, got, "input=%q", tc.input)
		}
	}
}
