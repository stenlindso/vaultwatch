package rollup

import (
	"fmt"
	"io"
	"strings"
)

const barWidth = 20

// FormatText writes a human-readable rollup table to w.
func FormatText(w io.Writer, result Result) {
	if len(result.Entries) == 0 {
		fmt.Fprintln(w, "No rollup data available.")
		return
	}

	fmt.Fprintf(w, "Rollup Summary (%d environment(s))\n", result.TotalEnvs)
	fmt.Fprintln(w, strings.Repeat("-", 60))
	fmt.Fprintf(w, "%-35s %6s  %s\n", "PREFIX", "ENVS", "COVERAGE")
	fmt.Fprintln(w, strings.Repeat("-", 60))

	for _, e := range result.Entries {
		bar := coverageBar(e.Coverage)
		prefix := e.Prefix
		if len(prefix) > 34 {
			prefix = prefix[:31] + "..."
		}
		fmt.Fprintf(w, "%-35s %6d  %s %.0f%%\n", prefix, e.Count, bar, e.Coverage*100)
	}
}

// FormatOneLiner returns a compact single-line summary.
func FormatOneLiner(result Result) string {
	if len(result.Entries) == 0 {
		return "rollup: no data"
	}
	return fmt.Sprintf("rollup: %d prefix group(s) across %d env(s)",
		len(result.Entries), result.TotalEnvs)
}

func coverageBar(coverage float64) string {
	filled := int(coverage * barWidth)
	if filled > barWidth {
		filled = barWidth
	}
	return "[" + strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled) + "]"
}
