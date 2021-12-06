package cache

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	aws_sdk "github.com/aws/aws-sdk-go-v2/aws"
	aws_sdk_ec2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	aws_sdk_eks "github.com/aws/aws-sdk-go-v2/service/eks"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	kubernetes_rest "k8s.io/client-go/rest"
	awsiamtoken "sigs.k8s.io/aws-iam-authenticator/pkg/token"

	"github.com/aporeto-se/cloud-operator/common/strbuilder"
)

// Cache this
type Cache struct {
	ec2 *aws_sdk_ec2.Client
	eks *aws_sdk_eks.Client

	Vpcs                 []*Vpc
	Instances            []*Instance
	Clusters             []*Cluster
	RoleAccounts         []*RoleAccount
	awsiamtokenGenerator awsiamtoken.Generator
}

func (t *Cache) init(ctx context.Context) error {

	zap.L().Debug("entering init")

	var err error

	t.awsiamtokenGenerator, err = awsiamtoken.NewGenerator(true, false)
	if err != nil {
		return err
	}

	roleAccountMap := make(map[string]*RoleAccount)
	clusterMap := make(map[string]*Cluster)
	vpcMap := make(map[string]*Vpc)

	// Get AWS VPCs, Clusters and Instances

	vpcList, err := t.ec2.DescribeVpcs(ctx, &aws_sdk_ec2.DescribeVpcsInput{})
	if err != nil {
		zap.L().Debug("returning init with error(s)")
		return err
	}

	awsSubnets, err := t.ec2.DescribeSubnets(ctx, &aws_sdk_ec2.DescribeSubnetsInput{})
	if err != nil {
		zap.L().Debug("returning init with error(s)")
		return err
	}

	clusterList, err := t.eks.ListClusters(ctx, &aws_sdk_eks.ListClustersInput{})
	if err != nil {
		zap.L().Debug("returning init with error(s)")
		return err
	}

	awsInstances, err := t.ec2.DescribeInstances(ctx, &aws_sdk_ec2.DescribeInstancesInput{})
	if err != nil {
		zap.L().Debug("returning init with error(s)")
		return err
	}

	// Iterate AWS VPC and
	// 1: Create new local VPC
	// 2: Store VPC in map using key vpcID
	// 3: Store VPC in slice

	for _, awsVpc := range vpcList.Vpcs {
		vpc := newVpc(&awsVpc)
		vpcMap[*vpc.VpcId] = vpc
		t.Vpcs = append(t.Vpcs, vpc)
	}

	for _, awsSubnet := range awsSubnets.Subnets {
		subnet := newSubnet(&awsSubnet)
		vpc := vpcMap[*subnet.VpcId]
		if vpc == nil {
			return fmt.Errorf("Missing VPC for vpcID %s", *subnet.VpcId)
		}
		vpc.Subnets = append(vpc.Subnets, subnet)
	}

	// Iterate AWS EKS Clusters and
	// 1: Create new local cluster
	// 2: Attach cluster to its VPC and its VPC to the cluster
	// 3: Attach the cluster to its Role Account and its Role Account to the cluster
	// 4: Create Role Account if necessary and store in map using role account name as the key
	// 5: Store cluster in a map using its name as the key
	// 6: Store cluster in a slice

	for _, name := range clusterList.Clusters {

		zap.L().Debug(fmt.Sprintf("Processing cluster %s", name))

		awsCluster, err := t.eks.DescribeCluster(ctx, &aws_sdk_eks.DescribeClusterInput{
			Name: aws_sdk.String(name),
		})
		if err != nil {
			zap.L().Debug("returning init with error(s)")
			return err
		}

		cluster := newCluster(awsCluster.Cluster)

		vpc := vpcMap[*cluster.ResourcesVpcConfig.VpcId]
		if vpc == nil {
			return fmt.Errorf("Missing VPC for vpcID %s", *cluster.ResourcesVpcConfig.VpcId)
		}

		cluster.Vpc = vpc
		vpc.Clusters = append(vpc.Clusters, cluster)

		listNodegroups, err := t.eks.ListNodegroups(ctx, &aws_sdk_eks.ListNodegroupsInput{
			ClusterName: cluster.Name,
		})
		if err != nil {
			zap.L().Debug("returning init with error(s)")
			return err
		}

		for _, nodeGroupName := range listNodegroups.Nodegroups {

			describeNodegroup, err := t.eks.DescribeNodegroup(ctx, &aws_sdk_eks.DescribeNodegroupInput{
				ClusterName:   cluster.Name,
				NodegroupName: aws_sdk.String(nodeGroupName),
			})
			if err != nil {
				zap.L().Debug("returning init with error(s)")
				return err
			}

			roleAccountName := arnToRole(*describeNodegroup.Nodegroup.NodeRole)
			roleAccount := roleAccountMap[roleAccountName]
			if roleAccount == nil {
				roleAccount = newRoleAccount(roleAccountName)
				t.RoleAccounts = append(t.RoleAccounts, roleAccount)

			}

			cluster.RoleAccounts = append(cluster.RoleAccounts, roleAccount)
			roleAccount.Clusters = append(roleAccount.Clusters, cluster)
			roleAccountMap[roleAccountName] = roleAccount

		}

		clusterMap[name] = cluster
		t.Clusters = append(t.Clusters, cluster)

	}

	// Iterate AWS Instances and
	// 1: Attach Instance to its VPC and its VPC to the Instance
	// 2: Attach Instance to its Role Account
	// 3: If Role Account does not exist, create Role Account and store in map using Role Account name as the key
	// 4: Attach Instance to its Cluster and its Cluster to the Instance

	for _, awsReservation := range awsInstances.Reservations {
		for _, awsInstance := range awsReservation.Instances {

			instance := newInstance(&awsInstance)

			vpc := vpcMap[*instance.VpcId]
			if vpc == nil {
				return fmt.Errorf("Missing VPC for vpcID %s", *instance.VpcId)
			}

			instance.Vpc = vpc
			vpc.Instances = append(vpc.Instances, instance)

			// Map instance to its cluster by iterating its tags and looking for the key eks:cluster-name
			// If present the value is the name of the cluster which we have in a map

			var cluster *Cluster

			for _, tag := range instance.Tags {
				if *tag.Key == "eks:cluster-name" {
					cluster = clusterMap[*tag.Value]
					if cluster != nil {
						instance.Cluster = cluster
						cluster.Instances = append(cluster.Instances, instance)
						zap.L().Debug(fmt.Sprintf("instanceID %s is part of cluster %s", *instance.InstanceId, *cluster.Name))
					} else {
						zap.L().Warn(fmt.Sprintf("instanceID %s is missing its cluster", *instance.InstanceId))
					}
					break
				}
			}

			if instance.IamInstanceProfile != nil {

				roleAccountName := arnToRole(*instance.IamInstanceProfile.Arn)
				roleAccount := roleAccountMap[roleAccountName]
				if roleAccount == nil {
					roleAccount = newRoleAccount(roleAccountName)
					t.RoleAccounts = append(t.RoleAccounts, roleAccount)
				}

				instance.RoleAccount = roleAccount
				roleAccountMap[roleAccountName] = roleAccount

				if cluster != nil {
					roleAccount.ClusterInstances = append(roleAccount.ClusterInstances, instance)
				} else {
					roleAccount.ComputeInstances = append(roleAccount.ComputeInstances, instance)
				}

			}

			t.Instances = append(t.Instances, instance)

		}
	}

	zap.L().Debug("returning init")
	return nil
}

