// Package schedule provides cron-style scheduling for periodic Vault path scans.
package schedule

import (
	"errors"
	"fmt"
	"time"
)

// Entry represents a single scheduled scan job.
type Entry struct {
	Environment string
	Interval    time.Duration
	LastRun     time.Time
	NextRun     time.Time
}

// Registry holds all scheduled entries keyed by environment name.
type Registry struct {
	entries map[string]*Entry
	now     func() time.Time
}

// New creates an empty Registry.
func New() *Registry {
	return &Registry{
		entries: make(map[string]*Entry),
		now:     time.Now,
	}
}

// Register adds or replaces a scheduled entry for the given environment.
func (r *Registry) Register(env string, interval time.Duration) error {
	if env == "" {
		return errors.New("schedule: environment name must not be empty")
	}
	if interval <= 0 {
		return fmt.Errorf("schedule: interval must be positive, got %s", interval)
	}
	now := r.now()
	r.entries[env] = &Entry{
		Environment: env,
		Interval:    interval,
		NextRun:     now.Add(interval),
	}
	return nil
}

// Due returns all entries whose NextRun is at or before now.
func (r *Registry) Due() []*Entry {
	now := r.now()
	var due []*Entry
	for _, e := range r.entries {
		if !e.NextRun.After(now) {
			due = append(due, e)
		}
	}
	return due
}

// Advance marks an entry as having run, updating LastRun and NextRun.
func (r *Registry) Advance(env string) error {
	e, ok := r.entries[env]
	if !ok {
		return fmt.Errorf("schedule: no entry for environment %q", env)
	}
	now := r.now()
	e.LastRun = now
	e.NextRun = now.Add(e.Interval)
	return nil
}

// List returns all registered entries.
func (r *Registry) List() []*Entry {
	out := make([]*Entry, 0, len(r.entries))
	for _, e := range r.entries {
		out = append(out, e)
	}
	return out
}

// Remove deletes the entry for the given environment.
func (r *Registry) Remove(env string) {
	delete(r.entries, env)
}
