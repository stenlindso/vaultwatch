// Package rollup provides prefix-level aggregation of Vault secret paths
// across multiple environments.
//
// Given a map of environment names to their respective secret path lists,
// Rollup groups paths by a configurable prefix depth and computes coverage
// metrics showing how many environments share each prefix group.
//
// This is useful for identifying common secret namespaces and spotting
// environments that are missing expected path prefixes.
package rollup
