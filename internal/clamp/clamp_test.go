package clamp

import (
	"testing"
)

func TestNew_ValidOptions(t *testing.T) {
	c, err := New(Options{Min: 0, Max: 100})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Clamper")
	}
}

func TestNew_InvalidOptions(t *testing.T) {
	_, err := New(Options{Min: 100, Max: 0})
	if err == nil {
		t.Fatal("expected error for Min >= Max")
	}
}

func TestNew_EqualMinMax(t *testing.T) {
	_, err := New(Options{Min: 50, Max: 50})
	if err == nil {
		t.Fatal("expected error when Min == Max")
	}
}

func TestOne_WithinRange(t *testing.T) {
	c, _ := New(DefaultOptions())
	if got := c.One(50); got != 50 {
		t.Errorf("expected 50, got %.2f", got)
	}
}

func TestOne_BelowMin(t *testing.T) {
	c, _ := New(DefaultOptions())
	if got := c.One(-10); got != 0 {
		t.Errorf("expected 0, got %.2f", got)
	}
}

func TestOne_AboveMax(t *testing.T) {
	c, _ := New(DefaultOptions())
	if got := c.One(150); got != 100 {
		t.Errorf("expected 100, got %.2f", got)
	}
}

func TestAll_ClampsEachValue(t *testing.T) {
	c, _ := New(DefaultOptions())
	input := []float64{-5, 50, 120}
	got := c.All(input)
	expected := []float64{0, 50, 100}
	for i, v := range expected {
		if got[i] != v {
			t.Errorf("index %d: expected %.2f, got %.2f", i, v, got[i])
		}
	}
}

func TestNormalize_Midpoint(t *testing.T) {
	c, _ := New(DefaultOptions())
	if got := c.Normalize(50); got != 0.5 {
		t.Errorf("expected 0.5, got %.4f", got)
	}
}

func TestNormalize_AtMin(t *testing.T) {
	c, _ := New(DefaultOptions())
	if got := c.Normalize(0); got != 0.0 {
		t.Errorf("expected 0.0, got %.4f", got)
	}
}

func TestNormalize_AtMax(t *testing.T) {
	c, _ := New(DefaultOptions())
	if got := c.Normalize(100); got != 1.0 {
		t.Errorf("expected 1.0, got %.4f", got)
	}
}

func TestNormalize_BelowMin_ClampsFirst(t *testing.T) {
	c, _ := New(DefaultOptions())
	if got := c.Normalize(-999); got != 0.0 {
		t.Errorf("expected 0.0 after clamping, got %.4f", got)
	}
}
