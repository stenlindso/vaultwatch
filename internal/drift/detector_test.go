package drift

import (
	"strings"
	"testing"
)

func TestDetect_NoDrift(t *testing.T) {
	d := NewDetector("prod", "v1")
	paths := []string{"secret/app/db", "secret/app/api"}
	r := d.Detect(paths, paths)

	if r.HasDrift() {
		t.Errorf("expected no drift, got added=%v removed=%v", r.Added, r.Removed)
	}
	if r.Environment != "prod" {
		t.Errorf("expected env prod, got %s", r.Environment)
	}
	if r.Label != "v1" {
		t.Errorf("expected label v1, got %s", r.Label)
	}
}

func TestDetect_Added(t *testing.T) {
	d := NewDetector("staging", "base")
	baseline := []string{"secret/app/db"}
	current := []string{"secret/app/db", "secret/app/cache"}

	r := d.Detect(baseline, current)

	if !r.HasDrift() {
		t.Fatal("expected drift")
	}
	if len(r.Added) != 1 || r.Added[0] != "secret/app/cache" {
		t.Errorf("unexpected added paths: %v", r.Added)
	}
	if len(r.Removed) != 0 {
		t.Errorf("expected no removed paths, got: %v", r.Removed)
	}
}

func TestDetect_Removed(t *testing.T) {
	d := NewDetector("dev", "initial")
	baseline := []string{"secret/app/db", "secret/app/legacy"}
	current := []string{"secret/app/db"}

	r := d.Detect(baseline, current)

	if len(r.Removed) != 1 || r.Removed[0] != "secret/app/legacy" {
		t.Errorf("unexpected removed paths: %v", r.Removed)
	}
	if len(r.Added) != 0 {
		t.Errorf("expected no added paths, got: %v", r.Added)
	}
}

func TestDetect_BothDirections(t *testing.T) {
	d := NewDetector("prod", "v2")
	baseline := []string{"secret/a", "secret/b"}
	current := []string{"secret/b", "secret/c"}

	r := d.Detect(baseline, current)

	if len(r.Added) != 1 || len(r.Removed) != 1 {
		t.Errorf("expected 1 added and 1 removed, got +%d -%d", len(r.Added), len(r.Removed))
	}
}

func TestSummary_NoDrift(t *testing.T) {
	d := NewDetector("prod", "v1")
	r := d.Detect([]string{"secret/x"}, []string{"secret/x"})

	if !strings.Contains(r.Summary(), "no drift") {
		t.Errorf("expected 'no drift' in summary, got: %s", r.Summary())
	}
}

func TestSummary_WithDrift(t *testing.T) {
	d := NewDetector("prod", "v1")
	r := d.Detect([]string{"secret/old"}, []string{"secret/new"})

	if !strings.Contains(r.Summary(), "drift detected") {
		t.Errorf("expected 'drift detected' in summary, got: %s", r.Summary())
	}
}

func TestDetect_EmptyBaseline(t *testing.T) {
	d := NewDetector("qa", "empty")
	r := d.Detect([]string{}, []string{"secret/app/new"})

	if len(r.Added) != 1 {
		t.Errorf("expected 1 added path, got %d", len(r.Added))
	}
}
