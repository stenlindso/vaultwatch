package summary

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"
)

// FormatText writes a human-readable summary table to w.
func FormatText(w io.Writer, r *Report) {
	fmt.Fprintf(w, "Generated: %s\n", r.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(w, "Total paths: %d | Added: %d | Removed: %d\n\n",
		r.TotalPaths, r.TotalAdded, r.TotalRemoved)

	if len(r.Environments) == 0 {
		fmt.Fprintln(w, "No environments found.")
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "ENVIRONMENT\tPATHS\tADDED\tREMOVED\tLAST SCANNED")
	fmt.Fprintln(tw, strings.Repeat("-", 60))

	for _, e := range r.Environments {
		fmt.Fprintf(tw, "%s\t%d\t%d\t%d\t%s\n",
			e.Environment,
			e.TotalPaths,
			e.Added,
			e.Removed,
			e.LastScanned.Format(time.RFC3339),
		)
	}
	tw.Flush()
}

// FormatOneLiner returns a compact single-line summary string.
func FormatOneLiner(r *Report) string {
	if !r.HasChanges() {
		return fmt.Sprintf("[OK] %d environments, %d paths — no changes detected",
			len(r.Environments), r.TotalPaths)
	}
	return fmt.Sprintf("[CHANGED] %d environments, %d paths — +%d/-%d",
		len(r.Environments), r.TotalPaths, r.TotalAdded, r.TotalRemoved)
}
