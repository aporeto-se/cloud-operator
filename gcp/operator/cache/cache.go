package cache

import (
	"context"
	"strings"

	"go.uber.org/zap"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" // register GCP auth provider

	gcp_compute "google.golang.org/api/compute/v1"
	gke_service "google.golang.org/api/container/v1"
)

// Cache server
type Cache struct {
	project         string
	zone            string
	Instances       []*Instance
	Clusters        []*Cluster
	ServiceAccounts []*ServiceAccount
}

func (t *Cache) init(ctx context.Context) error {

	zap.L().Debug("entering init")

	var instances []*Instance
	var clusters []*Cluster
	var serviceAccounts []*ServiceAccount

	gcp, err := gcp_compute.NewService(ctx)
	if err != nil {
		zap.L().Debug("returning init with error(s)")
		return err
	}

	gke, err := gke_service.NewService(ctx)
	if err != nil {
		zap.L().Debug("returning init with error(s)")
		return err
	}

	gcpInstances, err := gcp.Instances.List(t.project, t.zone).Context(ctx).Do()
	if err != nil {
		zap.L().Debug("returning init with error(s)")
		return err
	}

	gcpClusters, err := gke.Projects.Zones.Clusters.List(t.project, t.zone).Context(ctx).Do()
	if err != nil {
		zap.L().Debug("returning init with error(s)")
		return err
	}

	createdByToInstanceMap := make(map[string]*Instance)
	emailToServiceAccountMap := make(map[string]*ServiceAccount)

	for _, gcpInstance := range gcpInstances.Items {

		// A new instance is always created as each iteration will be a new instance
		instance := newInstance(gcpInstance)

		_, isKubernetes := instance.Labels["goog-gke-node"]

		// // An instance can either be part of a Kubernetes cluster or not
		// if isKubernetes {
		// 	instance.IsKubernetes = true
		// } else {
		// 	instance.IsInstance = true
		// }

		// If the instance is part of Kubernetes cluster it will have a metadata tag called created-by.
		// We are checking for this tag on each instance but its probably just going to be on instances
		// that are part of a Kubernetes cluster.
		for _, item := range gcpInstance.Metadata.Items {

			// We iterate the tags and if we find one we are looking for we drop out of the loop
			if item.Key == "created-by" {
				instance.CreatedBy = basename(*item.Value)
				break
			}
		}

		// Now we interate the instance Service Accounts. Instances only have a single Service Account but
		// for whatever reason the GCP API/SDK allows for many. We handle it anyways.
		for _, gcpServiceAccount := range gcpInstance.ServiceAccounts {

			// Because Service Accounts can be assigned to multiple instances the
			// service account may already exist. Hence we use a map and key on the
			// service account email which is unique. If the service account exist we
			// fetch it and update it. If not we create it new and store it in the map.

			serviceAccount := emailToServiceAccountMap[gcpServiceAccount.Email]

			if serviceAccount == nil {
				serviceAccount = newServiceAccount(gcpServiceAccount)
			}

			// A Service Account may be assigned to instances that are both part of a Kubernetes cluster
			// and normal instances. So we add the Instance to either the Kuberbetes Instances or the Compute
			// Instances. We only do one per iteration.
			if isKubernetes {
				serviceAccount.KubernetesInstances = append(serviceAccount.KubernetesInstances, instance)
			} else {
				serviceAccount.ComputeInstances = append(serviceAccount.ComputeInstances, instance)
			}

			// We add the instance to the service account, the service account to the instance, the
			// service account to the service account slice and service account map

			serviceAccounts = append(serviceAccounts, serviceAccount)
			instance.ServiceAccounts = append(instance.ServiceAccounts, serviceAccount)
			emailToServiceAccountMap[gcpServiceAccount.Email] = serviceAccount
		}

		// We store the instance in the instances slice and if the instance has the CreatedBy attribute set
		// then we store it in a map. This will be used when we iterate the clusters to map the cluster to its
		// instance(s) and the reverse
		instances = append(instances, instance)

		if instance.CreatedBy != "" {
			createdByToInstanceMap[instance.CreatedBy] = instance
		}

	}

	for _, gkeCluster := range gcpClusters.Clusters {

		// Each iteration is a unique cluster so we create a new cluster wrapper
		cluster := newCluster(gkeCluster)

		// We use the cluster InstanceGroups to find the cluster's instances (if any) by looking the instance up
		// with the CreatedBy attribute
		for _, instanceGroupURL := range cluster.InstanceGroupUrls {
			cluster.CreatedBy = basename(instanceGroupURL)
			instance := createdByToInstanceMap[cluster.CreatedBy]
			if instance != nil {
				// If the instance exist then we add the instance to the cluster and the cluster to the instance
				cluster.Instances = append(cluster.Instances, instance)
				instance.Clusters = append(instance.Clusters, cluster)

				if instance.ServiceAccounts != nil {
					cluster.ServiceAccounts = append(cluster.ServiceAccounts, instance.ServiceAccounts...)
				}

			}
		}

		clusters = append(clusters, cluster)
	}

	zap.L().Debug("returning init")

	t.Instances = instances
	t.Clusters = clusters
	t.ServiceAccounts = serviceAccounts

	return nil
}

func basename(input string) string {
	x := strings.Split(input, "/")
	return x[len(x)-1]
}
