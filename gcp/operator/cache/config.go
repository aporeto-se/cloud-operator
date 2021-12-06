package cache

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap"
)

// Config ...
type Config struct {
	Project string
	Zone    string
}

// NewConfig returns a new entity instance
func NewConfig() *Config {
	return &Config{}
}

// SetProject sets attribute and returns self
func (t *Config) SetProject(project string) *Config {
	t.Project = project
	return t
}

// SetZone sets attribute and returns self
func (t *Config) SetZone(zone string) *Config {
	t.Zone = zone
	return t
}

// Build returns new Cache from config or error
func (t *Config) Build(ctx context.Context) (*Cache, error) {

	zap.L().Debug("entering Build")

	var errors *multierror.Error

	project := t.Project
	zone := t.Zone

	if project == "" {
		zap.L().Debug("returning Build with error(s)")
		errors = multierror.Append(errors, fmt.Errorf("project is required"))
	}

	if zone == "" {
		zap.L().Debug("returning Build with error(s)")
		errors = multierror.Append(errors, fmt.Errorf("zone is required"))
	}

	err := errors.ErrorOrNil()
	if err != nil {
		return nil, err
	}

	c := &Cache{
		project: project,
		zone:    zone,
	}

	err = c.init(ctx)
	if err != nil {
		zap.L().Debug("returning Build with error(s)")
		return nil, err
	}

	zap.L().Debug("returning Build")
	return c, nil
}
