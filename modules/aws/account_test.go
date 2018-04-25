package aws

import "testing"

func TestExtractAccountIdFromValidArn(t *testing.T) {
	t.Parallel()

	expectedAccountId := "123456789012"
	arn := "arn:aws:iam::" + expectedAccountId + ":user/test"

	actualAccountId, err := extractAccountIdFromArn(arn)
	if err != nil {
		t.Fatalf("Unexpected error while extracting account id from arn %s: %s", arn, err)
	}

	if actualAccountId != expectedAccountId {
		t.Fatalf("Did not get expected account id. Expected: %s. Actual: %s.", expectedAccountId, actualAccountId)
	}
}

func TestExtractAccountIdFromInvalidArn(t *testing.T) {
	t.Parallel()

	_, err := extractAccountIdFromArn("invalidArn")
	if err == nil {
		t.Fatalf("Expected an error when extracting an account id from an invalid ARN, but got nil")
	}
}

