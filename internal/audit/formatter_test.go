package audit_test

import (
	"strings"
	"testing"
	"time"

	"github.com/vaultwatch/internal/audit"
	"github.com/vaultwatch/internal/diff"
)

func makeReport() *audit.Report {
	return &audit.Report{
		EnvironmentA: "prod",
		EnvironmentB: "staging",
		Timestamp:    time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Result: diff.Result{
			OnlyInA: []string{"secret/prod-only"},
			OnlyInB: []string{"secret/staging-only"},
			InBoth:  []string{"secret/shared"},
		},
	}
}

func TestFormat_ContainsEnvironmentNames(t *testing.T) {
	var sb strings.Builder
	audit.Format(&sb, makeReport())
	out := sb.String()

	for _, want := range []string{"prod", "staging", "secret/prod-only", "secret/staging-only", "secret/shared"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestFormat_EmptySections(t *testing.T) {
	var sb strings.Builder
	r := &audit.Report{
		EnvironmentA: "a",
		EnvironmentB: "b",
		Timestamp:    time.Now(),
		Result:       diff.Result{InBoth: []string{"secret/x"}},
	}
	audit.Format(&sb, r)
	out := sb.String()

	if !strings.Contains(out, "(none)") {
		t.Error("expected (none) for empty sections")
	}
}

func TestFormat_Timestamp(t *testing.T) {
	var sb strings.Builder
	audit.Format(&sb, makeReport())
	out := sb.String()

	if !strings.Contains(out, "2024-01-15") {
		t.Errorf("expected formatted timestamp in output, got:\n%s", out)
	}
}
