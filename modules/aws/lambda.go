package aws

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// LambdaOptions contains additional parameters for InvokeFunctionWithParams().
// It contains a subset of the fields found in the lambda.InvokeInput struct.
type LambdaOptions struct {
	// FunctionName is a required field containing the lambda function name.
	FunctionName *string

	// InvocationType can be one of "RequestResponse" or "DryRun".
	//    * RequestResponse (default) - Invoke the function synchronously.
	//    Keep the connection open until the function returns a response
	//    or times out.
	//
	//    * DryRun - Validate parameter values and verify that the user or
	//    role has permission to invoke the function.
	InvocationType *string

	// Lambda function input; will be converted to JSON.
	Payload interface{}
}

// LambdaOutput contains the output from InvokeFunctionWithParams().  The
// fields may or may not have a value depending on the invocation type and
// whether an error occurred or not.
type LambdaOutput struct {
	// If present, indicates that an error occurred during function execution.
	// Error details are included in the response payload.
	FunctionError *string

	// The response from the function, or an error object.
	Payload []byte

	// The HTTP status code for a successful request is in the 200 range.
	// For RequestResponse invocation type, the status code is 200.
	// For the DryRun invocation type, the status code is 204.
	StatusCode *int64
}

// InvokeFunction invokes a lambda function.
func InvokeFunction(t testing.TestingT, region, functionName string, payload interface{}) []byte {
	input := &LambdaOptions{
		FunctionName: &functionName,
		Payload:      &payload,
	}
	out, err := InvokeFunctionWithParams(t, region, input)
	require.NoError(t, err)
	return out.Payload
}

// InvokeFunctionE invokes a lambda function.
func InvokeFunctionE(t testing.TestingT, region, functionName string, payload interface{}) ([]byte, error) {
	input := &LambdaOptions{
		FunctionName: &functionName,
		Payload:      &payload,
	}
	out, err := InvokeFunctionWithParams(t, region, input)
	if err != nil {
		return nil, err
	}

	if out.FunctionError != nil {
		return out.Payload, &FunctionError{Message: *out.FunctionError, StatusCode: *out.StatusCode, Payload: out.Payload}
	}

	return out.Payload, nil
}

// InvokeFunctionWithParams invokes a lambda function using parameters
// supplied in the LambdaOptions struct and returns values in a LambdaOutput
// struct.
func InvokeFunctionWithParams(t testing.TestingT, region string, input *LambdaOptions) (*LambdaOutput, error) {
	lambdaClient, err := NewLambdaClientE(t, region)
	if err != nil {
		return nil, err
	}

	// The function name is a required field in LambdaOptions. If missing,
	// report the error.
	if input.FunctionName == nil {
		msg := "LambdaOptions.FunctionName is a required field"
		return &LambdaOutput{FunctionError: &msg}, errors.New(msg)
	}

	// Verify the InvocationType is one of the allowed values and report
	// an error if its not.  By default the InvocationType will be
	// "RequestResponse".
	invocationType := lambda.InvocationTypeRequestResponse
	if input.InvocationType != nil {
		switch *input.InvocationType {
		case
			lambda.InvocationTypeRequestResponse,
			lambda.InvocationTypeDryRun:
			invocationType = *input.InvocationType
		default:
			msg := fmt.Sprintf("LambdaOptions.InvocationType, if specified, must either be \"%s\" or \"%s\"",
				lambda.InvocationTypeRequestResponse,
				lambda.InvocationTypeDryRun)
			return &LambdaOutput{FunctionError: &msg}, errors.New(msg)
		}
	}

	invokeInput := &lambda.InvokeInput{
		FunctionName:   input.FunctionName,
		InvocationType: &invocationType,
	}

	if input.Payload != nil {
		payloadJson, err := json.Marshal(input.Payload)
		if err != nil {
			return nil, err
		}
		invokeInput.Payload = payloadJson
	}

	out, err := lambdaClient.Invoke(invokeInput)

	// As this function supports different invocation types, so it must
	// support different combinations of output.
	lambdaOutput := LambdaOutput{
		FunctionError: out.FunctionError,
		Payload:       out.Payload,
		StatusCode:    out.StatusCode,
	}
	return &lambdaOutput, err
}

type FunctionError struct {
	Message    string
	StatusCode int64
	Payload    []byte
}

func (err *FunctionError) Error() string {
	return fmt.Sprintf("%s error invoking lambda function: %v", err.Message, err.Payload)
}

// NewLambdaClient creates a new Lambda client.
func NewLambdaClient(t testing.TestingT, region string) *lambda.Lambda {
	client, err := NewLambdaClientE(t, region)
	require.NoError(t, err)
	return client
}

// NewLambdaClientE creates a new Lambda client.
func NewLambdaClientE(t testing.TestingT, region string) (*lambda.Lambda, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return lambda.New(sess), nil
}
