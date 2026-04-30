// Package stale identifies secret paths that have not been observed
// across recent snapshots, flagging them as potentially stale.
package stale

import (
	"fmt"
	"sort"
	"time"
)

// Entry represents a path that may be stale.
type Entry struct {
	Path        string
	Environment string
	LastSeen    time.Time
	AgeDays     int
}

// Result holds the stale detection output for one environment.
type Result struct {
	Environment string
	Entries     []Entry
}

// HasStale returns true if any stale entries were found.
func (r Result) HasStale() bool {
	return len(r.Entries) > 0
}

// Options controls stale detection behaviour.
type Options struct {
	// ThresholdDays is the minimum age in days before a path is considered stale.
	ThresholdDays int
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{ThresholdDays: 30}
}

// Snapshot is the minimal interface the detector needs from a snapshot.
type Snapshot struct {
	Environment string
	Paths       []string
	CapturedAt  time.Time
}

// Detect compares current paths against a reference snapshot taken at an
// earlier point in time and returns paths absent from the current set that
// exceed the age threshold.
func Detect(current, reference Snapshot, opts Options) Result {
	if opts.ThresholdDays <= 0 {
		opts = DefaultOptions()
	}

	currentSet := toSet(current.Paths)
	threshold := time.Duration(opts.ThresholdDays) * 24 * time.Hour
	now := time.Now()

	var entries []Entry
	for _, p := range reference.Paths {
		if _, ok := currentSet[p]; ok {
			continue
		}
		age := now.Sub(reference.CapturedAt)
		if age >= threshold {
			entries = append(entries, Entry{
				Path:        p,
				Environment: reference.Environment,
				LastSeen:    reference.CapturedAt,
				AgeDays:     int(age.Hours() / 24),
			})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})

	return Result{Environment: current.Environment, Entries: entries}
}

// Summary returns a human-readable one-liner for the result.
func Summary(r Result) string {
	if !r.HasStale() {
		return fmt.Sprintf("%s: no stale paths detected", r.Environment)
	}
	return fmt.Sprintf("%s: %d stale path(s) detected", r.Environment, len(r.Entries))
}

func toSet(paths []string) map[string]struct{} {
	s := make(map[string]struct{}, len(paths))
	for _, p := range paths {
		s[p] = struct{}{}
	}
	return s
}
