package cache

import (
	aws_sdk_eks_types "github.com/aws/aws-sdk-go-v2/service/eks/types"
)

// Cluster AWS Kubernetes Cluster
type Cluster struct {
	*aws_sdk_eks_types.Cluster
	Vpc          *Vpc
	Instances    []*Instance
	RoleAccounts []*RoleAccount
}

func newCluster(cluster *aws_sdk_eks_types.Cluster) *Cluster {
	return &Cluster{
		Cluster: cluster,
	}
}

// ComputeInstancesLen returns length of Instances
func (t *Cluster) ComputeInstancesLen() int {
	return len(t.Instances)
}

// RoleAccountLen returns length of RoleAccounts
func (t *Cluster) RoleAccountLen() int {
	return len(t.RoleAccounts)
}
