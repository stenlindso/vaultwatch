// Package clamp provides utilities for bounding numeric values within
// configurable min/max ranges, used to normalize scores and metrics across
// environments in vaultwatch reports.
package clamp

import "fmt"

// Options configures the clamping behavior.
type Options struct {
	Min float64
	Max float64
}

// DefaultOptions returns a standard 0–100 range suitable for scoring.
func DefaultOptions() Options {
	return Options{Min: 0, Max: 100}
}

// Clamper applies min/max bounding to float64 values.
type Clamper struct {
	opts Options
}

// New creates a Clamper with the given options.
// Returns an error if Min >= Max.
func New(opts Options) (*Clamper, error) {
	if opts.Min >= opts.Max {
		return nil, fmt.Errorf("clamp: Min (%.2f) must be less than Max (%.2f)", opts.Min, opts.Max)
	}
	return &Clamper{opts: opts}, nil
}

// One clamps a single value to [Min, Max].
func (c *Clamper) One(v float64) float64 {
	if v < c.opts.Min {
		return c.opts.Min
	}
	if v > c.opts.Max {
		return c.opts.Max
	}
	return v
}

// All clamps each value in the slice and returns a new slice.
func (c *Clamper) All(vals []float64) []float64 {
	out := make([]float64, len(vals))
	for i, v := range vals {
		out[i] = c.One(v)
	}
	return out
}

// Normalize maps v from [Min, Max] to [0.0, 1.0].
// Values outside the range are clamped before normalizing.
func (c *Clamper) Normalize(v float64) float64 {
	clamped := c.One(v)
	return (clamped - c.opts.Min) / (c.opts.Max - c.opts.Min)
}
