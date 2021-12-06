package cache

import (
	gke_service "google.golang.org/api/container/v1"
)

// Cluster ...
type Cluster struct {
	*gke_service.Cluster
	Instances       []*Instance
	ServiceAccounts []*ServiceAccount
	CreatedBy       string
}

func newCluster(cluster *gke_service.Cluster) *Cluster {
	return &Cluster{
		Cluster: cluster,
	}
}
