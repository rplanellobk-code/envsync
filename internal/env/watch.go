package env

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"
)

// WatchEvent describes a change detected on a remote environment.
type WatchEvent struct {
	Environment string
	PreviousSum [32]byte
	CurrentSum  [32]byte
	ChangedAt   time.Time
}

// Watcher polls a Vault at a fixed interval and emits WatchEvents when the
// contents of a tracked environment change.
type Watcher struct {
	vault      *Vault
	passphrase string
	interval   time.Duration
	sums       map[string][32]byte
}

// NewWatcher creates a Watcher for the given Vault.
// interval controls how often the remote is polled.
func NewWatcher(v *Vault, passphrase string, interval time.Duration) *Watcher {
	return &Watcher{
		vault:      v,
		passphrase: passphrase,
		interval:   interval,
		sums:       make(map[string][32]byte),
	}
}

// Watch starts polling for changes to env in the background.
// Events are sent on the returned channel; the channel is closed when ctx is
// cancelled or an unrecoverable error occurs.
func (w *Watcher) Watch(ctx context.Context, env string) (<-chan WatchEvent, error) {
	// Seed the initial checksum so the first tick does not fire a false positive.
	sum, err := w.checksum(env)
	if err != nil {
		return nil, fmt.Errorf("watch: initial pull failed: %w", err)
	}
	w.sums[env] = sum

	ch := make(chan WatchEvent, 4)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				newSum, err := w.checksum(env)
				if err != nil {
					// Transient error — keep polling.
					continue
				}
				prev := w.sums[env]
				if newSum != prev {
					w.sums[env] = newSum
					select {
					case ch <- WatchEvent{
						Environment: env,
						PreviousSum: prev,
						CurrentSum:  newSum,
						ChangedAt:   t,
					}:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()
	return ch, nil
}

func (w *Watcher) checksum(env string) ([32]byte, error) {
	vals, err := w.vault.Pull(env, w.passphrase)
	if err != nil {
		return [32]byte{}, err
	}
	return sha256.Sum256([]byte(Serialize(vals))), nil
}
