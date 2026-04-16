package diff

import (
	"reflect"
	"testing"
)

func TestCompare_Disjoint(t *testing.T) {
	a := []string{"secret/foo", "secret/bar"}
	b := []string{"secret/baz"}
	res := Compare(a, b)

	if !reflect.DeepEqual(res.OnlyInA, []string{"secret/bar", "secret/foo"}) {
		t.Errorf("OnlyInA mismatch: %v", res.OnlyInA)
	}
	if !reflect.DeepEqual(res.OnlyInB, []string{"secret/baz"}) {
		t.Errorf("OnlyInB mismatch: %v", res.OnlyInB)
	}
	if len(res.InBoth) != 0 {
		t.Errorf("expected empty InBoth, got %v", res.InBoth)
	}
}

func TestCompare_Overlap(t *testing.T) {
	a := []string{"secret/shared", "secret/only-a"}
	b := []string{"secret/shared", "secret/only-b"}
	res := Compare(a, b)

	if !reflect.DeepEqual(res.InBoth, []string{"secret/shared"}) {
		t.Errorf("InBoth mismatch: %v", res.InBoth)
	}
	if !reflect.DeepEqual(res.OnlyInA, []string{"secret/only-a"}) {
		t.Errorf("OnlyInA mismatch: %v", res.OnlyInA)
	}
	if !reflect.DeepEqual(res.OnlyInB, []string{"secret/only-b"}) {
		t.Errorf("OnlyInB mismatch: %v", res.OnlyInB)
	}
}

func TestCompare_Identical(t *testing.T) {
	paths := []string{"secret/a", "secret/b"}
	res := Compare(paths, paths)
	if len(res.OnlyInA) != 0 || len(res.OnlyInB) != 0 {
		t.Errorf("expected no differences, got %+v", res)
	}
	if len(res.InBoth) != 2 {
		t.Errorf("expected 2 in common, got %d", len(res.InBoth))
	}
}
