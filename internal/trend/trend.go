// Package trend analyzes change frequency and direction across snapshot history.
package trend

import (
	"sort"
	"time"
)

// Direction indicates whether secrets are growing, shrinking, or stable.
type Direction string

const (
	DirectionGrowing  Direction = "growing"
	DirectionShrinking Direction = "shrinking"
	DirectionStable   Direction = "stable"
)

// DataPoint represents the secret path count at a point in time.
type DataPoint struct {
	Timestamp time.Time
	Count     int
	Env       string
}

// Result holds the computed trend for a single environment.
type Result struct {
	Env       string
	Direction Direction
	Delta     int // net change from first to last data point
	AvgChange float64 // average change per interval
	Samples   int
}

// Analyze computes the trend for each environment from a slice of data points.
// Points for each environment are sorted by timestamp before analysis.
func Analyze(points []DataPoint) []Result {
	byEnv := make(map[string][]DataPoint)
	for _, p := range points {
		byEnv[p.Env] = append(byEnv[p.Env], p)
	}

	results := make([]Result, 0, len(byEnv))
	for env, pts := range byEnv {
		sort.Slice(pts, func(i, j int) bool {
			return pts[i].Timestamp.Before(pts[j].Timestamp)
		})
		results = append(results, compute(env, pts))
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Env < results[j].Env
	})
	return results
}

func compute(env string, pts []DataPoint) Result {
	if len(pts) < 2 {
		return Result{Env: env, Direction: DirectionStable, Samples: len(pts)}
	}

	delta := pts[len(pts)-1].Count - pts[0].Count
	intervals := len(pts) - 1
	avg := float64(delta) / float64(intervals)

	var dir Direction
	switch {
	case delta > 0:
		dir = DirectionGrowing
	case delta < 0:
		dir = DirectionShrinking
	default:
		dir = DirectionStable
	}

	return Result{
		Env:       env,
		Direction: dir,
		Delta:     delta,
		AvgChange: avg,
		Samples:   len(pts),
	}
}
