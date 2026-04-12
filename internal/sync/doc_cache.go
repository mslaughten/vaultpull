// Package sync provides the core synchronisation logic for vaultpull.
//
// # Secret Cache
//
// SecretCache stores fetched Vault secrets on disk as JSON files inside a
// configurable directory. Each entry is keyed by a SHA-256 hash of the Vault
// path so that arbitrary path strings map safely to file names.
//
// Entries expire after a caller-supplied TTL. Expired entries are removed
// automatically on the next Get call. Use Invalidate to evict an entry
// immediately (e.g. after a successful push or rotate operation).
//
// The cache directory is created with permission 0700 and individual entry
// files are written with permission 0600 to prevent other users from reading
// cached secrets.
package sync
