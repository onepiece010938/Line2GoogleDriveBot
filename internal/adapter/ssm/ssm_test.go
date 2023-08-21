package ssm

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type SSMGetParameterImpl struct{}

func (dt SSMGetParameterImpl) GetParameter(ctx context.Context,
	params *ssm.GetParameterInput,
	optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {

	parameter := &types.Parameter{Value: aws.String("parameter-value")}

	output := &ssm.GetParameterOutput{
		Parameter: parameter,
	}

	return output, nil
}

type Config struct {
	ParameterName string `json:"ParameterName"`
}

var configFileName = "config.json"

var globalConfig Config

func populateConfiguration(t *testing.T) error {
	content, err := os.ReadFile(configFileName)
	if err != nil {
		return err
	}

	text := string(content)

	err = json.Unmarshal([]byte(text), &globalConfig)
	if err != nil {
		return err
	}

	if globalConfig.ParameterName == "" {
		msg := "You must supply a value for ParameterName in " + configFileName
		return errors.New(msg)
	}

	return nil
}

func TestFindParameter(t *testing.T) {
	thisTime := time.Now()
	nowString := thisTime.Format("2006-01-02 15:04:05 Monday")
	t.Log("Starting unit test at " + nowString)

	err := populateConfiguration(t)
	if err != nil {
		t.Fatal(err)
	}

	api := &SSMGetParameterImpl{}

	// input := &ssm.GetParameterInput{
	// 	Name: &globalConfig.ParameterName,
	// }

	resp, err := testSSM.FindParameter(context.Background(), *api, globalConfig.ParameterName)
	if err != nil {
		t.Log("Got an error ...:")
		t.Log(err)
		return
	}

	t.Log("Parameter value: " + resp)
}

var testSSM *SSM

func TestMain(m *testing.M) {

	testSSM = NewSSM()

	os.Exit(m.Run())
}
