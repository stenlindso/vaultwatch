package rollup

import (
	"strings"
	"testing"
)

func buildResult() Result {
	return Result{
		TotalEnvs: 3,
		Entries: []Entry{
			{Prefix: "secret/app", Envs: []string{"dev", "prod", "staging"}, Count: 3, Coverage: 1.0},
			{Prefix: "secret/infra", Envs: []string{"prod", "staging"}, Count: 2, Coverage: 0.667},
			{Prefix: "kv/data", Envs: []string{"dev"}, Count: 1, Coverage: 0.333},
		},
	}
}

func TestFormatText_ContainsPrefixes(t *testing.T) {
	var sb strings.Builder
	FormatText(&sb, buildResult())
	out := sb.String()

	for _, prefix := range []string{"secret/app", "secret/infra", "kv/data"} {
		if !strings.Contains(out, prefix) {
			t.Errorf("expected output to contain %q", prefix)
		}
	}
}

func TestFormatText_ContainsEnvCount(t *testing.T) {
	var sb strings.Builder
	FormatText(&sb, buildResult())
	out := sb.String()

	if !strings.Contains(out, "3 environment") {
		t.Errorf("expected env count in header, got:\n%s", out)
	}
}

func TestFormatText_EmptyResult(t *testing.T) {
	var sb strings.Builder
	FormatText(&sb, Result{})
	out := sb.String()

	if !strings.Contains(out, "No rollup data") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormatOneLiner_WithData(t *testing.T) {
	line := FormatOneLiner(buildResult())
	if !strings.Contains(line, "3 prefix group") {
		t.Errorf("unexpected one-liner: %s", line)
	}
	if !strings.Contains(line, "3 env") {
		t.Errorf("expected env count in one-liner: %s", line)
	}
}

func TestFormatOneLiner_Empty(t *testing.T) {
	line := FormatOneLiner(Result{})
	if !strings.Contains(line, "no data") {
		t.Errorf("expected 'no data' in one-liner: %s", line)
	}
}
