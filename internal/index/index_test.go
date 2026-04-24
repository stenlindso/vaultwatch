package index

import (
	"testing"
)

func TestAdd_And_PathsForEnv(t *testing.T) {
	idx := New()
	idx.Add("prod", "secret/db/password")
	idx.Add("prod", "secret/api/key")
	idx.Add("staging", "secret/db/password")

	paths := idx.PathsForEnv("prod")
	if len(paths) != 2 {
		t.Fatalf("expected 2 paths for prod, got %d", len(paths))
	}
	if paths[0] != "secret/api/key" {
		t.Errorf("expected sorted first path, got %q", paths[0])
	}
}

func TestEnvsForPath(t *testing.T) {
	idx := New()
	idx.Add("prod", "secret/db/password")
	idx.Add("staging", "secret/db/password")
	idx.Add("dev", "secret/db/password")

	envs := idx.EnvsForPath("secret/db/password")
	if len(envs) != 3 {
		t.Fatalf("expected 3 envs, got %d", len(envs))
	}
	if envs[0] != "dev" {
		t.Errorf("expected sorted first env to be dev, got %q", envs[0])
	}
}

func TestSearch_ReturnsMatches(t *testing.T) {
	idx := New()
	idx.Add("prod", "secret/db/password")
	idx.Add("prod", "secret/api/key")
	idx.Add("staging", "secret/db/host")

	results := idx.Search("db")
	if len(results) != 2 {
		t.Fatalf("expected 2 results for 'db', got %d", len(results))
	}
	for _, r := range results {
		if r.Path != "secret/db/password" && r.Path != "secret/db/host" {
			t.Errorf("unexpected path in results: %q", r.Path)
		}
	}
}

func TestSearch_NoMatches(t *testing.T) {
	idx := New()
	idx.Add("prod", "secret/db/password")

	results := idx.Search("nonexistent")
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
}

func TestAdd_IgnoresEmpty(t *testing.T) {
	idx := New()
	idx.Add("", "secret/db/password")
	idx.Add("prod", "")

	if len(idx.Environments()) != 0 {
		t.Errorf("expected no environments after empty adds")
	}
}

func TestAdd_DeduplicatesEnvForPath(t *testing.T) {
	idx := New()
	idx.Add("prod", "secret/db/password")
	idx.Add("prod", "secret/db/password")

	envs := idx.EnvsForPath("secret/db/password")
	if len(envs) != 1 {
		t.Errorf("expected 1 env (deduped), got %d", len(envs))
	}
}

func TestStats(t *testing.T) {
	idx := New()
	idx.Add("prod", "secret/a")
	idx.Add("staging", "secret/a")
	idx.Add("prod", "secret/b")

	stats := idx.Stats()
	if stats == "" {
		t.Error("expected non-empty stats string")
	}
}
