package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"

	helper "github.com/aporeto-se/cloud-operator/aws/functions"
	operator_types "github.com/aporeto-se/cloud-operator/common/types"
)

func handler(ctx context.Context) (*operator_types.Report, error) {

	operator, err := helper.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	report := operator.Run(ctx, nil)

	return report, report.Errors()
}

func main() {
	lambda.Start(handler)
}
