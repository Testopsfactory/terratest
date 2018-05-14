package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
)

// Get the ARN of a KMS Customer Master Key (CMK) in the given region with the given ID. The ID can be an alias, such
// as "alias/my-cmk".
func GetCmkArn(t *testing.T, region string, cmkId string) string {
	out, err := GetCmkArnE(t, region, cmkId)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// Get the ARN of a KMS Customer Master Key (CMK) in the given region with the given ID. The ID can be an alias, such
// as "alias/my-cmk".
func GetCmkArnE(t *testing.T, region string, cmkId string) (string, error) {
	kmsClient, err := NewKmsClientE(t, region)
	if err != nil {
		return "", err
	}

	result, err := kmsClient.DescribeKey(&kms.DescribeKeyInput{
		KeyId: aws.String(cmkId),
	})

	if err != nil {
		return "", err
	}

	return *result.KeyMetadata.Arn, nil
}

// Create a KMS client
func NewKmsClient(t *testing.T, region string) *kms.KMS {
	client, err := NewKmsClientE(t, region)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// Create a KMS client
func NewKmsClientE(t *testing.T, region string) (*kms.KMS, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return kms.New(sess), nil
}
