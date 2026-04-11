// Package sync provides the core synchronisation logic for vaultpull.
//
// # Watch
//
// The Watcher type continuously polls HashiCorp Vault and re-writes local
// .env files whenever the configured interval elapses.
//
// Basic usage:
//
//	syncer := sync.New(client, cfg)
//	watcher := sync.NewWatcher(syncer, 30*time.Second)
//	if err := watcher.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
//		log.Fatal(err)
//	}
//
// The first sync is executed immediately on Run; subsequent syncs fire on
// each tick. Errors from individual sync cycles are logged but do not stop
// the watch loop — only context cancellation terminates it.
package sync
