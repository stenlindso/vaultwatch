package report_test

import (
	"testing"

	"github.com/vaultwatch/internal/diff"
	"github.com/vaultwatch/internal/policy"
	"github.com/vaultwatch/internal/report"
)

func makeEnvReport(added, removed []string, violations []policy.Violation) report.EnvReport {
	return report.EnvReport{
		BaseEnv:   "prod",
		TargetEnv: "staging",
		Diff:      diff.Result{Added: added, Removed: removed},
		Violations: violations,
	}
}

func TestReport_HasAnyDiffs_True(t *testing.T) {
	b := report.NewBuilder()
	b.AddEnvReport(makeEnvReport([]string{"secret/new"}, nil, nil))
	r := b.Build()
	if !r.HasAnyDiffs() {
		t.Error("expected HasAnyDiffs to be true")
	}
}

func TestReport_HasAnyDiffs_False(t *testing.T) {
	b := report.NewBuilder()
	b.AddEnvReport(makeEnvReport(nil, nil, nil))
	r := b.Build()
	if r.HasAnyDiffs() {
		t.Error("expected HasAnyDiffs to be false")
	}
}

func TestReport_HasAnyViolations_True(t *testing.T) {
	v := []policy.Violation{{Path: "secret/bad", Rule: "deny", Message: "denied"}}
	b := report.NewBuilder()
	b.AddEnvReport(makeEnvReport(nil, nil, v))
	r := b.Build()
	if !r.HasAnyViolations() {
		t.Error("expected HasAnyViolations to be true")
	}
}

func TestReport_HasAnyViolations_False(t *testing.T) {
	b := report.NewBuilder()
	b.AddEnvReport(makeEnvReport(nil, nil, nil))
	r := b.Build()
	if r.HasAnyViolations() {
		t.Error("expected HasAnyViolations to be false")
	}
}

func TestBuilder_MultipleEnvs(t *testing.T) {
	b := report.NewBuilder()
	b.AddEnvReport(makeEnvReport([]string{"a"}, nil, nil))
	b.AddEnvReport(makeEnvReport(nil, []string{"b"}, nil))
	r := b.Build()
	if len(r.Environments) != 2 {
		t.Fatalf("expected 2 env reports, got %d", len(r.Environments))
	}
}
