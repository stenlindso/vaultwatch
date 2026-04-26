// Package recommend provides suggestions for improving Vault secret path
// organization based on observed patterns, policy violations, drift, and scores.
package recommend

import (
	"fmt"
	"sort"
	"strings"
)

// Priority indicates how urgently a recommendation should be addressed.
type Priority string

const (
	PriorityCritical Priority = "critical"
	PriorityHigh     Priority = "high"
	PriorityMedium   Priority = "medium"
	PriorityLow      Priority = "low"
)

// Category groups recommendations by concern area.
type Category string

const (
	CategoryDrift    Category = "drift"
	CategoryPolicy   Category = "policy"
	CategoryScore    Category = "score"
	CategoryRetention Category = "retention"
	CategoryStructure Category = "structure"
)

// Recommendation represents a single actionable suggestion.
type Recommendation struct {
	Priority    Priority
	Category    Category
	Environment string
	Message     string
	Detail      string
}

// Input holds the data used to generate recommendations.
type Input struct {
	// DriftedPaths maps environment name to paths that have drifted.
	DriftedPaths map[string][]string
	// PolicyViolations maps environment name to violated rule names.
	PolicyViolations map[string][]string
	// Scores maps environment name to its computed health score (0–100).
	Scores map[string]int
	// StaleDays maps environment name to days since last snapshot.
	StaleDays map[string]int
	// PathCounts maps environment name to total path count.
	PathCounts map[string]int
}

// Result holds all generated recommendations.
type Result struct {
	Items []Recommendation
}

// HasAny returns true if there is at least one recommendation.
func (r *Result) HasAny() bool {
	return len(r.Items) > 0
}

// ByPriority returns recommendations filtered to a specific priority.
func (r *Result) ByPriority(p Priority) []Recommendation {
	var out []Recommendation
	for _, item := range r.Items {
		if item.Priority == p {
			out = append(out, item)
		}
	}
	return out
}

// Generate produces recommendations based on the provided input data.
func Generate(in Input) *Result {
	var items []Recommendation

	// Drift recommendations
	for env, paths := range in.DriftedPaths {
		if len(paths) == 0 {
			continue
		}
		priority := PriorityMedium
		if len(paths) > 20 {
			priority = PriorityHigh
		}
		items = append(items, Recommendation{
			Priority:    priority,
			Category:    CategoryDrift,
			Environment: env,
			Message:     fmt.Sprintf("%d drifted path(s) detected", len(paths)),
			Detail:      fmt.Sprintf("Review: %s", strings.Join(truncate(paths, 3), ", ")),
		})
	}

	// Policy violation recommendations
	for env, violations := range in.PolicyViolations {
		if len(violations) == 0 {
			continue
		}
		items = append(items, Recommendation{
			Priority:    PriorityCritical,
			Category:    CategoryPolicy,
			Environment: env,
			Message:     fmt.Sprintf("%d policy violation(s) found", len(violations)),
			Detail:      fmt.Sprintf("Violated rules: %s", strings.Join(truncate(violations, 3), ", ")),
		})
	}

	// Score-based recommendations
	for env, score := range in.Scores {
		switch {
		case score < 40:
			items = append(items, Recommendation{
				Priority:    PriorityHigh,
				Category:    CategoryScore,
				Environment: env,
				Message:     fmt.Sprintf("Low health score: %d/100", score),
				Detail:      "Investigate drift, policy violations, and snapshot freshness.",
			})
		case score < 70:
			items = append(items, Recommendation{
				Priority:    PriorityMedium,
				Category:    CategoryScore,
				Environment: env,
				Message:     fmt.Sprintf("Moderate health score: %d/100", score),
				Detail:      "Consider reviewing recent changes and policy compliance.",
			})
		}
	}

	// Stale snapshot recommendations
	for env, days := range in.StaleDays {
		if days >= 7 {
			priority := PriorityMedium
			if days >= 30 {
				priority = PriorityHigh
			}
			items = append(items, Recommendation{
				Priority:    priority,
				Category:    CategoryRetention,
				Environment: env,
				Message:     fmt.Sprintf("Snapshot is %d day(s) old", days),
				Detail:      "Run a fresh snapshot to ensure audit accuracy.",
			})
		}
	}

	// Structure recommendations based on path count anomalies
	for env, count := range in.PathCounts {
		if count == 0 {
			items = append(items, Recommendation{
				Priority:    PriorityLow,
				Category:    CategoryStructure,
				Environment: env,
				Message:     "No paths found in environment",
				Detail:      "Verify that the environment is correctly configured and accessible.",
			})
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return priorityOrder(items[i].Priority) < priorityOrder(items[j].Priority)
	})

	return &Result{Items: items}
}

func priorityOrder(p Priority) int {
	switch p {
	case PriorityCritical:
		return 0
	case PriorityHigh:
		return 1
	case PriorityMedium:
		return 2
	default:
		return 3
	}
}

// truncate returns up to n items from s, appending "..." if truncated.
func truncate(s []string, n int) []string {
	if len(s) <= n {
		return s
	}
	return append(s[:n], fmt.Sprintf("(+%d more)", len(s)-n))
}
