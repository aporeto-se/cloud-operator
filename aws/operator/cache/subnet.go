package cache

import (
	aws_sdk_ec2_types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// Subnet AWS Subnet
type Subnet struct {
	*aws_sdk_ec2_types.Subnet
	Instances []*Instance
	Clusters  []*Cluster
}

func newSubnet(subnet *aws_sdk_ec2_types.Subnet) *Subnet {
	return &Subnet{
		Subnet: subnet,
	}
}
