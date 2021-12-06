package cache

// RoleAccount AWS Role Account
type RoleAccount struct {
	Name             string
	ClusterInstances []*Instance
	ComputeInstances []*Instance
	Clusters         []*Cluster
}

func newRoleAccount(name string) *RoleAccount {
	return &RoleAccount{
		Name: name,
	}
}

// ComputeInstancesLen returns length of Instances
func (t *RoleAccount) ComputeInstancesLen() int {
	return len(t.ComputeInstances)
}

// ClustersInstancesLen returns length of Clusters
func (t *RoleAccount) ClustersInstancesLen() int {
	return len(t.Clusters)
}
