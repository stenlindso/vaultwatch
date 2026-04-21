package baseline_test

import (
	"strings"
	"testing"
	"time"

	"github.com/youorg/vaultwatch/internal/baseline"
)

func makeDiffResult(added, removed []string) *baseline.DiffResult {
	return &baseline.DiffResult{
		Environment: "production",
		Label:       "v1",
		Added:       added,
		Removed:     removed,
		CheckedAt:   time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
	}
}

func TestFormatText_NoChanges(t *testing.T) {
	var buf strings.Builder
	d := makeDiffResult(nil, nil)
	baseline.FormatText(&buf, d)

	out := buf.String()
	if !strings.Contains(out, "No changes from baseline") {
		t.Errorf("expected no-changes message, got:\n%s", out)
	}
}

func TestFormatText_ContainsEnvAndLabel(t *testing.T) {
	var buf strings.Builder
	d := makeDiffResult([]string{"secret/new"}, nil)
	baseline.FormatText(&buf, d)

	out := buf.String()
	if !strings.Contains(out, "production") {
		t.Errorf("expected environment name in output")
	}
	if !strings.Contains(out, "v1") {
		t.Errorf("expected label in output")
	}
}

func TestFormatText_AddedPaths(t *testing.T) {
	var buf strings.Builder
	d := makeDiffResult([]string{"secret/x", "secret/y"}, nil)
	baseline.FormatText(&buf, d)

	out := buf.String()
	if !strings.Contains(out, "+ secret/x") {
		t.Errorf("expected added path secret/x in output")
	}
	if !strings.Contains(out, "+ secret/y") {
		t.Errorf("expected added path secret/y in output")
	}
}

func TestFormatText_RemovedPaths(t *testing.T) {
	var buf strings.Builder
	d := makeDiffResult(nil, []string{"secret/old"})
	baseline.FormatText(&buf, d)

	out := buf.String()
	if !strings.Contains(out, "- secret/old") {
		t.Errorf("expected removed path in output")
	}
}

func TestHasChanges_True(t *testing.T) {
	d := makeDiffResult([]string{"secret/a"}, nil)
	if !d.HasChanges() {
		t.Error("expected HasChanges to return true")
	}
}

func TestHasChanges_False(t *testing.T) {
	d := makeDiffResult(nil, nil)
	if d.HasChanges() {
		t.Error("expected HasChanges to return false")
	}
}
