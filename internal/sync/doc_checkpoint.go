// Package sync provides the CheckpointStore for persisting incremental sync
// state between vaultpull runs.
//
// A Checkpoint records:
//   - The env file path that was synced
//   - The UTC timestamp of the last successful sync
//   - A map of key → hash strings representing the secret values at sync time
//
// Usage:
//
//	store, err := sync.NewCheckpointStore(".vaultpull/checkpoints")
//	if err != nil { ... }
//
//	// After a successful sync:
//	err = store.Save(".env", hashes)
//
//	// On the next run, load to detect drift:
//	cp, err := store.Load(".env")
//	if cp != nil {
//	    // compare cp.Hashes with current vault values
//	}
package sync
