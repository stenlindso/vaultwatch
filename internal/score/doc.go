// Package score computes a weighted health score (0–100) and letter grade
// (A–F) for each Vault environment monitored by vaultwatch.
//
// The score is derived from three signals:
//
//   - Drift ratio: proportion of paths that differ from the reference state.
//   - Policy violations: number of rule violations detected by the policy engine.
//   - Snapshot staleness: how far past the configured maximum age the most
//     recent snapshot is.
//
// Use Compute to obtain a Result for a single environment, then integrate the
// results into reports or alerts via the report and alert packages.
package score
