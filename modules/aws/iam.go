package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"fmt"
	"github.com/aws/aws-sdk-go/service/iam"
)

// Get an AWS Session, checking that the user has credentials properly configured in their environment
func GetAuthenticatedSession(region string) (*session.Session, error) {
	sess, err := session.NewSession(aws.NewConfig().WithRegion(region))
	if err != nil {
		return nil, err
	}

	if _, err = sess.Config.Credentials.Get(); err != nil {
		return nil, CredentialsError{UnderlyingErr: err}
	}

	return sess, nil
}

// Create a new IAM client
func NewIamClient() (*iam.IAM, error) {
	sess, err := GetAuthenticatedSession(defaultRegion)
	if err != nil {
		return nil, err
	}
	return iam.New(sess), nil
}

type CredentialsError struct {
	UnderlyingErr error
}

func (err CredentialsError) Error() string {
	return fmt.Sprintf("Error finding AWS credentials. Did you set the AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables or configure an AWS profile? Underlying error: %v", err.UnderlyingErr)
}