// Package compare provides multi-environment path comparison across snapshots.
package compare

import (
	"fmt"
	"sort"

	"github.com/vaultwatch/internal/snapshot"
)

// EnvPaths maps environment names to their secret paths.
type EnvPaths map[string][]string

// Matrix holds pairwise diff results between environments.
type Matrix struct {
	Environments []string
	Cells        map[string]Cell
}

// Cell represents the diff between two environments.
type Cell struct {
	Left      string
	Right     string
	OnlyLeft  []string
	OnlyRight []string
	Shared    int
}

// Key returns a canonical key for a pair of environments.
func Key(a, b string) string {
	if a < b {
		return fmt.Sprintf("%s::%s", a, b)
	}
	return fmt.Sprintf("%s::%s", b, a)
}

// BuildMatrix compares all pairs of environments from the provided snapshots.
func BuildMatrix(envSnapshots map[string]*snapshot.Snapshot) *Matrix {
	envs := make([]string, 0, len(envSnapshots))
	for env := range envSnapshots {
		envs = append(envs, env)
	}
	sort.Strings(envs)

	cells := make(map[string]Cell)

	for i := 0; i < len(envs); i++ {
		for j := i + 1; j < len(envs); j++ {
			left := envs[i]
			right := envs[j]

			lSet := toSet(envSnapshots[left].Paths)
			rSet := toSet(envSnapshots[right].Paths)

			var onlyLeft, onlyRight []string
			shared := 0

			for p := range lSet {
				if rSet[p] {
					shared++
				} else {
					onlyLeft = append(onlyLeft, p)
				}
			}
			for p := range rSet {
				if !lSet[p] {
					onlyRight = append(onlyRight, p)
				}
			}

			sort.Strings(onlyLeft)
			sort.Strings(onlyRight)

			cells[Key(left, right)] = Cell{
				Left:      left,
				Right:     right,
				OnlyLeft:  onlyLeft,
				OnlyRight: onlyRight,
				Shared:    shared,
			}
		}
	}

	return &Matrix{Environments: envs, Cells: cells}
}

func toSet(paths []string) map[string]bool {
	s := make(map[string]bool, len(paths))
	for _, p := range paths {
		s[p] = true
	}
	return s
}
