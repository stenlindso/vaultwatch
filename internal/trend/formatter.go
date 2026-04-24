package trend

import (
	"fmt"
	"io"
	"strings"
)

// FormatText writes a human-readable trend report to w.
func FormatText(results []Result, w io.Writer) {
	if len(results) == 0 {
		fmt.Fprintln(w, "No trend data available.")
		return
	}

	fmt.Fprintln(w, "=== Secret Path Trends ===")
	for _, r := range results {
		arrow := directionArrow(r.Direction)
		fmt.Fprintf(w, "  %-20s %s %-10s  delta=%+d  avg=%+.1f/interval  samples=%d\n",
			r.Env, arrow, r.Direction, r.Delta, r.AvgChange, r.Samples)
	}
}

// FormatOneLiner returns a compact single-line summary of all trends.
func FormatOneLiner(results []Result) string {
	if len(results) == 0 {
		return "trends: no data"
	}
	parts := make([]string, 0, len(results))
	for _, r := range results {
		parts = append(parts, fmt.Sprintf("%s:%s(%+d)", r.Env, r.Direction, r.Delta))
	}
	return "trends: " + strings.Join(parts, ", ")
}

func directionArrow(d Direction) string {
	switch d {
	case DirectionGrowing:
		return "↑"
	case DirectionShrinking:
		return "↓"
	default:
		return "→"
	}
}
