package overlap

import (
	"testing"
)

func makeSnapshots() map[string][]string {
	return map[string][]string{
		"prod":    {"secret/db/password", "secret/api/key", "secret/shared/token"},
		"staging": {"secret/db/password", "secret/shared/token", "secret/staging/only"},
		"dev":     {"secret/shared/token", "secret/dev/debug"},
	}
}

func TestAnalyze_SharedPaths(t *testing.T) {
	r := Analyze(makeSnapshots())

	envs, ok := r.SharedPaths["secret/shared/token"]
	if !ok {
		t.Fatal("expected secret/shared/token in SharedPaths")
	}
	if len(envs) != 3 {
		t.Errorf("expected 3 envs for shared/token, got %d", len(envs))
	}
}

func TestAnalyze_UniqueByEnv(t *testing.T) {
	r := Analyze(makeSnapshots())

	stagingUniq := r.UniqueByEnv["staging"]
	if len(stagingUniq) != 1 || stagingUniq[0] != "secret/staging/only" {
		t.Errorf("unexpected staging unique paths: %v", stagingUniq)
	}

	devUniq := r.UniqueByEnv["dev"]
	if len(devUniq) != 1 || devUniq[0] != "secret/dev/debug" {
		t.Errorf("unexpected dev unique paths: %v", devUniq)
	}
}

func TestAnalyze_HasOverlap(t *testing.T) {
	r := Analyze(makeSnapshots())
	if !r.HasOverlap() {
		t.Error("expected HasOverlap to be true")
	}
}

func TestAnalyze_NoOverlap(t *testing.T) {
	snaps := map[string][]string{
		"prod": {"secret/a"},
		"dev":  {"secret/b"},
	}
	r := Analyze(snaps)
	if r.HasOverlap() {
		t.Error("expected HasOverlap to be false")
	}
}

func TestSharedAcrossAll(t *testing.T) {
	paths := SharedAcrossAll(makeSnapshots())
	if len(paths) != 1 || paths[0] != "secret/shared/token" {
		t.Errorf("expected only secret/shared/token, got %v", paths)
	}
}

func TestSharedAcrossAll_Empty(t *testing.T) {
	paths := SharedAcrossAll(map[string][]string{})
	if paths != nil {
		t.Errorf("expected nil for empty input, got %v", paths)
	}
}

func TestAnalyze_EmptySnapshots(t *testing.T) {
	r := Analyze(map[string][]string{})
	if r.HasOverlap() {
		t.Error("expected no overlap for empty input")
	}
	if len(r.SharedPaths) != 0 {
		t.Errorf("expected empty SharedPaths, got %v", r.SharedPaths)
	}
}
