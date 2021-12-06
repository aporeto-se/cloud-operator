package types

import (
	"github.com/aporeto-se/cloud-operator/gcp/types"
)

// CloudOperatorConfig AWS Implementation
type CloudOperatorConfig struct {
	types.CloudOperatorConfig
}

// NewCloudOperatorConfig returns new instance of CloudOperatorConfig
func NewCloudOperatorConfig() *CloudOperatorConfig {
	return &CloudOperatorConfig{}
}

// SetFromEnv sets attributes and types from env variables as defined in
// constants file. If attribute is not of the expected type an error will be
// returned. If child entities exist and are initialized (not nil) then a call
// to the child entities SetFromEnv() will be executed. Any errors will be aggregated
// and returned.
func (t *CloudOperatorConfig) SetFromEnv() error {

	return t.CloudOperatorConfig.SetFromEnv()
}
