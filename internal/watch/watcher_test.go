package watch

import (
	"testing"
	"time"
)

func TestDiff_Added(t *testing.T) {
	prev := []string{"secret/a", "secret/b"}
	curr := []string{"secret/a", "secret/b", "secret/c"}
	ev := diff("staging", prev, curr)
	if len(ev.Added) != 1 || ev.Added[0] != "secret/c" {
		t.Errorf("expected Added=[secret/c], got %v", ev.Added)
	}
	if len(ev.Removed) != 0 {
		t.Errorf("expected no removals, got %v", ev.Removed)
	}
}

func TestDiff_Removed(t *testing.T) {
	prev := []string{"secret/a", "secret/b"}
	curr := []string{"secret/a"}
	ev := diff("prod", prev, curr)
	if len(ev.Removed) != 1 || ev.Removed[0] != "secret/b" {
		t.Errorf("expected Removed=[secret/b], got %v", ev.Removed)
	}
	if len(ev.Added) != 0 {
		t.Errorf("expected no additions, got %v", ev.Added)
	}
}

func TestDiff_NoChange(t *testing.T) {
	paths := []string{"secret/x", "secret/y"}
	ev := diff("dev", paths, paths)
	if len(ev.Added) != 0 || len(ev.Removed) != 0 {
		t.Errorf("expected no changes, got added=%v removed=%v", ev.Added, ev.Removed)
	}
}

func TestDiff_EnvAndTimestamp(t *testing.T) {
	before := time.Now()
	ev := diff("qa", nil, nil)
	if ev.Env != "qa" {
		t.Errorf("expected env=qa, got %s", ev.Env)
	}
	if ev.At.Before(before) {
		t.Error("expected At to be set to current time")
	}
}
