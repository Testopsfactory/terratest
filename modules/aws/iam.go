package aws

import (
	"github.com/aws/aws-sdk-go/service/iam"
	"testing"
)

// Get the username fo the current IAM user
func GetIamCurrentUserName(t *testing.T) string {
	out, err := GetIamCurrentUserNameE(t)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// Get the username fo the current IAM user
func GetIamCurrentUserNameE(t *testing.T) (string, error) {
	iamClient, err := NewIamClient(defaultRegion)
	if err != nil {
		return "", err
	}

	resp, err := iamClient.GetUser(&iam.GetUserInput{})
	if err != nil {
		return "", err
	}

	return *resp.User.UserName, nil
}

// Get the ARN for the current IAM user
func GetIamCurrentUserArn(t *testing.T) string {
	out, err := GetIamCurrentUserArnE(t)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// Get the ARN for the current IAM user
func GetIamCurrentUserArnE(t *testing.T) (string, error) {
	iamClient, err := NewIamClient(defaultRegion)
	if err != nil {
		return "", err
	}

	resp, err := iamClient.GetUser(&iam.GetUserInput{})
	if err != nil {
		return "", err
	}

	return *resp.User.Arn, nil
}

// Create a new IAM client
func NewIamClient(region string) (*iam.IAM, error) {
	sess, err := GetAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return iam.New(sess), nil
}
