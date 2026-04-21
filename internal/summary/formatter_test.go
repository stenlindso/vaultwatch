package summary

import (
	"strings"
	"testing"
	"time"
)

func buildReport() *Report {
	return Summarize([]Input{
		{
			Environment: "prod",
			Paths:       []string{"secret/a", "secret/b"},
			Added:       []string{"secret/b"},
			Removed:     []string{},
			LastScanned: time.Now().UTC(),
		},
		{
			Environment: "staging",
			Paths:       []string{"secret/a"},
			Added:       []string{},
			Removed:     []string{"secret/c"},
			LastScanned: time.Now().UTC(),
		},
	})
}

func TestFormatText_ContainsEnvironments(t *testing.T) {
	r := buildReport()
	var sb strings.Builder
	FormatText(&sb, r)
	out := sb.String()

	for _, env := range []string{"prod", "staging"} {
		if !strings.Contains(out, env) {
			t.Errorf("expected output to contain %q", env)
		}
	}
}

func TestFormatText_ContainsTotals(t *testing.T) {
	r := buildReport()
	var sb strings.Builder
	FormatText(&sb, r)
	out := sb.String()

	if !strings.Contains(out, "Total paths") {
		t.Error("expected 'Total paths' in output")
	}
}

func TestFormatText_EmptyReport(t *testing.T) {
	r := Summarize(nil)
	var sb strings.Builder
	FormatText(&sb, r)
	out := sb.String()

	if !strings.Contains(out, "No environments") {
		t.Error("expected 'No environments' message for empty report")
	}
}

func TestFormatOneLiner_NoChanges(t *testing.T) {
	r := Summarize([]Input{
		{Environment: "prod", Paths: []string{"a"}, Added: nil, Removed: nil},
	})
	line := FormatOneLiner(r)
	if !strings.Contains(line, "no changes") {
		t.Errorf("expected 'no changes' in one-liner, got: %s", line)
	}
}

func TestFormatOneLiner_WithChanges(t *testing.T) {
	r := buildReport()
	line := FormatOneLiner(r)
	if !strings.Contains(line, "CHANGED") {
		t.Errorf("expected 'CHANGED' in one-liner, got: %s", line)
	}
}
