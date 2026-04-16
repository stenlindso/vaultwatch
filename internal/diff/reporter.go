package diff

import (
	"fmt"
	"io"
)

// Reporter formats a diff Result for human-readable output.
type Reporter struct {
	w io.Writer
}

// NewReporter creates a Reporter writing to w.
func NewReporter(w io.Writer) *Reporter {
	return &Reporter{w: w}
}

// Print writes a coloured summary of the diff result to the writer.
func (r *Reporter) Print(res Result, labelA, labelB string) {
	fmt.Fprintf(r.w, "\n=== Vault Path Diff: %s vs %s ===\n", labelA, labelB)

	if len(res.OnlyInA) == 0 && len(res.OnlyInB) == 0 {
		fmt.Fprintln(r.w, "✔  No differences found.")
		return
	}

	if len(res.OnlyInA) > 0 {
		fmt.Fprintf(r.w, "\n  Only in %s (%d):\n", labelA, len(res.OnlyInA))
		for _, p := range res.OnlyInA {
			fmt.Fprintf(r.w, "    - %s\n", p)
		}
	}

	if len(res.OnlyInB) > 0 {
		fmt.Fprintf(r.w, "\n  Only in %s (%d):\n", labelB, len(res.OnlyInB))
		for _, p := range res.OnlyInB {
			fmt.Fprintf(r.w, "    + %s\n", p)
		}
	}

	fmt.Fprintf(r.w, "\n  Common paths: %d\n", len(res.InBoth))
}
