package operator

// import (
// 	"context"

// 	"github.com/hashicorp/go-multierror"

// 	api_client "github.com/aporeto-se/prisma-sdk-go-v2/api"
// 	token "github.com/aporeto-se/prisma-sdk-go-v2/token/gcp"
// )

// // NewClient returns instance from env
// func NewClient(ctx context.Context) (*Client, error) {

// 	config := NewConfig()

// 	err := config.SetFromEnv()
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = InitLogging(config.PrismaConfig.GetLogLevel())
// 	if err != nil {
// 		return nil, err
// 	}

// 	var errors *multierror.Error

// 	api, err := config.PrismaConfig.GetAPI()
// 	if err != nil {
// 		errors = multierror.Append(errors, err)
// 	}

// 	namespace, err := config.PrismaConfig.GetNamespace()
// 	if err != nil {
// 		errors = multierror.Append(errors, err)
// 	}

// 	err = errors.ErrorOrNil()

// 	if err != nil {
// 		return nil, err
// 	}

// 	tokenProvider, err := token.NewConfig().
// 		SetAPI(api).
// 		SetNamespace(namespace).
// 		SetHTTPClient(config.GetHTTPClient()).
// 		Build(ctx)

// 	if err != nil {
// 		return nil, err
// 	}

// 	// We get our token now so that if there are any errors we can handle them now
// 	_, err = tokenProvider.GetToken(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	prismaClient, err := api_client.NewConfig().
// 		SetNamespace(namespace).
// 		SetTokenProvider(tokenProvider).
// 		SetHTTPClient(config.GetHTTPClient()).
// 		Build(ctx)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return config.SetPrismaClient(prismaClient).Build(ctx)
// }
