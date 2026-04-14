package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCoalesceStrategyFromString_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  CoalesceStrategy
	}{
		{"first", CoalesceFirst},
		{"", CoalesceFirst},
		{"longest", CoalesceLongest},
		{"nonzero", CoalesceNonZero},
	}
	for _, tc := range cases {
		s, err := CoalesceStrategyFromString(tc.input)
		require.NoError(t, err)
		assert.Equal(t, tc.want, s)
	}
}

func TestCoalesceStrategyFromString_Invalid(t *testing.T) {
	_, err := CoalesceStrategyFromString("bogus")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "bogus")
}

func TestCoalescer_Apply_FirstWins(t *testing.T) {
	c, err := NewCoalescer("first")
	require.NoError(t, err)

	a := map[string]string{"KEY": "alpha", "ONLY_A": "yes"}
	b := map[string]string{"KEY": "beta", "ONLY_B": "yes"}

	got := c.Apply(a, b)
	assert.Equal(t, "alpha", got["KEY"])
	assert.Equal(t, "yes", got["ONLY_A"])
	assert.Equal(t, "yes", got["ONLY_B"])
}

func TestCoalescer_Apply_Longest(t *testing.T) {
	c, err := NewCoalescer("longest")
	require.NoError(t, err)

	a := map[string]string{"KEY": "hi"}
	b := map[string]string{"KEY": "hello-world"}

	got := c.Apply(a, b)
	assert.Equal(t, "hello-world", got["KEY"])
}

func TestCoalescer_Apply_NonZero_PrefersNonEmpty(t *testing.T) {
	c, err := NewCoalescer("nonzero")
	require.NoError(t, err)

	a := map[string]string{"KEY": ""}
	b := map[string]string{"KEY": "filled"}

	got := c.Apply(a, b)
	assert.Equal(t, "filled", got["KEY"])
}

func TestCoalescer_Apply_NonZero_KeepsExistingNonEmpty(t *testing.T) {
	c, err := NewCoalescer("nonzero")
	require.NoError(t, err)

	a := map[string]string{"KEY": "original"}
	b := map[string]string{"KEY": "replacement"}

	got := c.Apply(a, b)
	assert.Equal(t, "original", got["KEY"])
}

func TestNewCoalescer_InvalidStrategy_ReturnsError(t *testing.T) {
	_, err := NewCoalescer("unknown")
	require.Error(t, err)
}
