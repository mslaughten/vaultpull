// Package sync provides the Snapshotter type for capturing and restoring
// named point-in-time snapshots of secret maps.
//
// A snapshot is a labelled JSON file persisted to a configurable directory.
// Snapshots can be used to compare the current state of secrets against a
// previously captured baseline, enabling drift detection and rollback without
// requiring a live Vault connection.
//
// Usage:
//
//	ss, err := sync.NewSnapshotter("/tmp/vaultpull/snapshots")
//	if err != nil { ... }
//
//	// Capture
//	_ = ss.Save("before-deploy", secrets)
//
//	// Restore
//	entry, _ := ss.Load("before-deploy")
package sync
