package watch

import (
	"context"
	"log"
	"time"
)

// Scheduler runs a Watcher on a fixed interval for a set of environments.
type Scheduler struct {
	watcher  *Watcher
	notifier *Notifier
	envs     []string
	interval time.Duration
}

// NewScheduler creates a Scheduler that polls the given environments.
func NewScheduler(w *Watcher, n *Notifier, envs []string, interval time.Duration) *Scheduler {
	return &Scheduler{
		watcher:  w,
		notifier: n,
		envs:     envs,
		interval: interval,
	}
}

// Run starts the polling loop and blocks until ctx is cancelled.
func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	log.Printf("scheduler: starting, interval=%s envs=%v", s.interval, s.envs)

	for {
		select {
		case <-ctx.Done():
			log.Println("scheduler: stopped")
			return
		case <-ticker.C:
			s.poll(ctx)
		}
	}
}

func (s *Scheduler) poll(ctx context.Context) {
	for _, env := range s.envs {
		events, err := s.watcher.Diff(ctx, env)
		if err != nil {
			log.Printf("scheduler: diff error env=%s: %v", env, err)
			continue
		}
		for _, ev := range events {
			if err := s.notifier.Handle(ev); err != nil {
				log.Printf("scheduler: notify error env=%s: %v", env, err)
			}
		}
	}
}
