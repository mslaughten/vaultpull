// Package sync provides the core synchronisation logic for vaultpull.
//
// # Merger
//
// The Merge function combines a map of values already present in a local .env
// file with a map of values fetched from Vault. The caller controls which set
// wins via a MergeStrategy:
//
//   - MergeStrategyOverwrite   – only Vault values are written (default).
//   - MergeStrategyKeepExisting – local values are preserved; Vault only adds
//     keys that are not already present.
//   - MergeStrategyVaultWins   – both sets are merged, but Vault values take
//     precedence over local ones for keys that exist in both.
//
// The strategy can be parsed from a CLI flag string with MergeStrategyFromString.
package sync
