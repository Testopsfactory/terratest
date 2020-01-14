package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/lambda"
)

// InvokeFunction invokes a lambda function.
func InvokeFunction(t *testing.T, region, functionName string) FunctionResult {
	out, err := InvokeFunctionE(t, region, functionName)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// InvokeFunctionE invokes a lambda function.
func InvokeFunctionE(t *testing.T, region, functionName string) (FunctionResult, error) {
	lambdaClient, err := NewLambdaClientE(t, region)
	if err != nil {
		return FunctionResult{}, err
	}

	invokeInput := &lambda.InvokeInput{
		FunctionName: &functionName,
	}

	out, err := lambdaClient.Invoke(invokeInput)
	if err != nil {
		return FunctionResult{}, err
	}

	result := FunctionResult{
		Payload:    out.Payload,
		StatusCode: *out.StatusCode,
	}

	if out.FunctionError != nil {
		result.FunctionError = *out.FunctionError
	}

	return result, err
}

type FunctionResult struct {
	FunctionError string
	Payload       []byte
	StatusCode    int64
}

// NewLambdaClient creates a new Lambda client.
func NewLambdaClient(t *testing.T, region string) *lambda.Lambda {
	client, err := NewLambdaClientE(t, region)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// NewLambdaClientE creates a new Lambda client.
func NewLambdaClientE(t *testing.T, region string) (*lambda.Lambda, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return lambda.New(sess), nil
}
