// Package alert provides threshold-based alerting for Vault path changes
// detected during watch cycles or audit runs.
package alert

import (
	"fmt"
	"time"
)

// Severity represents the urgency level of an alert.
type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"
)

// Alert represents a triggered threshold alert.
type Alert struct {
	Environment string
	Severity    Severity
	Message     string
	PathCount   int
	TriggeredAt time.Time
}

// Threshold defines the rules for triggering alerts based on path change counts.
type Threshold struct {
	WarningCount  int
	CriticalCount int
}

// DefaultThreshold returns a sensible default threshold configuration.
func DefaultThreshold() Threshold {
	return Threshold{
		WarningCount:  5,
		CriticalCount: 20,
	}
}

// Evaluate checks the number of changed paths against the threshold and
// returns an Alert if a threshold is exceeded, or nil otherwise.
func Evaluate(env string, changedCount int, t Threshold) *Alert {
	var sev Severity
	var msg string

	switch {
	case changedCount >= t.CriticalCount:
		sev = SeverityCritical
		msg = fmt.Sprintf("%d paths changed in %q — exceeds critical threshold (%d)", changedCount, env, t.CriticalCount)
	case changedCount >= t.WarningCount:
		sev = SeverityWarning
		msg = fmt.Sprintf("%d paths changed in %q — exceeds warning threshold (%d)", changedCount, env, t.WarningCount)
	default:
		return nil
	}

	return &Alert{
		Environment: env,
		Severity:    sev,
		Message:     msg,
		PathCount:   changedCount,
		TriggeredAt: time.Now().UTC(),
	}
}
