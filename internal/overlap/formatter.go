package overlap

import (
	"fmt"
	"sort"
	"strings"
)

// FormatText renders an AnalysisResult as a human-readable text report.
func FormatText(result AnalysisResult) string {
	var sb strings.Builder

	envs := sortedKeys(result.ByEnvironment)

	sb.WriteString(fmt.Sprintf("Overlap Analysis (%d environments)\n", len(envs)))
	sb.WriteString(strings.Repeat("-", 40) + "\n")

	if len(result.SharedByAll) > 0 {
		sb.WriteString(fmt.Sprintf("\nShared across all environments (%d paths):\n", len(result.SharedByAll)))
		shared := sortedSlice(result.SharedByAll)
		for _, p := range shared {
			sb.WriteString(fmt.Sprintf("  %s\n", p))
		}
	} else {
		sb.WriteString("\nNo paths shared across all environments.\n")
	}

	if len(result.SharedPairs) > 0 {
		sb.WriteString(fmt.Sprintf("\nPairwise overlaps (%d pairs):\n", len(result.SharedPairs)))
		pairs := sortedPairKeys(result.SharedPairs)
		for _, key := range pairs {
			paths := result.SharedPairs[key]
			sb.WriteString(fmt.Sprintf("  [%s] — %d shared path(s)\n", key, len(paths)))
			sorted := sortedSlice(paths)
			for _, p := range sorted {
				sb.WriteString(fmt.Sprintf("    %s\n", p))
			}
		}
	}

	sb.WriteString("\nUnique paths per environment:\n")
	for _, env := range envs {
		uniq := result.ByEnvironment[env].UniqueToEnv
		sb.WriteString(fmt.Sprintf("  %s: %d unique\n", env, len(uniq)))
	}

	return sb.String()
}

// FormatOneLiner returns a compact single-line summary of the overlap result.
func FormatOneLiner(result AnalysisResult) string {
	if !result.HasOverlap {
		return fmt.Sprintf("overlap: none across %d environments", len(result.ByEnvironment))
	}
	return fmt.Sprintf(
		"overlap: %d shared-by-all, %d pairs with overlap across %d environments",
		len(result.SharedByAll),
		len(result.SharedPairs),
		len(result.ByEnvironment),
	)
}

func sortedKeys(m map[string]EnvOverlap) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedPairKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedSlice(in []string) []string {
	out := make([]string, len(in))
	copy(out, in)
	sort.Strings(out)
	return out
}
