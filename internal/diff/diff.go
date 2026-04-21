package diff

import "sort"

// Result holds the outcome of comparing two sets of secret paths.
type Result struct {
	OnlyInA []string
	OnlyInB []string
	InBoth  []string
}

// IsEmpty reports whether the Result contains no differences between
// the two sets (i.e. OnlyInA and OnlyInB are both empty).
func (r Result) IsEmpty() bool {
	return len(r.OnlyInA) == 0 && len(r.OnlyInB) == 0
}

// Compare returns a Result describing the symmetric difference between
// two slices of secret paths (e.g. from two environments).
func Compare(a, b []string) Result {
	aSet := toSet(a)
	bSet := toSet(b)

	var res Result
	for k := range aSet {
		if bSet[k] {
			res.InBoth = append(res.InBoth, k)
		} else {
			res.OnlyInA = append(res.OnlyInA, k)
		}
	}
	for k := range bSet {
		if !aSet[k] {
			res.OnlyInB = append(res.OnlyInB, k)
		}
	}

	sort.Strings(res.OnlyInA)
	sort.Strings(res.OnlyInB)
	sort.Strings(res.InBoth)
	return res
}

func toSet(paths []string) map[string]bool {
	s := make(map[string]bool, len(paths))
	for _, p := range paths {
		s[p] = true
	}
	return s
}
