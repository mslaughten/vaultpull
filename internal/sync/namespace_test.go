package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewNamespaceMapper_Valid(t *testing.T) {
	m, err := NewNamespaceMapper([]string{"prod=PROD", "staging=STG"})
	require.NoError(t, err)
	assert.NotNil(t, m)
}

func TestNewNamespaceMapper_MissingSeparator(t *testing.T) {
	_, err := NewNamespaceMapper([]string{"prodPROD"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid rule")
}

func TestNewNamespaceMapper_EmptyNamespace(t *testing.T) {
	_, err := NewNamespaceMapper([]string{"=PREFIX"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "empty namespace")
}

func TestNamespaceMapper_Apply_RewritesMatchingKeys(t *testing.T) {
	m, err := NewNamespaceMapper([]string{"prod=PROD"})
	require.NoError(t, err)

	src := map[string]string{
		"prod/DB_HOST": "db.example.com",
		"prod/API_KEY": "secret",
		"staging/DB_HOST": "stg.example.com",
	}

	out, err := m.Apply(src)
	require.NoError(t, err)

	assert.Equal(t, "db.example.com", out["PROD_DB_HOST"])
	assert.Equal(t, "secret", out["PROD_API_KEY"])
	// unmatched key is unchanged
	assert.Equal(t, "stg.example.com", out["staging/DB_HOST"])
}

func TestNamespaceMapper_Apply_EmptyPrefix_StripsNamespace(t *testing.T) {
	m, err := NewNamespaceMapper([]string{"prod="})
	require.NoError(t, err)

	out, err := m.Apply(map[string]string{"prod/TOKEN": "abc"})
	require.NoError(t, err)
	assert.Equal(t, "abc", out["TOKEN"])
}

func TestNamespaceMapper_Apply_NoRules_Passthrough(t *testing.T) {
	m, err := NewNamespaceMapper(nil)
	require.NoError(t, err)

	src := map[string]string{"prod/KEY": "val"}
	out, err := m.Apply(src)
	require.NoError(t, err)
	assert.Equal(t, src, out)
}

func TestNamespaceMapper_Apply_MultipleRules(t *testing.T) {
	m, err := NewNamespaceMapper([]string{"prod=PROD", "dev=DEV"})
	require.NoError(t, err)

	src := map[string]string{
		"prod/HOST": "prod-host",
		"dev/HOST":  "dev-host",
	}

	out, err := m.Apply(src)
	require.NoError(t, err)
	assert.Equal(t, "prod-host", out["PROD_HOST"])
	assert.Equal(t, "dev-host", out["DEV_HOST"])
}
