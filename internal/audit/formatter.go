package audit

import (
	"fmt"
	"io"
	"strings"
)

// Format writes a human-readable audit report to w.
func Format(w io.Writer, r *Report) {
	fmt.Fprintf(w, "Audit Report\n")
	fmt.Fprintf(w, "  Environments : %s  vs  %s\n", r.EnvironmentA, r.EnvironmentB)
	fmt.Fprintf(w, "  Timestamp    : %s\n", r.Timestamp.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintln(w, strings.Repeat("-", 48))

	printSection(w, fmt.Sprintf("Only in %s (%d)", r.EnvironmentA, len(r.Result.OnlyInA)), r.Result.OnlyInA, "-")
	printSection(w, fmt.Sprintf("Only in %s (%d)", r.EnvironmentB, len(r.Result.OnlyInB)), r.Result.OnlyInB, "+")
	printSection(w, fmt.Sprintf("In both (%d)", len(r.Result.InBoth)), r.Result.InBoth, " ")
}

func printSection(w io.Writer, header string, paths []string, prefix string) {
	fmt.Fprintf(w, "\n%s\n", header)
	if len(paths) == 0 {
		fmt.Fprintln(w, "  (none)")
		return
	}
	for _, p := range paths {
		fmt.Fprintf(w, "  %s %s\n", prefix, p)
	}
}
