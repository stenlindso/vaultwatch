// Package alert implements threshold-based alerting for VaultWatch.
//
// It evaluates the number of changed secret paths in a given environment
// against configurable warning and critical thresholds. When a threshold
// is exceeded, an Alert is produced and routed to registered handlers via
// the Dispatcher.
//
// Example usage:
//
//	d := alert.NewDispatcher(alert.DefaultThreshold(), alert.StdoutHandler())
//	d.Dispatch("production", len(changedPaths))
package alert
