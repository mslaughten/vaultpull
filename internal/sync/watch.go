package sync

import (
	"context"
	"log"
	"time"
)

// Watcher periodically re-runs a Syncer to keep local .env files
// up to date with Vault secrets.
type Watcher struct {
	syncer   *Syncer
	interval time.Duration
	logger   *log.Logger
}

// WatcherOption configures a Watcher.
type WatcherOption func(*Watcher)

// WithWatchLogger sets the logger used by the Watcher.
func WithWatchLogger(l *log.Logger) WatcherOption {
	return func(w *Watcher) {
		w.logger = l
	}
}

// NewWatcher creates a Watcher that triggers syncer every interval.
func NewWatcher(s *Syncer, interval time.Duration, opts ...WatcherOption) *Watcher {
	w := &Watcher{
		syncer:   s,
		interval: interval,
		logger:   log.Default(),
	}
	for _, o := range opts {
		o(w)
	}
	return w
}

// Run starts the watch loop. It blocks until ctx is cancelled.
// The first sync is executed immediately, then on each tick.
func (w *Watcher) Run(ctx context.Context) error {
	if err := w.tick(ctx); err != nil {
		w.logger.Printf("[watch] initial sync error: %v", err)
	}

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Println("[watch] stopping")
			return ctx.Err()
		case <-ticker.C:
			if err := w.tick(ctx); err != nil {
				w.logger.Printf("[watch] sync error: %v", err)
			}
		}
	}
}

func (w *Watcher) tick(ctx context.Context) error {
	w.logger.Println("[watch] running sync")
	result := w.syncer.Run(ctx)
	w.logger.Printf("[watch] sync complete: %s", result.Summary())
	if result.HasErrors() {
		return result.Errors[0]
	}
	return nil
}
