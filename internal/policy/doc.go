// Package policy provides rule-based auditing for Vault secret paths.
//
// Rules can be of two types:
//   - required: a path matching the pattern must exist in the snapshot
//   - deny: no path matching the pattern may exist in the snapshot
//
// Use LoadFromFile or DefaultRules to obtain a rule set, NewChecker to
// evaluate paths against those rules, and NewEvaluator for higher-level
// per-environment evaluation across multiple snapshots.
package policy
