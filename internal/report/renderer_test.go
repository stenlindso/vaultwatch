package report_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/vaultwatch/internal/diff"
	"github.com/vaultwatch/internal/policy"
	"github.com/vaultwatch/internal/report"
)

func buildReport() report.Report {
	b := report.NewBuilder()
	b.AddEnvReport(report.EnvReport{
		BaseEnv:   "prod",
		TargetEnv: "staging",
		Diff:      diff.Result{Added: []string{"secret/new"}, Removed: []string{"secret/old"}},
		Violations: []policy.Violation{{Path: "secret/bad", Rule: "deny", Message: "not allowed"}},
	})
	return b.Build()
}

func TestRenderText_ContainsEnvNames(t *testing.T) {
	var buf bytes.Buffer
	err := report.Render(&buf, buildReport(), report.FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"prod", "staging", "secret/new", "secret/old", "not allowed"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}

func TestRenderJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	err := report.Render(&buf, buildReport(), report.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out report.Report
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(out.Environments) != 1 {
		t.Errorf("expected 1 environment, got %d", len(out.Environments))
	}
}

func TestRenderText_NoChanges(t *testing.T) {
	b := report.NewBuilder()
	b.AddEnvReport(report.EnvReport{BaseEnv: "a", TargetEnv: "b"})
	var buf bytes.Buffer
	_ = report.Render(&buf, b.Build(), report.FormatText)
	if !strings.Contains(buf.String(), "No changes or violations") {
		t.Error("expected no-changes message")
	}
}
