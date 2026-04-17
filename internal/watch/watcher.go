package watch

import (
	"context"
	"log"
	"time"

	"github.com/vaultwatch/internal/vault"
)

// Watcher polls Vault paths at a fixed interval and emits change events.
type Watcher struct {
	lister   *vault.Lister
	interval time.Duration
	env      string
}

// Event represents a detected change in Vault paths.
type Event struct {
	Env     string
	Added   []string
	Removed []string
	At      time.Time
}

// NewWatcher creates a Watcher for the given environment.
func NewWatcher(lister *vault.Lister, env string, interval time.Duration) *Watcher {
	return &Watcher{lister: lister, env: env, interval: interval}
}

// Watch polls Vault and sends Events on the returned channel until ctx is cancelled.
func (w *Watcher) Watch(ctx context.Context) <-chan Event {
	ch := make(chan Event, 4)
	go func() {
		defer close(ch)
		prev, err := w.lister.ListPaths(ctx, "secret/")
		if err != nil {
			log.Printf("watch: initial list failed: %v", err)
			return
		}
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				curr, err := w.lister.ListPaths(ctx, "secret/")
				if err != nil {
					log.Printf("watch: list failed: %v", err)
					continue
				}
				ev := diff(w.env, prev, curr)
				if len(ev.Added) > 0 || len(ev.Removed) > 0 {
					ch <- ev
				}
				prev = curr
			}
		}
	}()
	return ch
}

 prev, curr []string) Event {
	prevSet := toSet(prev)
	currSet := toSet(curr)
	ev := Event{Env: env, At: time.Now()}
	for p := range currSet {
		if !prevSet[p] {
			ev.Added = append(ev.Added, p)
		}
	}
	for p := range prevSet {
		if !currSet[p] {
			ev.Removed = append(ev.Removed, p)
		}
	}
	return ev
}

func toSet(paths []string) map[string]bool {
	s := make(map[string]bool, len(paths))
	for _, p := range paths {
		s[p] = true
	}
	return s
}
