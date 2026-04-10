// Package sync provides the high-level orchestration layer for vaultpull.
//
// It ties together the vault and envfile packages:
//
//  1. List secret paths from a Vault KV mount.
//  2. Optionally filter paths by a namespace prefix.
//  3. Read each secret's key/value pairs.
//  4. Write the pairs to a corresponding .env file in the output directory.
//
// Typical usage:
//
//	client, _ := vault.NewClient(cfg)
//	s := sync.New(client, sync.Options{
//	    MountPath: "secret",
//	    Namespace: "myapp",
//	    OutputDir: ".",
//	})
//	result, err := s.Run()
package sync
