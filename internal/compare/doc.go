// Package compare implements multi-environment secret path matrix comparison.
//
// It builds a pairwise diff matrix from a set of environment snapshots,
// identifying paths unique to each environment and paths shared across them.
// This is useful for auditing consistency across dev, staging, and production.
package compare
