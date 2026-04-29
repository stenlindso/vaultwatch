package verify

import (
	"fmt"
	"strings"
)

// FormatText renders a slice of Results as a human-readable text report.
func FormatText(results []Result) string {
	if len(results) == 0 {
		return "verify: no results\n"
	}

	var sb strings.Builder
	sb.WriteString("=== Verification Report ===\n")

	for _, r := range results {
		fmt.Fprintf(&sb, "\n[%s] %s\n", r.Environment, r.Timestamp.Format("2006-01-02 15:04:05 UTC"))

		if len(r.Present) > 0 {
			sb.WriteString("  Present:\n")
			for _, p := range r.Present {
				fmt.Fprintf(&sb, "    ✓ %s\n", p)
			}
		}

		if len(r.Missing) > 0 {
			sb.WriteString("  Missing:\n")
			for _, p := range r.Missing {
				fmt.Fprintf(&sb, "    ✗ %s\n", p)
			}
		}

		if !r.HasMissing() {
			sb.WriteString("  All expected paths present.\n")
		}
	}

	sb.WriteString("\n")
	return sb.String()
}

// FormatOneLiner returns a compact single-line summary across all results.
func FormatOneLiner(results []Result) string {
	if len(results) == 0 {
		return "verify: no results"
	}
	totalPresent, totalMissing := 0, 0
	for _, r := range results {
		totalPresent += len(r.Present)
		totalMissing += len(r.Missing)
	}
	status := "OK"
	if totalMissing > 0 {
		status = "FAIL"
	}
	return fmt.Sprintf("verify: envs=%d present=%d missing=%d status=%s",
		len(results), totalPresent, totalMissing, status)
}
