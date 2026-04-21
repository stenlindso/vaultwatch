// Package filter provides flexible path filtering for Vault secret paths.
//
// Filters are composed of one or more rules, each specifying a match type
// (prefix, suffix, or regex) and whether matching paths should be included
// or excluded from results.
//
// Exclusion rules take priority. If no inclusion rules are defined, all paths
// are included by default (only exclusions are applied).
package filter
