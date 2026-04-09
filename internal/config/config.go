package config

import (
	"errors"
	"os"
)

// Config holds the runtime configuration for vaultpull.
type Config struct {
	// VaultAddr is the address of the Vault server.
	VaultAddr string

	// VaultToken is the authentication token.
	VaultToken string

	// Namespace is the Vault namespace (path prefix) to filter secrets from.
	Namespace string

	// MountPath is the KV secrets engine mount path.
	MountPath string

	// OutputFile is the path to the .env file to write secrets into.
	OutputFile string

	// KVVersion is the KV engine version (1 or 2).
	KVVersion int
}

// Validate checks that all required configuration fields are set.
func (c *Config) Validate() error {
	if c.VaultAddr == "" {
		return errors.New("vault address is required (VAULT_ADDR or --vault-addr)")
	}
	if c.VaultToken == "" {
		return errors.New("vault token is required (VAULT_TOKEN or --vault-token)")
	}
	if c.MountPath == "" {
		return errors.New("mount path is required (--mount)")
	}
	if c.KVVersion != 1 && c.KVVersion != 2 {
		return errors.New("kv-version must be 1 or 2")
	}
	return nil
}

// FromEnv populates missing fields from environment variables.
func (c *Config) FromEnv() {
	if c.VaultAddr == "" {
		c.VaultAddr = os.Getenv("VAULT_ADDR")
	}
	if c.VaultToken == "" {
		c.VaultToken = os.Getenv("VAULT_TOKEN")
	}
	if c.Namespace == "" {
		c.Namespace = os.Getenv("VAULT_NAMESPACE")
	}
}
