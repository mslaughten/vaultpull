package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate_MissingVaultAddr(t *testing.T) {
	c := &Config{
		VaultToken: "token",
		MountPath:  "secret",
		KVVersion:  2,
	}
	err := c.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "vault address")
}

func TestValidate_MissingToken(t *testing.T) {
	c := &Config{
		VaultAddr: "http://127.0.0.1:8200",
		MountPath: "secret",
		KVVersion: 2,
	}
	err := c.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "vault token")
}

func TestValidate_InvalidKVVersion(t *testing.T) {
	c := &Config{
		VaultAddr:  "http://127.0.0.1:8200",
		VaultToken: "token",
		MountPath:  "secret",
		KVVersion:  3,
	}
	err := c.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "kv-version")
}

func TestValidate_Valid(t *testing.T) {
	c := &Config{
		VaultAddr:  "http://127.0.0.1:8200",
		VaultToken: "token",
		MountPath:  "secret",
		KVVersion:  2,
	}
	assert.NoError(t, c.Validate())
}

func TestFromEnv_PopulatesFields(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://vault:8200")
	t.Setenv("VAULT_TOKEN", "env-token")
	t.Setenv("VAULT_NAMESPACE", "my-ns")

	c := &Config{}
	c.FromEnv()

	assert.Equal(t, "http://vault:8200", c.VaultAddr)
	assert.Equal(t, "env-token", c.VaultToken)
	assert.Equal(t, "my-ns", c.Namespace)
}

func TestFromEnv_DoesNotOverrideExisting(t *testing.T) {
	os.Setenv("VAULT_ADDR", "http://vault:8200")
	defer os.Unsetenv("VAULT_ADDR")

	c := &Config{VaultAddr: "http://explicit:8200"}
	c.FromEnv()

	assert.Equal(t, "http://explicit:8200", c.VaultAddr)
}
