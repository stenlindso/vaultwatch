package quota

import (
	"strings"
	"testing"
)

func makePaths(n int) []string {
	paths := make([]string, n)
	for i := range paths {
		paths[i] = fmt.Sprintf("secret/path/%d", i)
	}
	return paths
}

func TestEvaluate_WithinLimits(t *testing.T) {
	envPaths := map[string][]string{
		"prod": {"secret/a", "secret/b"},
		"dev":  {"secret/x"},
	}
	report := Evaluate(envPaths, Policy{MaxPaths: 10, WarnAt: 5})
	if report.HasExceeded() {
		t.Error("expected no exceeded environments")
	}
	if report.HasWarnings() {
		t.Error("expected no warnings")
	}
	if len(report.Statuses) != 2 {
		t.Fatalf("expected 2 statuses, got %d", len(report.Statuses))
	}
}

func TestEvaluate_Exceeded(t *testing.T) {
	envPaths := map[string][]string{
		"prod": {"a", "b", "c", "d", "e"},
	}
	report := Evaluate(envPaths, Policy{MaxPaths: 3, WarnAt: 2})
	if !report.HasExceeded() {
		t.Error("expected exceeded to be true")
	}
	if report.Statuses[0].PathCount != 5 {
		t.Errorf("expected PathCount 5, got %d", report.Statuses[0].PathCount)
	}
}

func TestEvaluate_Warning(t *testing.T) {
	envPaths := map[string][]string{
		"staging": {"a", "b", "c"},
	}
	report := Evaluate(envPaths, Policy{MaxPaths: 10, WarnAt: 3})
	if report.HasExceeded() {
		t.Error("expected no exceeded")
	}
	if !report.HasWarnings() {
		t.Error("expected warning to be triggered")
	}
	if !report.Statuses[0].Warning {
		t.Error("expected staging status to have Warning=true")
	}
}

func TestEvaluate_NoLimit(t *testing.T) {
	envPaths := map[string][]string{
		"prod": {"a", "b", "c", "d", "e", "f", "g", "h"},
	}
	report := Evaluate(envPaths, Policy{MaxPaths: 0, WarnAt: 0})
	if report.HasExceeded() {
		t.Error("expected no exceeded when MaxPaths=0")
	}
	if report.HasWarnings() {
		t.Error("expected no warnings when WarnAt=0")
	}
}

func TestEvaluate_SortedByEnvironment(t *testing.T) {
	envPaths := map[string][]string{
		"prod":    {"a"},
		"dev":     {"b"},
		"staging": {"c"},
	}
	report := Evaluate(envPaths, Policy{MaxPaths: 5})
	names := make([]string, len(report.Statuses))
	for i, s := range report.Statuses {
		names[i] = s.Environment
	}
	if names[0] != "dev" || names[1] != "prod" || names[2] != "staging" {
		t.Errorf("expected sorted order, got %v", names)
	}
}

func TestSummary_OK(t *testing.T) {
	report := Evaluate(map[string][]string{"dev": {"a"}}, Policy{MaxPaths: 10})
	s := Summary(report)
	if !strings.Contains(s, "OK") {
		t.Errorf("expected OK in summary, got: %s", s)
	}
}

func TestSummary_WithIssues(t *testing.T) {
	envPaths := map[string][]string{
		"prod": {"a", "b", "c", "d"},
		"dev":  {"x", "y", "z"},
	}
	report := Evaluate(envPaths, Policy{MaxPaths: 2, WarnAt: 2})
	s := Summary(report)
	if !strings.Contains(s, "exceeded") {
		t.Errorf("expected 'exceeded' in summary, got: %s", s)
	}
}
