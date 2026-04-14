package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewScoper_Valid(t *testing.T) {
	s, err := NewScoper([]string{"DB_:database", "API_*=api", "LOG_LEVEL=logging"})
	require.NoError(t, err)
	assert.NotNil(t, s)
	assert.Len(t, s.rules, 3)
}

func TestNewScoper_MissingSeparator(t *testing.T) {
	_, err := NewScoper([]string{"DB_database"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing '=' separator")
}

func TestNewScoper_EmptyPattern(t *testing.T) {
	_, err := NewScoper([]string{"=myscope"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "empty pattern")
}

func TestNewScoper_EmptyScope(t *testing.T) {
	_, err := NewScoper([]string{"MY_KEY="})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "empty scope")
}

func TestScoper_Apply_ExactMatch(t *testing.T) {
	s, err := NewScoper([]string{"LOG_LEVEL=logging"})
	require.NoError(t, err)

	secrets := map[string]string{
		"LOG_LEVEL": "info",
		"OTHER_KEY": "value",
	}
	result := s.Apply(secrets)
	assert.Contains(t, result["logging"], "LOG_LEVEL")
	assert.Contains(t, result["default"], "OTHER_KEY")
}

func TestScoper_Apply_PrefixMatch(t *testing.T) {
	s, err := NewScoper([]string{"DB_:=database"})
	require.NoError(t, err)

	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV":  "production",
	}
	result := s.Apply(secrets)
	assert.ElementsMatch(t, []string{"DB_HOST", "DB_PORT"}, result["database"])
	assert.Contains(t, result["default"], "APP_ENV")
}

func TestScoper_Apply_GlobMatch(t *testing.T) {
	s, err := NewScoper([]string{"API_*=api"})
	require.NoError(t, err)

	secrets := map[string]string{
		"API_KEY":    "abc",
		"API_SECRET": "xyz",
		"OTHER":      "val",
	}
	result := s.Apply(secrets)
	assert.ElementsMatch(t, []string{"API_KEY", "API_SECRET"}, result["api"])
	assert.Contains(t, result["default"], "OTHER")
}

func TestScoper_Apply_NoRules_AllDefault(t *testing.T) {
	s, err := NewScoper([]string{})
	require.NoError(t, err)

	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	result := s.Apply(secrets)
	assert.Len(t, result["default"], 2)
}

func TestScoper_Apply_EmptySecrets(t *testing.T) {
	s, err := NewScoper([]string{"DB_:=database"})
	require.NoError(t, err)

	result := s.Apply(map[string]string{})
	assert.Empty(t, result)
}
