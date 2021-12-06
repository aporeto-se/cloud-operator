package cache

import (
	"strings"

	gcp_compute "google.golang.org/api/compute/v1"
)

// ServiceAccount ...
type ServiceAccount struct {
	*gcp_compute.ServiceAccount
	NamespaceName       string
	ComputeInstances    []*Instance
	KubernetesInstances []*Instance
}

func newServiceAccount(serviceAccount *gcp_compute.ServiceAccount) *ServiceAccount {

	namespaceName := strings.Split(serviceAccount.Email, "@")[0]

	return &ServiceAccount{
		ServiceAccount: serviceAccount,
		NamespaceName:  namespaceName,
	}
}

// ComputeInstancesLen returns length of ComputeInstances
func (t *ServiceAccount) ComputeInstancesLen() int {
	return len(t.ComputeInstances)
}

// ClustersLen returns length of KubernetesInstances
func (t *ServiceAccount) ClustersLen() int {
	return len(t.KubernetesInstances)
}
