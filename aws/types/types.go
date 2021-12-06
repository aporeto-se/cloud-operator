package types

import (
	"fmt"
	"os"

	"github.com/aporeto-se/cloud-operator/common/types"
)

// ================================================================================================

// CloudOperatorConfig config for AWS
type CloudOperatorConfig struct {

	// Parent
	types.CloudOperatorConfig

	// AWS Region
	AWSRegion string
}

// NewCloudOperatorConfig returns new intance of entity
func NewCloudOperatorConfig() *CloudOperatorConfig {
	return &CloudOperatorConfig{}
}

// SetFromEnv sets attributes and types from env variables as defined in
// constants file. If attribute is not of the expected type an error will be
// returned. If child entities exist and are initialized (not nil) then a call
// to the child entities SetFromEnv() will be executed. Any errors will be aggregated
// and returned.
func (t *CloudOperatorConfig) SetFromEnv() error {

	awsRegion := os.Getenv(AWSRegionEnv)

	if awsRegion != "" {
		t.AWSRegion = awsRegion
	}

	return t.CloudOperatorConfig.SetFromEnv()
}

// SetAWSRegion sets attribute and returns self
func (t *CloudOperatorConfig) SetAWSRegion(awsRegion string) *CloudOperatorConfig {
	t.AWSRegion = awsRegion
	return t
}

// GetAWSRegion returns attribute or error
func (t *CloudOperatorConfig) GetAWSRegion() (string, error) {
	var err error
	if t.AWSRegion == "" {
		err = fmt.Errorf("attribute AWSRegion (env var %s) is required", AWSRegionEnv)
	}
	return t.AWSRegion, err
}
