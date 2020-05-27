package aws

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetParameter retrieves the latest version of SSM Parameter at keyName with decryption.
func GetParameter(t testing.TestingT, awsRegion string, keyName string) string {
	keyValue, err := GetParameterE(t, awsRegion, keyName)
	require.NoError(t, err)
	return keyValue
}

// GetParameterE retrieves the latest version of SSM Parameter at keyName with decryption.
func GetParameterE(t testing.TestingT, awsRegion string, keyName string) (string, error) {
	ssmClient, err := NewSsmClientE(t, awsRegion)
	if err != nil {
		return "", err
	}

	resp, err := ssmClient.GetParameter(&ssm.GetParameterInput{Name: aws.String(keyName), WithDecryption: aws.Bool(true)})
	if err != nil {
		return "", err
	}

	parameter := *resp.Parameter
	return *parameter.Value, nil
}

// PutParameter creates new version of SSM Parameter at keyName with keyValue as SecureString.
func PutParameter(t testing.TestingT, awsRegion string, keyName string, keyDescription string, keyValue string) int64 {
	version, err := PutParameterE(t, awsRegion, keyName, keyDescription, keyValue)
	require.NoError(t, err)
	return version
}

// PutParameterE creates new version of SSM Parameter at keyName with keyValue as SecureString.
func PutParameterE(t testing.TestingT, awsRegion string, keyName string, keyDescription string, keyValue string) (int64, error) {
	ssmClient, err := NewSsmClientE(t, awsRegion)
	if err != nil {
		return 0, err
	}

	resp, err := ssmClient.PutParameter(&ssm.PutParameterInput{Name: aws.String(keyName), Description: aws.String(keyDescription), Value: aws.String(keyValue), Type: aws.String("SecureString")})
	if err != nil {
		return 0, err
	}

	return *resp.Version, nil
}

// NewSsmClient creates a SSM client.
func NewSsmClient(t testing.TestingT, region string) *ssm.SSM {
	client, err := NewSsmClientE(t, region)
	require.NoError(t, err)
	return client
}

// NewSsmClientE creates an SSM client.
func NewSsmClientE(t testing.TestingT, region string) (*ssm.SSM, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return ssm.New(sess), nil
}

// WaitForSsmInstanceE waits until the instance get registered to the SSM inventory.
func WaitForSsmInstanceE(t testing.TestingT, awsRegion, instanceID string, timeout time.Duration) error {
	timeBetweenRetries := 2 * time.Second
	maxRetries := int(timeout.Seconds() / timeBetweenRetries.Seconds())
	description := fmt.Sprintf("Waiting for %s to appear in the SSM inventory", instanceID)

	input := &ssm.GetInventoryInput{
		Filters: []*ssm.InventoryFilter{
			{
				Key:    aws.String("AWS:InstanceInformation.InstanceId"),
				Type:   aws.String("Equal"),
				Values: aws.StringSlice([]string{instanceID}),
			},
		},
	}
	_, err := retry.DoWithRetryE(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		client := NewSsmClient(t, awsRegion)
		req, resp := client.GetInventoryRequest(input)

		if err := req.Send(); err != nil {
			return "", err
		}

		if len(resp.Entities) != 1 {
			return "", fmt.Errorf("%s is not in the SSM inventory", instanceID)
		}

		return "", nil
	})

	return err
}

// WaitForSsmInstance waits until the instance get registered to the SSM inventory.
func WaitForSsmInstance(t testing.TestingT, awsRegion, instanceID string, timeout time.Duration) {
	err := WaitForSsmInstanceE(t, awsRegion, instanceID, timeout)
	require.NoError(t, err)
}

// CheckSsmCommand checks that you can run the given command on the given instance through AWS SSM. Returns stdout and stderr.
func CheckSsmCommand(t testing.TestingT, awsRegion, instanceID, command string, timeout time.Duration) (string, string) {
	stdout, stderr, err := CheckSsmCommandE(t, awsRegion, instanceID, command, timeout)
	require.NoErrorf(t, err, "failed to execute '%s' on %s (%v):]\n  stdout: %#v\n  stderr: %#v", command, instanceID, err, stdout, stderr)
	return stdout, stderr
}

type result struct {
	Stdout string
	Stderr string
}

// CheckSsmCommandE checks that you can run the given command on the given instance through AWS SSM. Returns the stdout, stderr and an error if one occurs.
func CheckSsmCommandE(t testing.TestingT, awsRegion, instanceID, command string, timeout time.Duration) (string, string, error) {
	logger.Logf(t, "Running command '%s' on EC2 instance with ID '%s'", command, instanceID)

	timeBetweenRetries := 2 * time.Second
	maxRetries := int(timeout.Seconds() / timeBetweenRetries.Seconds())

	// Now that we know the instance in the SSM inventory, we can send the command
	client, err := NewSsmClientE(t, awsRegion)
	if err != nil {
		return "", "", err
	}
	comment := "Terratest SSM"
	documentName := "AWS-RunShellScript"
	req, resp := client.SendCommandRequest(&ssm.SendCommandInput{
		Comment:      &comment,
		DocumentName: &documentName,
		InstanceIds:  []*string{&instanceID},
		Parameters: map[string][]*string{
			"commands": []*string{&command},
		},
	})
	if err := req.Send(); err != nil {
		return "", "", err
	}

	// Wait for the result
	description := "Waiting for the result of the command"
	retryableErrors := map[string]string{
		"InvocationDoesNotExist": "InvocationDoesNotExist",
		"bad status: Pending":    "bad status: Pending",
		"bad status: InProgress": "bad status: InProgress",
		"bad status: Delayed":    "bad status: Delayed",
	}
	var stdout, stderr string
	_, err = retry.DoWithRetryableErrorsE(t, description, retryableErrors, maxRetries, timeBetweenRetries, func() (string, error) {
		req, resp := client.GetCommandInvocationRequest(&ssm.GetCommandInvocationInput{
			CommandId:  resp.Command.CommandId,
			InstanceId: &instanceID,
		})
		if err := req.Send(); err != nil {
			return "", err
		}

		// Remove the SSM prefix from stderr
		stderr = aws.StringValue(resp.StandardErrorContent)

		stdout = aws.StringValue(resp.StandardOutputContent)

		status := aws.StringValue(resp.Status)

		if status == "Success" {
			return "", nil
		}

		if status == "Failed" {
			return "", fmt.Errorf("Failed")
		}

		return "", fmt.Errorf("bad status: %s", status)
	})

	if err != nil {
		return stdout, stderr, err.(retry.FatalError).Underlying
	}

	return stdout, stderr, nil
}