// KubeConfig returns Kubernetes Clientset for specified cluster
func (t *Cache) KubeConfig(cluster *Cluster) (*kubernetes.Clientset, error) {

	zap.L().Debug("entering KubeConfig")

	name := *cluster.Name
	cert := *cluster.CertificateAuthority.Data
	endpoint := *cluster.Endpoint

	opts := &awsiamtoken.GetTokenOptions{
		ClusterID: name,
	}

	tok, err := t.awsiamtokenGenerator.GetWithOptions(opts)
	if err != nil {
		zap.L().Debug("returning KubeConfig with error(s)")
		return nil, err
	}
	ca, err := base64.StdEncoding.DecodeString(cert)
	if err != nil {
		zap.L().Debug("returning KubeConfig with error(s)")
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(
		&kubernetes_rest.Config{
			Host:        endpoint,
			BearerToken: tok.Token,
			TLSClientConfig: kubernetes_rest.TLSClientConfig{
				CAData: ca,
			},
		},
	)

	if err != nil {
		zap.L().Debug("returning KubeConfig with error(s)")
		return nil, err
	}

	zap.L().Debug("returning KubeConfig")
	return clientset, nil
}

// ReportString returns cache as a printable string. Useful for debugging.
func (t *Cache) ReportString() string {

	s := strbuilder.NewStrBuilder("\n")

	for _, x := range t.Instances {

		clusterName := ""

		if x.Cluster != nil {
			clusterName = *x.Cluster.Name
		}

		s.A(fmt.Sprintf("Instance %s role=\"%s\" vpc=\"%s\" cluster=\"%s\"",
			*x.InstanceId, x.RoleAccount.Name, *x.VpcId, clusterName))

	}

	for _, x := range t.Clusters {
		s.A(fmt.Sprintf("Cluster %s roles=[%s] vpc=\"%s\" instances=[%s]",
			*x.Name, roleAccountsToCommaString(x.RoleAccounts),
			*x.Vpc.VpcId, instanceIDsToCommaString(x.Instances)))
	}

	for _, x := range t.RoleAccounts {
		s.A(fmt.Sprintf("Role %s clusters=[%s] computeInstances=[%s] clusterInstances=[%s]",
			x.Name, clusterToCommaString(x.Clusters),
			instanceIDsToCommaString(x.ComputeInstances),
			instanceIDsToCommaString(x.ClusterInstances)))
	}

	return s.Build()
}

func roleAccountsToCommaString(input []*RoleAccount) string {
	result := ""
	for _, x := range input {
		if result == "" {
			result = x.Name
		} else {
			result = result + ", " + x.Name
		}
	}
	return result
}

func instanceIDsToCommaString(input []*Instance) string {
	result := ""
	for _, x := range input {
		if result == "" {
			result = *x.InstanceId
		} else {
			result = result + ", " + *x.InstanceId
		}
	}
	return result
}

func clusterToCommaString(input []*Cluster) string {
	result := ""
	for _, x := range input {
		if result == "" {
			result = *x.Name
		} else {
			result = result + ", " + *x.Name
		}
	}
	return result
}

// func vpcToCommaString(input []*Vpc) string {
// 	result := ""
// 	for _, x := range input {
// 		if result == "" {
// 			result = *x.VpcId
// 		} else {
// 			result = result + ", " + *x.VpcId
// 		}
// 	}
// 	return result
// }

func arnToRole(input string) string {
	x := strings.Split(input, "/")
	return x[len(x)-1]
}
