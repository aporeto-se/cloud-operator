package types

import (
	"fmt"
	"os"

	"github.com/aporeto-se/cloud-operator/aws/types"
)

// CloudOperatorConfig AWS Implementation
type CloudOperatorConfig struct {
	types.CloudOperatorConfig
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
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

	accessKeyID := os.Getenv(AccessKeyIDEnv)
	secretAccessKey := os.Getenv(SecretAccessKeyEnv)
	sessionToken := os.Getenv(SessionTokenEnv)

	if accessKeyID != "" {
		t.AccessKeyID = accessKeyID
	}

	if secretAccessKey != "" {
		t.SecretAccessKey = secretAccessKey
	}

	if sessionToken != "" {
		t.SessionToken = sessionToken
	}

	return t.CloudOperatorConfig.SetFromEnv()
}

// GetAccessKeyID returns attribute or error
func (t *CloudOperatorConfig) GetAccessKeyID() (string, error) {
	var err error
	if t.AccessKeyID == "" {
		err = fmt.Errorf("attribute Region (env var %s) is required", AccessKeyIDEnv)
	}
	return t.AccessKeyID, err
}

// GetSecretAccessKey returns attribute or error
func (t *CloudOperatorConfig) GetSecretAccessKey() (string, error) {
	var err error
	if t.SecretAccessKey == "" {
		err = fmt.Errorf("attribute Region (env var %s) is required", SecretAccessKeyEnv)
	}
	return t.SecretAccessKey, err
}

// GetSessionToken returns attribute or error
func (t *CloudOperatorConfig) GetSessionToken() (string, error) {
	var err error
	if t.SessionToken == "" {
		err = fmt.Errorf("attribute Region (env var %s) is required", SessionTokenEnv)
	}
	return t.SessionToken, err
}
