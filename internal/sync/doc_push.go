// Package sync provides the core synchronisation logic for vaultpull.
//
// # Push
//
// The push sub-feature is the inverse of the standard pull workflow: it reads
// a local .env file and writes every key/value pair back into HashiCorp Vault
// at the specified secret path.
//
// Usage:
//
//	pusher := sync.NewPusher(vaultClient)
//	result := pusher.Push(ctx, ".env", "secret/myapp/production")
//	if result.Err != nil {
//		log.Fatal(result.Err)
//	}
//	fmt.Printf("pushed %d keys to %s\n", result.Written, result.VaultPath)
//
// Both KV v1 and KV v2 engines are supported; the underlying vault.Client
// handles payload wrapping automatically.
package sync
