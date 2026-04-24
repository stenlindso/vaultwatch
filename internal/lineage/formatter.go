package lineage

import (
	"fmt"
	"io"
	"strings"
)

// FormatText writes a human-readable lineage report for the given records.
func FormatText(w io.Writer, env string, records []Record) {
	fmt.Fprintf(w, "Lineage report — environment: %s\n", env)
	fmt.Fprintln(w, strings.Repeat("-", 50))

	if len(records) == 0 {
		fmt.Fprintln(w, "  (no lineage data)")
		return
	}

	for _, r := range records {
		fmt.Fprintf(w, "  path: %s\n", r.Path)
		for _, e := range r.History {
			fmt.Fprintf(w, "    [%s] source=%s\n",
				e.SeenAt.Format("2006-01-02T15:04:05Z"),
				e.Source,
			)
		}
	}
	fmt.Fprintln(w, strings.Repeat("-", 50))
	fmt.Fprintf(w, "Total paths: %d\n", len(records))
}
