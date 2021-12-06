package cache

import (
	gcp_compute "google.golang.org/api/compute/v1"
)

// Instance GCP Compute instance wrapper
type Instance struct {
	*gcp_compute.Instance
	ServiceAccounts []*ServiceAccount
	CreatedBy       string
	Clusters        []*Cluster
}

func newInstance(instance *gcp_compute.Instance) *Instance {
	return &Instance{
		Instance: instance,
	}
}

// ClustersLen returns length of Clusters
func (t *Instance) ClustersLen() int {
	return len(t.Clusters)
}
