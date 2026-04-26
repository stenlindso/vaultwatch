package classify

import (
	"strings"
	"testing"
)

func buildGrouped() map[Level][]Result {
	return map[Level][]Result{
		LevelCritical: {
			{Path: "prod/app/db-password", Level: LevelCritical, Rule: `password`},
		},
		LevelSecret: {
			{Path: "prod/app/tls-key", Level: LevelSecret, Rule: `key`},
			{Path: "prod/app/api-cert", Level: LevelSecret, Rule: `cert`},
		},
		LevelPublic: {
			{Path: "prod/app/version", Level: LevelPublic},
		},
	}
}

func TestFormatText_ContainsLevels(t *testing.T) {
	var sb strings.Builder
	FormatText(&sb, buildGrouped())
	out := sb.String()
	for _, want := range []string{"CRITICAL", "SECRET", "PUBLIC"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}

func TestFormatText_ContainsPaths(t *testing.T) {
	var sb strings.Builder
	FormatText(&sb, buildGrouped())
	out := sb.String()
	if !strings.Contains(out, "prod/app/db-password") {
		t.Error("expected path in output")
	}
	if !strings.Contains(out, "matched:") {
		t.Error("expected matched rule annotation in output")
	}
}

func TestFormatText_EmptyGroupsOmitted(t *testing.T) {
	var sb strings.Builder
	FormatText(&sb, map[Level][]Result{})
	out := sb.String()
	if strings.Contains(out, "CRITICAL") || strings.Contains(out, "PUBLIC") {
		t.Error("expected no output for empty groups")
	}
}

func TestFormatSummary_Counts(t *testing.T) {
	grouped := buildGrouped()
	summary := FormatSummary(grouped)
	if !strings.Contains(summary, "critical=1") {
		t.Errorf("expected critical=1 in summary, got %q", summary)
	}
	if !strings.Contains(summary, "secret=2") {
		t.Errorf("expected secret=2 in summary, got %q", summary)
	}
	if !strings.Contains(summary, "public=1") {
		t.Errorf("expected public=1 in summary, got %q", summary)
	}
}

func TestFormatSummary_AllZero(t *testing.T) {
	summary := FormatSummary(map[Level][]Result{})
	if !strings.Contains(summary, "critical=0") {
		t.Errorf("expected zeros in summary, got %q", summary)
	}
}
