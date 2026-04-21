// Package summary provides aggregated statistics across multiple environments
// and snapshot comparisons for high-level reporting.
package summary

import (
	"fmt"
	"sort"
	"time"
)

// EnvStat holds aggregated statistics for a single environment.
type EnvStat struct {
	Environment string
	TotalPaths  int
	Added       int
	Removed     int
	LastScanned time.Time
}

// Report is a cross-environment summary produced by Summarize.
type Report struct {
	GeneratedAt  time.Time
	Environments []EnvStat
	TotalAdded   int
	TotalRemoved int
	TotalPaths   int
}

// HasChanges returns true if any environment has additions or removals.
func (r *Report) HasChanges() bool {
	return r.TotalAdded > 0 || r.TotalRemoved > 0
}

// Input describes per-environment data needed to build a summary.
type Input struct {
	Environment string
	Paths       []string
	Added       []string
	Removed     []string
	LastScanned time.Time
}

// Summarize aggregates a slice of Input values into a Report.
func Summarize(inputs []Input) *Report {
	report := &Report{
		GeneratedAt: time.Now().UTC(),
	}

	for _, in := range inputs {
		stat := EnvStat{
			Environment: in.Environment,
			TotalPaths:  len(in.Paths),
			Added:       len(in.Added),
			Removed:     len(in.Removed),
			LastScanned: in.LastScanned,
		}
		report.Environments = append(report.Environments, stat)
		report.TotalAdded += stat.Added
		report.TotalRemoved += stat.Removed
		report.TotalPaths += stat.TotalPaths
	}

	sort.Slice(report.Environments, func(i, j int) bool {
		return report.Environments[i].Environment < report.Environments[j].Environment
	})

	return report
}

// String returns a compact one-line description of the report.
func (r *Report) String() string {
	return fmt.Sprintf("summary: envs=%d total_paths=%d added=%d removed=%d",
		len(r.Environments), r.TotalPaths, r.TotalAdded, r.TotalRemoved)
}
