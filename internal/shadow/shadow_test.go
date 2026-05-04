package shadow

import (
	"testing"
)

func makeSnapshots(data map[string][]string) map[string][]string {
	return data
}

func TestAnalyze_NoShadowing(t *testing.T) {
	snaps := makeSnapshots(map[string][]string{
		"prod": {"secret/app/db"},
		"staging": {"secret/app/cache"},
	})
	r := Analyze([]string{"prod", "staging"}, snaps)
	if r.HasShadows() {
		t.Errorf("expected no shadows, got %v", r.Shadowed)
	}
	if len(r.Unique) != 2 {
		t.Errorf("expected 2 unique paths, got %d", len(r.Unique))
	}
}

func TestAnalyze_ShadowedPath(t *testing.T) {
	snaps := makeSnapshots(map[string][]string{
		"prod":    {"secret/app/db", "secret/app/token"},
		"staging": {"secret/app/db", "secret/app/cache"},
	})
	r := Analyze([]string{"prod", "staging"}, snaps)
	if !r.HasShadows() {
		t.Fatal("expected shadows to be detected")
	}
	shadowedEnvs, ok := r.Shadowed["secret/app/db"]
	if !ok {
		t.Fatal("expected secret/app/db to be shadowed")
	}
	if len(shadowedEnvs) != 1 || shadowedEnvs[0] != "staging" {
		t.Errorf("expected staging to be shadowed, got %v", shadowedEnvs)
	}
}

func TestAnalyze_MultipleShadowedEnvs(t *testing.T) {
	snaps := makeSnapshots(map[string][]string{
		"prod":    {"secret/shared"},
		"staging": {"secret/shared"},
		"dev":     {"secret/shared"},
	})
	r := Analyze([]string{"prod", "staging", "dev"}, snaps)
	envs := r.Shadowed["secret/shared"]
	if len(envs) != 2 {
		t.Errorf("expected 2 shadowed envs, got %d: %v", len(envs), envs)
	}
	if envs[0] != "staging" || envs[1] != "dev" {
		t.Errorf("unexpected shadow order: %v", envs)
	}
}

func TestAnalyze_UniquePathsNotShadowed(t *testing.T) {
	snaps := makeSnapshots(map[string][]string{
		"prod": {"secret/only-prod"},
		"dev":  {"secret/only-dev"},
	})
	r := Analyze([]string{"prod", "dev"}, snaps)
	if r.Unique["secret/only-prod"] != "prod" {
		t.Errorf("expected secret/only-prod to be unique to prod")
	}
	if r.Unique["secret/only-dev"] != "dev" {
		t.Errorf("expected secret/only-dev to be unique to dev")
	}
}

func TestAnalyze_EmptySnapshots(t *testing.T) {
	r := Analyze([]string{"prod", "staging"}, map[string][]string{})
	if r.HasShadows() {
		t.Error("expected no shadows for empty snapshots")
	}
	if len(r.Unique) != 0 {
		t.Errorf("expected no unique paths, got %d", len(r.Unique))
	}
}
