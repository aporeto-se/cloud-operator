package cache

import (
	"context"
	"fmt"
	"net/http"

	aws_sdk_config "github.com/aws/aws-sdk-go-v2/config"
	aws_sdk_ec2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	aws_sdk_eks "github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap"
)

// Config ...
type Config struct {
	AWSRegion  string
	HTTPClient *http.Client
}

// NewConfig returns a new entity instance
func NewConfig() *Config {
	return &Config{}
}

// SetRegion sets attribute and returns self
func (t *Config) SetRegion(awsRegion string) *Config {
	t.AWSRegion = awsRegion
	return t
}

// SetHTTPClient sets entity and returns self
func (t *Config) SetHTTPClient(httpClient *http.Client) *Config {
	t.HTTPClient = httpClient
	return t
}

// GetHTTPClient returns entity. If entity is nil, entity will be initialized
func (t *Config) GetHTTPClient() *http.Client {
	if t.HTTPClient == nil {
		t.HTTPClient = &http.Client{}
	}
	return t.HTTPClient
}

// Build returns new Cache from config or error
func (t *Config) Build(ctx context.Context) (*Cache, error) {

	zap.L().Debug("entering Build")

	var errors *multierror.Error

	if t.AWSRegion == "" {
		errors = multierror.Append(errors, fmt.Errorf("attribute AWSRegion is required"))
	}

	err := errors.ErrorOrNil()
	if err != nil {
		zap.L().Debug("returning Build with error(s)")
		return nil, err
	}

	awsConfig, err := aws_sdk_config.LoadDefaultConfig(ctx)
	if err != nil {
		zap.L().Debug("returning NewOperator -> error(s)")
		return nil, err
	}

	awsConfig.HTTPClient = t.GetHTTPClient()

	awsConfig.Region = t.AWSRegion

	c := &Cache{
		ec2: aws_sdk_ec2.NewFromConfig(awsConfig),
		eks: aws_sdk_eks.NewFromConfig(awsConfig),
	}

	err = c.init(ctx)
	if err != nil {
		zap.L().Debug("returning Build with error(s)")
		return nil, err
	}

	zap.L().Debug("returning Build")
	return c, nil
}
