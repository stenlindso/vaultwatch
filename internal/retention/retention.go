// Package retention provides policies for managing snapshot and history lifecycle.
package retention

import (
	"fmt"
	"time"
)

// Policy defines the rules for retaining or discarding entries.
type Policy struct {
	MaxAge   time.Duration
	MaxCount int
}

// Entry represents a retainable item with a timestamp and identifier.
type Entry struct {
	ID        string
	Env       string
	CreatedAt time.Time
}

// Result holds the outcome of applying a retention policy.
type Result struct {
	Retained []Entry
	Pruned   []Entry
}

// HasPruned returns true if any entries were pruned.
func (r Result) HasPruned() bool {
	return len(r.Pruned) > 0
}

// Summary returns a human-readable summary of the result.
func (r Result) Summary() string {
	return fmt.Sprintf("retained=%d pruned=%d", len(r.Retained), len(r.Pruned))
}

// Apply evaluates the given entries against the policy and returns a Result.
// Entries should be provided in descending order (newest first).
func Apply(p Policy, entries []Entry) Result {
	var retained, pruned []Entry
	now := time.Now()

	for i, e := range entries {
		prunedByAge := p.MaxAge > 0 && now.Sub(e.CreatedAt) > p.MaxAge
		prunedByCount := p.MaxCount > 0 && i >= p.MaxCount

		if prunedByAge || prunedByCount {
			pruned = append(pruned, e)
		} else {
			retained = append(retained, e)
		}
	}

	return Result{Retained: retained, Pruned: pruned}
}

// DefaultPolicy returns a sensible default retention policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxAge:   30 * 24 * time.Hour,
		MaxCount: 50,
	}
}
