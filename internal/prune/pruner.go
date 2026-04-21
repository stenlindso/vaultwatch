// Package prune provides functionality for identifying and removing stale
// snapshot entries based on age or count thresholds.
package prune

import (
	"fmt"
	"sort"
	"time"
)

// Entry represents a snapshot entry eligible for pruning evaluation.
type Entry struct {
	Environment string
	Label       string
	CreatedAt   time.Time
}

// Options configures pruning behaviour.
type Options struct {
	// MaxAge removes entries older than this duration. Zero means no age limit.
	MaxAge time.Duration
	// MaxCount keeps only the N most recent entries per environment. Zero means no limit.
	MaxCount int
}

// Result holds the outcome of a prune operation.
type Result struct {
	Pruned  []Entry
	Retained []Entry
}

// HasPruned reports whether any entries were marked for removal.
func (r Result) HasPruned() bool {
	return len(r.Pruned) > 0
}

// Summary returns a human-readable summary of the prune result.
func (r Result) Summary() string {
	return fmt.Sprintf("pruned %d entries, retained %d entries", len(r.Pruned), len(r.Retained))
}

// Evaluate applies the given Options to entries and returns a Result indicating
// which entries should be pruned and which should be retained.
func Evaluate(entries []Entry, opts Options, now time.Time) Result {
	if len(entries) == 0 {
		return Result{}
	}

	pruneSet := make(map[int]bool)

	// Age-based pruning.
	if opts.MaxAge > 0 {
		cutoff := now.Add(-opts.MaxAge)
		for i, e := range entries {
			if e.CreatedAt.Before(cutoff) {
				pruneSet[i] = true
			}
		}
	}

	// Count-based pruning per environment.
	if opts.MaxCount > 0 {
		byEnv := make(map[string][]int)
		for i, e := range entries {
			byEnv[e.Environment] = append(byEnv[e.Environment], i)
		}
		for _, indices := range byEnv {
			sort.Slice(indices, func(a, b int) bool {
				return entries[indices[a]].CreatedAt.After(entries[indices[b]].CreatedAt)
			})
			if len(indices) > opts.MaxCount {
				for _, idx := range indices[opts.MaxCount:] {
					pruneSet[idx] = true
				}
			}
		}
	}

	var result Result
	for i, e := range entries {
		if pruneSet[i] {
			result.Pruned = append(result.Pruned, e)
		} else {
			result.Retained = append(result.Retained, e)
		}
	}
	return result
}
