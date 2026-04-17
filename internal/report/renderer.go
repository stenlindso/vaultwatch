package report

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Format controls the output format of a rendered report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Render writes the report to w in the specified format.
func Render(w io.Writer, r Report, f Format) error {
	switch f {
	case FormatJSON:
		return renderJSON(w, r)
	default:
		return renderText(w, r)
	}
}

func renderJSON(w io.Writer, r Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func renderText(w io.Writer, r Report) error {
	fmt.Fprintf(w, "VaultWatch Report — %s\n", r.GeneratedAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintln(w, strings.Repeat("=", 50))
	for _, e := range r.Environments {
		fmt.Fprintf(w, "\n[%s → %s]\n", e.BaseEnv, e.TargetEnv)
		if len(e.Diff.Added) > 0 {
			fmt.Fprintln(w, "  Added:")
			for _, p := range e.Diff.Added {
				fmt.Fprintf(w, "    + %s\n", p)
			}
		}
		if len(e.Diff.Removed) > 0 {
			fmt.Fprintln(w, "  Removed:")
			for _, p := range e.Diff.Removed {
				fmt.Fprintf(w, "    - %s\n", p)
			}
		}
		if len(e.Violations) > 0 {
			fmt.Fprintln(w, "  Policy Violations:")
			for _, v := range e.Violations {
				fmt.Fprintf(w, "    ! [%s] %s — %s\n", v.Rule, v.Path, v.Message)
			}
		}
		if len(e.Diff.Added) == 0 && len(e.Diff.Removed) == 0 && len(e.Violations) == 0 {
			fmt.Fprintln(w, "  No changes or violations.")
		}
	}
	return nil
}
