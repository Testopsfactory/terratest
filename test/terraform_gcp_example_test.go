package test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

var projectId = os.Getenv("GOOGLE_CLOUD_PROJECT_ID")

func TestTerraformGcpExample(t *testing.T) {
	t.Parallel()

	// Give this bucket a unique name so we can distinguish it from any other bucket in your GCP account
	expectedName := fmt.Sprintf("terratest-gcp-example-%s", strings.ToLower(random.UniqueId()))

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/terraform-gcp-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"bucket_name": expectedName,
			"project":     projectId,
		},
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of an output variable
	bucketURL := terraform.Output(t, terraformOptions, "url")

	expectedURL := fmt.Sprintf("gs://%s", expectedName)

	// Verify that our expected name tag is one of the tags
	assert.Equal(t, expectedURL, bucketURL)

	// Verify that the Storage Bucket exists
	gcp.AssertStorageBucketExists(t, expectedName)
}
