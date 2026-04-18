package watch

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestScheduler_RunAndStop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Watcher with no snapshot manager — Diff will return empty, no panic.
	w := &Watcher{snapshots: nil, lister: nil}

	var mu sync.Mutex
	var out strings.Builder

	n := &Notifier{
		handlers: []EventHandler{
			func(ev Event) error {
				mu.Lock()
				out.WriteString(ev.Env)
				mu.Unlock()
				return nil
			},
		},
	}

	_ = w
	_ = n

	// Scheduler with a very short interval; cancel quickly.
	s := NewScheduler(w, n, []string{"dev"}, 10*time.Millisecond)

	done := make(chan struct{})
	go func() {
		s.Run(ctx)
		close(done)
	}()

	time.Sleep(30 * time.Millisecond)
	cancel()

	select {
	case <-done:
		// ok
	case <-time.After(500 * time.Millisecond):
		t.Fatal("scheduler did not stop after context cancellation")
	}
}

func TestNewScheduler_Fields(t *testing.T) {
	w := &Watcher{}
	n := &Notifier{}
	envs := []string{"prod", "staging"}
	interval := 5 * time.Minute

	s := NewScheduler(w, n, envs, interval)

	if s.interval != interval {
		t.Errorf("expected interval %v, got %v", interval, s.interval)
	}
	if len(s.envs) != 2 {
		t.Errorf("expected 2 envs, got %d", len(s.envs))
	}
}
