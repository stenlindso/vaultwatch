// Package audit provides functionality for comparing Vault secret path
// snapshots across environments and formatting the resulting diff reports.
//
// Usage:
//
//	manager := snapshot.NewManager("/var/vaultwatch")
//	auditor := audit.NewAuditor(manager)
//	report, err := auditor.Audit("prod", "staging")
//	if err != nil { ... }
//	audit.Format(os.Stdout, report)
package audit
