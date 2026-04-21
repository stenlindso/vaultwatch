package normalize

import (
	"testing"
)

func TestOne_StripLeadingSlash(t *testing.T) {
	n := New(Options{StripLeadingSlash: true})
	got := n.One("/secret/foo")
	if got != "secret/foo" {
		t.Errorf("expected 'secret/foo', got %q", got)
	}
}

func TestOne_StripTrailingSlash(t *testing.T) {
	n := New(Options{StripTrailingSlash: true})
	got := n.One("secret/foo/")
	if got != "secret/foo" {
		t.Errorf("expected 'secret/foo', got %q", got)
	}
}

func TestOne_Lowercase(t *testing.T) {
	n := New(Options{Lowercase: true})
	got := n.One("Secret/FOO")
	if got != "secret/foo" {
		t.Errorf("expected 'secret/foo', got %q", got)
	}
}

func TestOne_Clean(t *testing.T) {
	n := New(Options{Clean: true})
	got := n.One("secret//foo/../foo")
	if got != "secret/foo" {
		t.Errorf("expected 'secret/foo', got %q", got)
	}
}

func TestOne_DefaultOptions(t *testing.T) {
	n := New(DefaultOptions())
	got := n.One("/secret//foo/")
	if got != "secret/foo" {
		t.Errorf("expected 'secret/foo', got %q", got)
	}
}

func TestAll_NormalizesEach(t *testing.T) {
	n := New(DefaultOptions())
	input := []string{"/a/b", "/c/d/"}
	got := n.All(input)
	expected := []string{"a/b", "c/d"}
	for i, g := range got {
		if g != expected[i] {
			t.Errorf("index %d: expected %q, got %q", i, expected[i], g)
		}
	}
}

func TestDeduplicate_RemovesDuplicates(t *testing.T) {
	n := New(DefaultOptions())
	input := []string{"/secret/foo", "secret/foo/", "/secret/bar"}
	got := n.Deduplicate(input)
	if len(got) != 2 {
		t.Errorf("expected 2 unique paths, got %d: %v", len(got), got)
	}
}

func TestDeduplicate_PreservesOrder(t *testing.T) {
	n := New(DefaultOptions())
	input := []string{"/b/c", "/a/b", "/b/c"}
	got := n.Deduplicate(input)
	if len(got) != 2 || got[0] != "b/c" || got[1] != "a/b" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestDeduplicate_Empty(t *testing.T) {
	n := New(DefaultOptions())
	got := n.Deduplicate([]string{})
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}
