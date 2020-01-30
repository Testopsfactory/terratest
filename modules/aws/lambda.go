package aws

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go/service/lambda"
)

// InvokeFunction invokes a lambda function.
func InvokeFunction(t *testing.T, region, functionName string, payload interface{}) []byte {
	out, err := InvokeFunctionE(t, region, functionName, payload)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// InvokeFunctionE invokes a lambda function.
func InvokeFunctionE(t *testing.T, region, functionName string, payload interface{}) ([]byte, error) {
	lambdaClient, err := NewLambdaClientE(t, region)
	if err != nil {
		return nil, err
	}

	invokeInput := &lambda.InvokeInput{
		FunctionName: &functionName,
	}

	if payload != nil {
		payloadJson, err := json.Marshal(payload)

		if err != nil {
			return nil, err
		} else {
			invokeInput.Payload = payloadJson
		}
	}

	out, err := lambdaClient.Invoke(invokeInput)
	if err != nil {
		return nil, err
	}

	if out.FunctionError != nil {
		return out.Payload, &FunctionError{Message: *out.FunctionError, StatusCode: *out.StatusCode}
	}

	return out.Payload, err
}

type FunctionError struct {
	Message    string
	StatusCode int64
}

func (err *FunctionError) Error() string {
	return err.Message
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
