package cache

import (
	aws_sdk_ec2_types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// Vpc AWS Vpc
type Vpc struct {
	*aws_sdk_ec2_types.Vpc
	Instances []*Instance
	Clusters  []*Cluster
	Subnets   []*Subnet
}

func newVpc(vpc *aws_sdk_ec2_types.Vpc) *Vpc {
	return &Vpc{
		Vpc: vpc,
	}
}
