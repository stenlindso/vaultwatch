// Package snapshot provides types and utilities for capturing, persisting,
// and loading point-in-time records of Vault secret paths per environment.
//
// A Snapshot records the environment name, capture timestamp, and the full
// list of discovered paths. Snapshots are stored as JSON files managed by
// a Manager, which organises them in a local directory keyed by environment
// name. This allows vaultwatch to compare a live Vault state against a
// previously saved baseline.
package snapshot
