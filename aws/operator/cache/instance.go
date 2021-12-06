package cache

import (
	aws_sdk_ec2_types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// Instance AWS Compute Instance
type Instance struct {
	*aws_sdk_ec2_types.Instance
	Vpc         *Vpc
	RoleAccount *RoleAccount
	Cluster     *Cluster
}

func newInstance(instance *aws_sdk_ec2_types.Instance) *Instance {
	return &Instance{
		Instance: instance,
	}
}
