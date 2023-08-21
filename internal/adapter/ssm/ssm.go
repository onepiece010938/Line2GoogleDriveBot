package ssm

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// SSMGetParameterAPI defines the interface for the GetParameter function.
// We use this interface to test the function using a mocked service.

type SSMGetParameterAPI interface {
	GetParameter(ctx context.Context,
		params *ssm.GetParameterInput,
		optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
}

func (s *SSM) FindParameter(c context.Context, api SSMGetParameterAPI, name string) (string, error) {
	input := &ssm.GetParameterInput{
		Name: &name,
	}
	results, err := api.GetParameter(c, input)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	fmt.Println(*results.Parameter.Value)

	return *results.Parameter.Value, nil
}

type SSM struct {
	Client *ssm.Client
}

func NewSSM() *SSM {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	client := ssm.NewFromConfig(cfg)
	return &SSM{
		Client: client,
	}
}
