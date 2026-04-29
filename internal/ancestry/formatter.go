package ancestry

import (
	"fmt"
	"strings"
)

// FormatText renders a human-readable summary of an ancestry Graph.
func FormatText(env string, g Graph) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Ancestry graph — env: %s\n", env))
	sb.WriteString(strings.Repeat("-", 40) + "\n")
	if len(g.Edges) == 0 {
		sb.WriteString("  (no edges recorded)\n")
		return sb.String()
	}
	for _, e := range g.Edges {
		sb.WriteString(fmt.Sprintf("  %s  →  %s\n", e.Parent, e.Child))
	}
	sb.WriteString(fmt.Sprintf("\nTotal edges: %d\n", len(g.Edges)))
	return sb.String()
}

// FormatOrphans renders a list of orphaned paths.
func FormatOrphans(env string, orphans []string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Orphaned paths — env: %s\n", env))
	sb.WriteString(strings.Repeat("-", 40) + "\n")
	if len(orphans) == 0 {
		sb.WriteString("  (none)\n")
		return sb.String()
	}
	for _, p := range orphans {
		sb.WriteString(fmt.Sprintf("  ! %s\n", p))
	}
	sb.WriteString(fmt.Sprintf("\nTotal orphans: %d\n", len(orphans)))
	return sb.String()
}
