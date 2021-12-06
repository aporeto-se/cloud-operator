package types

import (
	"fmt"
	"os"

	"github.com/aporeto-se/cloud-operator/common/types"
)

// ================================================================================================

// CloudOperatorConfig GCP Config
type CloudOperatorConfig struct {

	// Parent
	types.CloudOperatorConfig

	// Google Cloud Project ID
	GCloudProject string

	// Google Cloud Zone
	GCloudZone string
}

// SetFromEnv sets attributes and types from env variables as defined in
// constants file. If attribute is not of the expected type an error will be
// returned. If child entities exist and are initialized (not nil) then a call
// to the child entities SetFromEnv() will be executed. Any errors will be aggregated
// and returned.
func (t *CloudOperatorConfig) SetFromEnv() error {

	gCloudProject := os.Getenv(GCloudProjectEnv)
	gCloudZone := os.Getenv(GCloudZoneEnv)

	if gCloudProject != "" {
		t.GCloudProject = gCloudProject
	}

	if gCloudZone != "" {
		t.GCloudZone = gCloudZone
	}

	return t.CloudOperatorConfig.SetFromEnv()
}

// SetGCloudProject sets attribute and returns self
func (t *CloudOperatorConfig) SetGCloudProject(gCloudProject string) *CloudOperatorConfig {
	t.GCloudProject = gCloudProject
	return t
}

// GetGCloudProject returns attribute or error
func (t *CloudOperatorConfig) GetGCloudProject() (string, error) {
	var err error
	if t.GCloudProject == "" {
		err = fmt.Errorf("attribute GCloudProject (env var %s) is required", GCloudProjectEnv)
	}
	return t.GCloudProject, err
}

// SetGCloudZone sets attribute and returns self
func (t *CloudOperatorConfig) SetGCloudZone(gCloudZone string) *CloudOperatorConfig {
	t.GCloudZone = gCloudZone
	return t
}

// GetGCloudZone returns attribute or error
func (t *CloudOperatorConfig) GetGCloudZone() (string, error) {
	var err error
	if t.GCloudZone == "" {
		err = fmt.Errorf("attribute GCloudZone (env var %s) is required", GCloudZoneEnv)
	}
	return t.GCloudZone, err
}
