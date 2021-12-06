package helper

import (
	"context"
	"net/http"

	prisma_api "github.com/aporeto-se/prisma-sdk-go-v2/api"
	token "github.com/aporeto-se/prisma-sdk-go-v2/token/aws"
	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap"

	"github.com/aporeto-se/cloud-operator/aws/functions/types"
	operator "github.com/aporeto-se/cloud-operator/aws/operator"
)

// NewClient returns new Client
func NewClient(ctx context.Context) (*operator.Client, error) {

	// Logging has NOT been initialized yet

	cloudOperatorConfig := types.NewCloudOperatorConfig()
	err := cloudOperatorConfig.SetFromEnv()
	if err != nil {
		return nil, err
	}

	err = operator.InitLogging(cloudOperatorConfig.GetLogLevel())
	if err != nil {
		return nil, err
	}

	// Logging is now initialized

	zap.L().Debug("logging initialized : inside NewOperator")

	var errors *multierror.Error

	api, err := cloudOperatorConfig.GetAPI()
	if err != nil {
		errors = multierror.Append(errors, err)
	}

	namespace, err := cloudOperatorConfig.GetNamespace()
	if err != nil {
		errors = multierror.Append(errors, err)
	}

	accessKeyID, err := cloudOperatorConfig.GetAccessKeyID()
	if err != nil {
		errors = multierror.Append(errors, err)
	}

	secretAccessKey, err := cloudOperatorConfig.GetSecretAccessKey()
	if err != nil {
		errors = multierror.Append(errors, err)
	}

	sessionToken, err := cloudOperatorConfig.GetSessionToken()
	if err != nil {
		errors = multierror.Append(errors, err)
	}

	err = errors.ErrorOrNil()

	if err != nil {
		zap.L().Debug("returning NewOperator with error(s)")
		return nil, err
	}

	httpClient := &http.Client{}

	tokenprovider, err := token.NewConfig().
		SetAPI(api).
		SetNamespace(namespace).
		SetHTTPClient(httpClient).
		SetAccessKeyID(accessKeyID).
		SetSecretAccessKey(secretAccessKey).
		SetSessionToken(sessionToken).Build()
	if err != nil {
		zap.L().Debug("returning NewOperator with error(s)")
		return nil, err
	}

	prismaClient, err := prisma_api.NewConfig().
		SetNamespace(namespace).
		SetAPI(api).
		SetTokenProvider(tokenprovider).
		SetHTTPClient(httpClient).Build(ctx)
	if err != nil {
		zap.L().Debug("returning NewOperator with error(s)")
		return nil, err
	}

	operator, err := operator.NewConfig().
		SetCloudOperatorConfig(&cloudOperatorConfig.CloudOperatorConfig).
		SetPrismaClient(prismaClient).SetHTTPClient(httpClient).
		Build(ctx)
	if err != nil {
		zap.L().Debug("returning NewOperator with error(s)")
		return nil, err
	}

	zap.L().Debug("returning NewOperator")
	return operator, nil
}
