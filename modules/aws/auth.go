package aws

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/pquerna/otp/totp"
)

// Get an AWS Session, checking that the user has credentials properly configured in their environment
func NewAuthenticatedSession(region string) (*session.Session, error) {
	sess, err := session.NewSession(aws.NewConfig().WithRegion(region))
	if err != nil {
		return nil, err
	}

	if _, err = sess.Config.Credentials.Get(); err != nil {
		return nil, CredentialsError{UnderlyingErr: err}
	}

	return sess, nil
}

// CreateAwsSessionFromRole returns a new AWS session after assuming the role
// whose ARN is provided in roleARN.
func CreateAwsSessionFromRole(region string, roleARN string) (*session.Session, error) {
	sess, err := session.NewSession(aws.NewConfig().WithRegion(region))
	if err != nil {
		return nil, err
	}
	sess, err = AssumeRole(sess, roleARN)
	return sess, err
}

// AssumeRole mutates the provided session by obtaining new credentials by
// assuming the role provided in roleARN.
func AssumeRole(sess *session.Session, roleARN string) (*session.Session, error) {
	sess.Config.Credentials = stscreds.NewCredentials(sess, roleARN)
	return sess, nil
}

// Create a new AWS session using explicit credentials. This is useful if you want to create an IAM User dynamically and
// create an AWS session authenticated as the new IAM User.
func CreateAwsSessionWithCreds(region string, accessKeyId string, secretAccessKey string) (*session.Session, error) {
	creds := CreateAwsCredentials(accessKeyId, secretAccessKey)
	return session.NewSession(aws.NewConfig().WithRegion(region).WithCredentials(creds))
}

// Create a new AWS session authenticated using an MFA token retrieved using the given STS client and MFA Device
func CreateAwsSessionWithMfa(region string, stsClient *sts.STS, mfaDevice *iam.VirtualMFADevice) (*session.Session, error) {
	tokenCode, err := GetTimeBasedOneTimePassword(mfaDevice)
	if err != nil {
		return nil, err
	}

	output, err := stsClient.GetSessionToken(&sts.GetSessionTokenInput{
		SerialNumber: mfaDevice.SerialNumber,
		TokenCode:    aws.String(tokenCode),
	})
	if err != nil {
		return nil, err
	}

	accessKeyId := *output.Credentials.AccessKeyId
	secretAccessKey := *output.Credentials.SecretAccessKey
	sessionToken := *output.Credentials.SessionToken

	creds := CreateAwsCredentialsWithSessionToken(accessKeyId, secretAccessKey, sessionToken)
	return session.NewSession(aws.NewConfig().WithRegion(region).WithCredentials(creds))
}

// Create an AWS Credentials configuration with specific AWS credentials.
func CreateAwsCredentials(accessKeyId string, secretAccessKey string) *credentials.Credentials {
	creds := credentials.Value{AccessKeyID: accessKeyId, SecretAccessKey: secretAccessKey}
	return credentials.NewStaticCredentialsFromCreds(creds)
}

// Create an AWS Credentials configuration with temporary AWS credentials by including a session token (used for
// authenticating with MFA)
func CreateAwsCredentialsWithSessionToken(accessKeyId, secretAccessKey, sessionToken string) *credentials.Credentials {
	creds := credentials.Value{
		AccessKeyID:     accessKeyId,
		SecretAccessKey: secretAccessKey,
		SessionToken:    sessionToken,
	}
	return credentials.NewStaticCredentialsFromCreds(creds)
}

// Get a One-Time Password from the given mfaDevice. Per the RFC 6238 standard, this value will be different every 30 seconds.
func GetTimeBasedOneTimePassword(mfaDevice *iam.VirtualMFADevice) (string, error) {
	base32StringSeed := string(mfaDevice.Base32StringSeed)

	otp, err := totp.GenerateCode(base32StringSeed, time.Now())
	if err != nil {
		return "", err
	}

	return otp, nil
}

func ReadPasswordPolicyMinPasswordLength(iamClient *iam.IAM) (int, error) {
	output, err := iamClient.GetAccountPasswordPolicy(&iam.GetAccountPasswordPolicyInput{})
	if err != nil {
		return -1, err
	}

	return int(*output.PasswordPolicy.MinimumPasswordLength), nil
}

type CredentialsError struct {
	UnderlyingErr error
}

func (err CredentialsError) Error() string {
	return fmt.Sprintf("Error finding AWS credentials. Did you set the AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables or configure an AWS profile? Underlying error: %v", err.UnderlyingErr)
}
