package cache

import (
	"testing"
	"time"
)

func TestSet_And_Get_Valid(t *testing.T) {
	c := New(5 * time.Minute)
	paths := []string{"secret/a", "secret/b"}
	c.Set("staging", paths)

	got, ok := c.Get("staging")
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if len(got) != 2 {
		t.Errorf("expected 2 paths, got %d", len(got))
	}
}

func TestGet_Missing(t *testing.T) {
	c := New(5 * time.Minute)
	_, ok := c.Get("nonexistent")
	if ok {
		t.Error("expected miss for nonexistent key")
	}
}

func TestGet_Stale(t *testing.T) {
	c := New(1 * time.Millisecond)
	c.Set("prod", []string{"secret/x"})
	time.Sleep(5 * time.Millisecond)

	_, ok := c.Get("prod")
	if ok {
		t.Error("expected stale entry to be rejected")
	}
}

func TestInvalidate(t *testing.T) {
	c := New(5 * time.Minute)
	c.Set("dev", []string{"secret/dev"})
	c.Invalidate("dev")

	_, ok := c.Get("dev")
	if ok {
		t.Error("expected entry to be removed after invalidation")
	}
}

func TestFlush(t *testing.T) {
	c := New(5 * time.Minute)
	c.Set("env1", []string{"secret/a"})
	c.Set("env2", []string{"secret/b"})
	c.Flush()

	if c.Size() != 0 {
		t.Errorf("expected size 0 after flush, got %d", c.Size())
	}
}

func TestSize(t *testing.T) {
	c := New(5 * time.Minute)
	if c.Size() != 0 {
		t.Errorf("expected initial size 0, got %d", c.Size())
	}
	c.Set("a", []string{})
	c.Set("b", []string{})
	if c.Size() != 2 {
		t.Errorf("expected size 2, got %d", c.Size())
	}
}
