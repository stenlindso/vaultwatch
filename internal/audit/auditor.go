package audit

import (
	"fmt"
	"time"

	"github.com/vaultwatch/internal/diff"
	"github.com/vaultwatch/internal/snapshot"
)

// Report holds the result of an audit between two environments.
type Report struct {
	EnvironmentA string
	EnvironmentB string
	Timestamp    time.Time
	Result       diff.Result
}

// Auditor compares snapshots across environments.
type Auditor struct {
	manager *snapshot.Manager
}

// NewAuditor creates an Auditor backed by the given Manager.
func NewAuditor(m *snapshot.Manager) *Auditor {
	return &Auditor{manager: m}
}

// Audit loads snapshots for envA and envB and returns a diff Report.
func (a *Auditor) Audit(envA, envB string) (*Report, error) {
	snapshotA, err := a.manager.Load(envA)
	if err != nil {
		return nil, fmt.Errorf("loading snapshot for %q: %w", envA, err)
	}

	snapshotB, err := a.manager.Load(envB)
	if err != nil {
		return nil, fmt.Errorf("loading snapshot for %q: %w", envB, err)
	}

	result := diff.Compare(snapshotA.Paths, snapshotB.Paths)

	return &Report{
		EnvironmentA: envA,
		EnvironmentB: envB,
		Timestamp:    time.Now().UTC(),
		Result:       result,
	}, nil
}
