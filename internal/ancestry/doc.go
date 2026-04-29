// Package ancestry records and queries parent-child derivation relationships
// between Vault secret paths. It allows vaultwatch to surface orphaned secrets
// (children that no longer have a live parent) and to visualise how secrets
// were derived across environments.
package ancestry
