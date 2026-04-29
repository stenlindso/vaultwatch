package verify

import (
	"strings"
	"testing"
)

func TestCheck_AllPresent(t *testing.T) {
	v := New([]string{"secret/a", "secret/b"})
	r := v.Check("prod", []string{"secret/a", "secret/b", "secret/c"})

	if r.Environment != "prod" {
		t.Errorf("expected env prod, got %s", r.Environment)
	}
	if len(r.Present) != 2 {
		t.Errorf("expected 2 present, got %d", len(r.Present))
	}
	if len(r.Missing) != 0 {
		t.Errorf("expected 0 missing, got %d", len(r.Missing))
	}
	if r.HasMissing() {
		t.Error("HasMissing should be false")
	}
}

func TestCheck_SomeMissing(t *testing.T) {
	v := New([]string{"secret/a", "secret/b", "secret/c"})
	r := v.Check("staging", []string{"secret/a"})

	if !r.HasMissing() {
		t.Error("HasMissing should be true")
	}
	if len(r.Missing) != 2 {
		t.Errorf("expected 2 missing, got %d", len(r.Missing))
	}
}

func TestCheck_EmptyActual(t *testing.T) {
	v := New([]string{"secret/x"})
	r := v.Check("dev", []string{})

	if len(r.Missing) != 1 || r.Missing[0] != "secret/x" {
		t.Errorf("expected secret/x to be missing, got %v", r.Missing)
	}
}

func TestCheckAll_MultipleEnvs(t *testing.T) {
	v := New([]string{"secret/a", "secret/b"})
	envPaths := map[string][]string{
		"prod":    {"secret/a", "secret/b"},
		"staging": {"secret/a"},
	}
	results := v.CheckAll(envPaths)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	// results are sorted by env name: prod, staging
	if results[0].Environment != "prod" {
		t.Errorf("expected prod first, got %s", results[0].Environment)
	}
	if results[1].HasMissing() == false {
		t.Error("staging should have missing paths")
	}
}

func TestSummary_Format(t *testing.T) {
	v := New([]string{"secret/a"})
	r := v.Check("prod", []string{"secret/a"})
	s := r.Summary()
	if !strings.Contains(s, "prod") {
		t.Errorf("summary missing env name: %s", s)
	}
}

func TestFormatText_ContainsMissing(t *testing.T) {
	v := New([]string{"secret/a", "secret/b"})
	results := []Result{v.Check("prod", []string{"secret/a"})}
	out := FormatText(results)
	if !strings.Contains(out, "secret/b") {
		t.Error("expected missing path in output")
	}
	if !strings.Contains(out, "✗") {
		t.Error("expected ✗ marker for missing path")
	}
}

func TestFormatOneLiner_StatusFail(t *testing.T) {
	v := New([]string{"secret/a"})
	results := []Result{v.Check("prod", []string{})}
	line := FormatOneLiner(results)
	if !strings.Contains(line, "FAIL") {
		t.Errorf("expected FAIL status, got: %s", line)
	}
}

func TestFormatOneLiner_StatusOK(t *testing.T) {
	v := New([]string{"secret/a"})
	results := []Result{v.Check("prod", []string{"secret/a"})}
	line := FormatOneLiner(results)
	if !strings.Contains(line, "OK") {
		t.Errorf("expected OK status, got: %s", line)
	}
}
