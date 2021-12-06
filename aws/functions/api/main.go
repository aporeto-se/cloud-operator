package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	helper "github.com/aporeto-se/cloud-operator/aws/functions"
	operator_types "github.com/aporeto-se/cloud-operator/common/types"
)

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {

	if req.RequestContext.HTTP.Method != "POST" {
		return returnHelperErr()
	}

	operator, err := helper.NewClient(ctx)
	if err != nil {
		return returnError(err)
	}

	var filter *operator_types.Filter

	err = json.Unmarshal([]byte(req.Body), &filter)
	if err != nil {
		return returnError(err)
	}

	report := operator.Run(ctx, nil)
	err = report.Errors()

	body, _ := json.Marshal(report)

	return &events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, err

}

func getExampleConfig() string {
	// s, _ := json.Marshal(operator_types.NewExampleFilterMatchNames())
	// s, _ := json.Marshal(operator_types.NewExampleFilterMatchAny())
	s, _ := json.Marshal(operator_types.NewExampleFilterMatchAny())
	return string(s)
}

func returnError(err error) (*events.APIGatewayProxyResponse, error) {

	newErr := operator_types.NewAPIError(err)
	body := errToJSON(newErr)

	return &events.APIGatewayProxyResponse{
		Body:       body,
		StatusCode: newErr.StatusCode,
	}, err
}

func returnHelperErr() (*events.APIGatewayProxyResponse, error) {

	err := fmt.Errorf("Here is an example: " + getExampleConfig())
	body := errToJSON(err)

	return &events.APIGatewayProxyResponse{
		Body:       body,
		StatusCode: 401,
	}, nil
}

func errToJSON(err error) string {
	jsonBytes, _ := json.Marshal(err)
	return string(jsonBytes)
}

func main() {
	lambda.Start(handler)
}
