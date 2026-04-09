// Package vault provides a thin wrapper around the HashiCorp Vault API client
// tailored for vaultpull's use-cases:
//
//   - Creating an authenticated client from vaultpull configuration.
//   - Reading KV v1 and KV v2 secrets transparently.
//   - Filtering and manipulating secret paths by namespace prefix.
//
// Typical usage:
//
//	client, err := vault.NewClient(cfg.VaultAddr, cfg.Token, cfg.Namespace, cfg.KVVersion)
//	if err != nil { ... }
//
//	paths := vault.FilterByNamespace(allPaths, cfg.Namespace)
//	for _, p := range paths {
//		data, err := client.ReadSecret(p)
//		...
//	}
package vault
