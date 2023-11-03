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
	"github.com/stretchr/testify/assert"
)

type SSMGetParameterImpl struct{}

func (dt SSMGetParameterImpl) GetParameter(ctx context.Context,
	params *ssm.GetParameterInput,
	optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {

	var parameter *types.Parameter

	if *params.Name == "secret-name" {
		parameter = &types.Parameter{Value: aws.String("secret-value")}
	}
	if *params.Name == "token-name" {
		parameter = &types.Parameter{Value: aws.String("token-value")}
	}

	output := &ssm.GetParameterOutput{
		Parameter: parameter,
	}

	return output, nil
}

func (dt SSMGetParameterImpl) GetParameters(ctx context.Context,
	params *ssm.GetParametersInput,
	optFns ...func(*ssm.Options)) (*ssm.GetParametersOutput, error) {
	var parameters []types.Parameter

	for _, name := range params.Names {

		if name == "secret-name" {

			parameter := &types.Parameter{
				Name:  aws.String(name),
				Value: aws.String("secret-value"),
			}
			parameters = append(parameters, *parameter)
		}
		if name == "token-name" {
			parameter := &types.Parameter{
				Name:  aws.String(name),
				Value: aws.String("token-value"),
			}
			parameters = append(parameters, *parameter)
		}

	}

	output := &ssm.GetParametersOutput{
		Parameters: parameters,
	}

	return output, nil
}

type Config struct {
	MockChannelSecret      string `json:"MockChannelSecret"`
	MockChannelAccessToken string `json:"MockChannelAccessToken"`
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

	if globalConfig.MockChannelSecret == "" {
		msg := "You must supply a value for MockChannelSecret in " + configFileName
		return errors.New(msg)
	}
	if globalConfig.MockChannelAccessToken == "" {
		msg := "You must supply a value for MockChannelAccessToken in " + configFileName
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

	respSecret, err := testSSM.FindParameter(context.Background(), *api, globalConfig.MockChannelSecret)
	assert.NoError(t, err, "Error fetching MockChannelSecret")

	t.Log("MockChannelSecret value: " + respSecret)
	assert.Equal(t, "secret-value", respSecret, "Unexpected value for MockChannelSecret")

	respToken, err := testSSM.FindParameter(context.Background(), *api, globalConfig.MockChannelAccessToken)
	assert.NoError(t, err, "Error fetching MockChannelAccessToken")

	t.Log("MockChannelAccessToken value: " + respToken)
	assert.Equal(t, "token-value", respToken, "Unexpected value for MockChannelAccessToken")
}

func TestFindParameters(t *testing.T) {
	thisTime := time.Now()
	nowString := thisTime.Format("2006-01-02 15:04:05 Monday")
	t.Log("Starting unit test at " + nowString)

	// Creating an instance of the SSMGetParameterImpl
	api := &SSMGetParameterImpl{}

	// Mocking the configuration
	err := populateConfiguration(t)
	assert.NoError(t, err, "Error populating configuration")

	// Setting up the testify assertion instance

	// List of parameter names to fetch
	paramNames := []string{
		globalConfig.MockChannelSecret,
		globalConfig.MockChannelAccessToken,
	}

	// Testing FindParameters
	respParameters, err := testSSM.FindParameters(context.Background(), *api, paramNames)
	assert.NoError(t, err, "Error fetching parameters")

	// Asserting on the response
	expectedValues := map[string]string{
		globalConfig.MockChannelSecret:      "secret-value",
		globalConfig.MockChannelAccessToken: "token-value",
	}

	for paramName, expectedValue := range expectedValues {
		assert.Equal(t, expectedValue, respParameters[paramName], "Unexpected value for parameter: "+paramName)
	}

}

var testSSM *SSM

func TestMain(m *testing.M) {

	testSSM = NewSSM()

	os.Exit(m.Run())
}
