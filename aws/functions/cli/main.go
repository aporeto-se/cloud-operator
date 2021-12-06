package main

import (
	"context"
	"encoding/json"
	"fmt"

	helper "github.com/aporeto-se/cloud-operator/aws/functions"
)

func main() {

	ctx := context.Background()

	err := run(ctx)
	if err != nil {
		panic(err)
	}

}

func run(ctx context.Context) error {

	operator, err := helper.NewClient(ctx)
	if err != nil {
		return err
	}

	report := operator.Run(ctx, nil)
	jsonReport, _ := json.Marshal(report)
	fmt.Println(string(jsonReport))

	return err
}
