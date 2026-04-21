// Package tags provides tagging support for Vault secret paths.
//
// Tags allow operators to annotate paths with arbitrary labels (e.g. "pii",
// "critical", "deprecated") and later filter or report on those subsets.
// Tag data is stored as JSON files alongside snapshot data and is keyed by
// environment name.
package tags
