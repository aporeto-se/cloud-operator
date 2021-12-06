package operator

import (
	"context"
	"net/http"

	prisma_api "github.com/aporeto-se/prisma-sdk-go-v2/api"

	"github.com/aporeto-se/cloud-operator/gcp/types"
)

// Config this config
type Config struct {
	CloudOperatorConfig *types.CloudOperatorConfig
	PrismaClient        *prisma_api.Client
	HTTPClient          *http.Client
}

// NewConfig returns new entity instance
func NewConfig() *Config {
	return &Config{}
}

// SetCloudOperatorConfig sets entity and returns self
func (t *Config) SetCloudOperatorConfig(cloudOperatorConfig *types.CloudOperatorConfig) *Config {
	t.CloudOperatorConfig = cloudOperatorConfig
	return t
}

// SetPrismaClient sets entity and returns self
func (t *Config) SetPrismaClient(prismaClient *prisma_api.Client) *Config {
	t.PrismaClient = prismaClient
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

// Build returns entity or error
func (t *Config) Build(ctx context.Context) (*Client, error) {
	return NewClient(ctx, t)
}
