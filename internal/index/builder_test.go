package index

import (
	"testing"
)

func TestBuildFromSnapshot_Basic(t *testing.T) {
	data := map[string][]string{
		"prod":    {"secret/a", "secret/b"},
		"staging": {"secret/a", "secret/c"},
	}
	idx := BuildFromSnapshot(data)

	prodPaths := idx.PathsForEnv("prod")
	if len(prodPaths) != 2 {
		t.Errorf("expected 2 prod paths, got %d", len(prodPaths))
	}
	envs := idx.EnvsForPath("secret/a")
	if len(envs) != 2 {
		t.Errorf("expected 2 envs for secret/a, got %d", len(envs))
	}
}

func TestBuildFromSnapshot_Empty(t *testing.T) {
	idx := BuildFromSnapshot(map[string][]string{})
	if len(idx.Environments()) != 0 {
		t.Error("expected empty index from empty snapshot")
	}
}

func TestMerge_CombinesIndexes(t *testing.T) {
	a := New()
	a.Add("prod", "secret/a")

	b := New()
	b.Add("staging", "secret/b")
	b.Add("prod", "secret/c")

	merged := Merge(a, b)
	envs := merged.Environments()
	if len(envs) != 2 {
		t.Errorf("expected 2 environments after merge, got %d", len(envs))
	}
	prodPaths := merged.PathsForEnv("prod")
	if len(prodPaths) != 2 {
		t.Errorf("expected 2 prod paths after merge, got %d", len(prodPaths))
	}
}

func TestMerge_Deduplicates(t *testing.T) {
	a := New()
	a.Add("prod", "secret/a")

	b := New()
	b.Add("prod", "secret/a")

	merged := Merge(a, b)
	paths := merged.PathsForEnv("prod")
	if len(paths) != 1 {
		t.Errorf("expected deduplication, got %d paths", len(paths))
	}
}

func TestSubset_FiltersEnvironments(t *testing.T) {
	idx := New()
	idx.Add("prod", "secret/a")
	idx.Add("staging", "secret/b")
	idx.Add("dev", "secret/c")

	sub := Subset(idx, []string{"prod", "dev"})
	envs := sub.Environments()
	if len(envs) != 2 {
		t.Errorf("expected 2 environments in subset, got %d", len(envs))
	}
	if len(sub.PathsForEnv("staging")) != 0 {
		t.Error("staging should not be in subset")
	}
}
