package filter_test

import (
	"testing"

	"github.com/your-org/vaultwatch/internal/filter"
)

func TestApply_PrefixInclusion(t *testing.T) {
	f := filter.New([]filter.Rule{
		{Type: "prefix", Pattern: "secret/prod", Exclude: false},
	})
	paths := []string{"secret/prod/db", "secret/staging/db", "secret/prod/api"}
	got, err := f.Apply(paths)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 paths, got %d: %v", len(got), got)
	}
}

func TestApply_SuffixExclusion(t *testing.T) {
	f := filter.New([]filter.Rule{
		{Type: "suffix", Pattern: "/tmp", Exclude: true},
	})
	paths := []string{"secret/prod/tmp", "secret/prod/db", "secret/staging/tmp"}
	got, err := f.Apply(paths)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 path, got %d: %v", len(got), got)
	}
	if got[0] != "secret/prod/db" {
		t.Errorf("expected secret/prod/db, got %s", got[0])
	}
}

func TestApply_RegexInclusion(t *testing.T) {
	f := filter.New([]filter.Rule{
		{Type: "regex", Pattern: `secret/(prod|staging)/api`, Exclude: false},
	})
	paths := []string{"secret/prod/api", "secret/staging/api", "secret/dev/api"}
	got, err := f.Apply(paths)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 paths, got %d: %v", len(got), got)
	}
}

func TestApply_NoRules_IncludesAll(t *testing.T) {
	f := filter.New(nil)
	paths := []string{"secret/a", "secret/b", "secret/c"}
	got, err := f.Apply(paths)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(paths) {
		t.Errorf("expected all %d paths, got %d", len(paths), len(got))
	}
}

func TestApply_InvalidRegex_ReturnsError(t *testing.T) {
	f := filter.New([]filter.Rule{
		{Type: "regex", Pattern: `[invalid`, Exclude: false},
	})
	_, err := f.Apply([]string{"secret/prod/api"})
	if err == nil {
		t.Error("expected error for invalid regex, got nil")
	}
}

func TestApply_ExclusionOverridesInclusion(t *testing.T) {
	f := filter.New([]filter.Rule{
		{Type: "prefix", Pattern: "secret/prod", Exclude: false},
		{Type: "suffix", Pattern: "/legacy", Exclude: true},
	})
	paths := []string{"secret/prod/api", "secret/prod/legacy", "secret/staging/api"}
	got, err := f.Apply(paths)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0] != "secret/prod/api" {
		t.Errorf("expected [secret/prod/api], got %v", got)
	}
}
