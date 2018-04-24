package aws

import (
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/aws"
	"testing"
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
	kmsClient, err := CreateKmsClient(region)
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
func CreateKmsClient(region string) (*kms.KMS, error) {
	sess, err := GetAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return kms.New(sess), nil
}
