// Package drift provides drift detection for Vault secret paths.
//
// It compares a labeled baseline snapshot of paths against a current
// set of paths for a given environment, producing a Result that
// describes which paths were added or removed since the baseline was
// captured.
//
// Typical usage:
//
//	detector := drift.NewDetector("prod", "v1")
//	result := detector.Detect(baselinePaths, currentPaths)
//	if result.HasDrift() {
//	    fmt.Println(result.Summary())
//	}
package drift
