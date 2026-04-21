package baseline

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

// DiffResult holds the output of a baseline diff operation.
type DiffResult struct {
	Environment string
	Label       string
	Added       []string
	Removed     []string
	CheckedAt   time.Time
}

// HasChanges returns true when there are any added or removed paths.
func (d *DiffResult) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0
}

// FormatText writes a human-readable diff summary to w.
func FormatText(w io.Writer, d *DiffResult) {
	fmt.Fprintf(w, "Baseline diff — env: %s  label: %s  checked: %s\n",
		d.Environment, d.Label, d.CheckedAt.Format(time.RFC3339))
	fmt.Fprintln(w, strings.Repeat("-", 60))

	if !d.HasChanges() {
		fmt.Fprintln(w, "No changes from baseline.")
		return
	}

	if len(d.Added) > 0 {
		sort.Strings(d.Added)
		fmt.Fprintln(w, "Added paths:")
		for _, p := range d.Added {
			fmt.Fprintf(w, "  + %s\n", p)
		}
	}

	if len(d.Removed) > 0 {
		sort.Strings(d.Removed)
		fmt.Fprintln(w, "Removed paths:")
		for _, p := range d.Removed {
			fmt.Fprintf(w, "  - %s\n", p)
		}
	}
}
