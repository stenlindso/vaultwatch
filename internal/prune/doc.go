// Package prune evaluates snapshot entries against age and count thresholds
// to determine which entries are stale and eligible for removal.
//
// It does not perform any I/O itself; callers are responsible for deleting
// the pruned entries from their respective stores.
package prune
