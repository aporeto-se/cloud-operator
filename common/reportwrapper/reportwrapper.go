package reportwrapper

import (
	"sync"

	"github.com/aporeto-se/cloud-operator/common/types"
)

// Wrapper ...
type Wrapper struct {
	reports *types.KubernetesReports
	sync.Mutex
}

// NewWrapper ...
func NewWrapper() *Wrapper {
	return &Wrapper{
		reports: types.NewKubernetesReports(),
	}
}

// AddKube ...
func (t *Wrapper) AddKube(report *types.KubernetesReport) {
	t.Lock()
	defer t.Unlock()
	t.reports.AddReports(report)
}

// Build ...
func (t *Wrapper) Build() *types.KubernetesReports {
	t.Lock()
	defer t.Unlock()
	return t.reports
}
