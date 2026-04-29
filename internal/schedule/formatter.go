package schedule

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

// FormatText writes a human-readable schedule table to w.
func FormatText(w io.Writer, entries []*Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "No scheduled environments.")
		return
	}

	sorted := make([]*Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Environment < sorted[j].Environment
	})

	fmt.Fprintln(w, "Scheduled Environments")
	fmt.Fprintln(w, strings.Repeat("-", 52))
	fmt.Fprintf(w, "%-20s %-12s %-18s\n", "Environment", "Interval", "Next Run")
	fmt.Fprintln(w, strings.Repeat("-", 52))

	for _, e := range sorted {
		nextRun := e.NextRun.UTC().Format(time.RFC3339)
		fmt.Fprintf(w, "%-20s %-12s %-18s\n", e.Environment, e.Interval.String(), nextRun)
	}
}

// FormatOneLiner returns a compact single-line summary.
func FormatOneLiner(entries []*Entry) string {
	if len(entries) == 0 {
		return "schedule: 0 environments registered"
	}
	envs := make([]string, 0, len(entries))
	for _, e := range entries {
		envs = append(envs, e.Environment)
	}
	sort.Strings(envs)
	return fmt.Sprintf("schedule: %d environment(s) — %s", len(entries), strings.Join(envs, ", "))
}
