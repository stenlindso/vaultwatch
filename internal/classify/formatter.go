package classify

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// FormatText writes a human-readable classification report to w.
func FormatText(w io.Writer, grouped map[Level][]Result) {
	levels := []Level{LevelCritical, LevelSecret, LevelInternal, LevelPublic}
	for _, lvl := range levels {
		results, ok := grouped[lvl]
		if !ok || len(results) == 0 {
			continue
		}
		sort.Slice(results, func(i, j int) bool {
			return results[i].Path < results[j].Path
		})
		fmt.Fprintf(w, "[%s] (%d paths)\n", strings.ToUpper(lvl.String()), len(results))
		for _, r := range results {
			if r.Rule != "" {
				fmt.Fprintf(w, "  %s  (matched: %s)\n", r.Path, r.Rule)
			} else {
				fmt.Fprintf(w, "  %s\n", r.Path)
			}
		}
	}
}

// FormatSummary returns a one-line summary of classification counts.
func FormatSummary(grouped map[Level][]Result) string {
	critical := len(grouped[LevelCritical])
	secret := len(grouped[LevelSecret])
	internal := len(grouped[LevelInternal])
	public := len(grouped[LevelPublic])
	return fmt.Sprintf(
		"critical=%d secret=%d internal=%d public=%d",
		critical, secret, internal, public,
	)
}
