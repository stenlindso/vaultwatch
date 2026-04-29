package schedule

import (
	"strings"
	"testing"
	"time"
)

func buildEntries() []*Entry {
	base := time.Date(2024, 6, 1, 9, 0, 0, 0, time.UTC)
	return []*Entry{
		{Environment: "prod", Interval: time.Hour, NextRun: base.Add(time.Hour)},
		{Environment: "staging", Interval: 30 * time.Minute, NextRun: base.Add(30 * time.Minute)},
	}
}

func TestFormatText_ContainsEnvironments(t *testing.T) {
	var sb strings.Builder
	FormatText(&sb, buildEntries())
	out := sb.String()
	for _, env := range []string{"prod", "staging"} {
		if !strings.Contains(out, env) {
			t.Errorf("expected output to contain %q", env)
		}
	}
}

func TestFormatText_ContainsInterval(t *testing.T) {
	var sb strings.Builder
	FormatText(&sb, buildEntries())
	out := sb.String()
	if !strings.Contains(out, "1h0m0s") {
		t.Error("expected output to contain interval string")
	}
}

func TestFormatText_EmptyEntries(t *testing.T) {
	var sb strings.Builder
	FormatText(&sb, nil)
	out := sb.String()
	if !strings.Contains(out, "No scheduled") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormatOneLiner_WithData(t *testing.T) {
	line := FormatOneLiner(buildEntries())
	if !strings.Contains(line, "2 environment") {
		t.Errorf("expected count in one-liner, got: %s", line)
	}
	if !strings.Contains(line, "prod") {
		t.Errorf("expected prod in one-liner, got: %s", line)
	}
}

func TestFormatOneLiner_Empty(t *testing.T) {
	line := FormatOneLiner(nil)
	if !strings.Contains(line, "0 environments") {
		t.Errorf("expected zero count, got: %s", line)
	}
}

func TestFormatText_SortedAlphabetically(t *testing.T) {
	var sb strings.Builder
	FormatText(&sb, buildEntries())
	out := sb.String()
	prodIdx := strings.Index(out, "prod")
	stagingIdx := strings.Index(out, "staging")
	if prodIdx > stagingIdx {
		t.Error("expected prod to appear before staging (alphabetical order)")
	}
}
