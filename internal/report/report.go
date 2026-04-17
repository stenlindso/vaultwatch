package report

import (
	"time"

	"github.com/vaultwatch/internal/diff"
	"github.com/vaultwatch/internal/policy"
)

// Report aggregates audit diff results and policy violations for a run.
type Report struct {
	GeneratedAt  time.Time         `json:"generated_at"`
	Environments []EnvReport       `json:"environments"`
}

// EnvReport holds the diff and policy results for a single environment pair.
type EnvReport struct {
	BaseEnv    string             `json:"base_env"`
	TargetEnv  string             `json:"target_env"`
	Diff       diff.Result        `json:"diff"`
	Violations []policy.Violation `json:"violations,omitempty"`
}

// Builder constructs a Report incrementally.
type Builder struct {
	report Report
}

// NewBuilder creates a new Builder with the current timestamp.
func NewBuilder() *Builder {
	return &Builder{
		report: Report{
			GeneratedAt: time.Now().UTC(),
		},
	}
}

// AddEnvReport appends an environment report entry.
func (b *Builder) AddEnvReport(e EnvReport) {
	b.report.Environments = append(b.report.Environments, e)
}

// Build returns the assembled Report.
func (b *Builder) Build() Report {
	return b.report
}

// HasAnyViolations returns true if any environment has policy violations.
func (r Report) HasAnyViolations() bool {
	for _, e := range r.Environments {
		if len(e.Violations) > 0 {
			return true
		}
	}
	return false
}

// HasAnyDiffs returns true if any environment has added or removed paths.
func (r Report) HasAnyDiffs() bool {
	for _, e := range r.Environments {
		if len(e.Diff.Added) > 0 || len(e.Diff.Removed) > 0 {
			return true
		}
	}
	return false
}
