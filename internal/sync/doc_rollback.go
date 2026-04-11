// Package sync provides rollback support for vaultpull sync operations.
//
// Before writing any env file, callers should take a Snapshot of the existing
// file on disk.  If the sync fails partway through, RollbackAll can be used to
// restore every file to its pre-sync state, ensuring that a partial run does
// not leave the working directory in an inconsistent state.
//
// Typical usage:
//
//	snap, err := sync.TakeSnapshot(".env")
//	if err != nil { ... }
//
//	// ... perform writes ...
//
//	if writeErr != nil {
//		sync.RollbackAll([]*sync.Snapshot{snap})
//	}
package sync
