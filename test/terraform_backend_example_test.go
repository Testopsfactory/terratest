package test

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

// An example of how to test the Terraform module in examples/terraform-backend-example using Terratest.
func TestTerraformBackendExample(t *testing.T) {
	t.Parallel()

	awsRegion := aws.GetRandomRegion(t, nil, nil)
	uniqueId := random.UniqueId()

	// Create an S3 bucket where we can store state
	bucketName := fmt.Sprintf("test-terraform-backend-example-%s", uniqueId)
	defer aws.DeleteS3Bucket(t, awsRegion, bucketName)
	aws.CreateS3Bucket(t, awsRegion, bucketName)

	// Deploy the module, configuring it to use the S3 bucket as an S3 backend
	terraformOptions := &terraform.Options{
		BackendConfig: map[string]interface{}{
			"bucket": bucketName,
			"key":    "terraform.tfstate",
			"region": awsRegion,
		},
	}

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// The module doesn't really *do* anything, so we just check a dummy output here and move on
	foo := terraform.OutputRequired(t, terraformOptions, "foo")
	require.Equal(t, "bar", foo)
}
