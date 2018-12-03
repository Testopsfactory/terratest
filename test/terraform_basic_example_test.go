package test

import (
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// An example of how to test the simple Terraform module in examples/terraform-basic-example using Terratest.
func TestTerraformBasicExample(t *testing.T) {
	t.Parallel()
	terraformDir := "../examples/terraform-basic-example"

	expectedText := "foo"

	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: terraformDir,

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"example": expectedText,
		},
		// Only create the example1 file
		Targets: []string{"local_file.example"},

		// Use the var files
		VarFiles: []string{"varfile.tfvars"},

		NoColor: true,
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of an output variable
	actualText := terraform.Output(t, terraformOptions, "example")

	// Verify we're getting back the variable we expect
	assert.Equal(t, expectedText, actualText)

	// Test for the second variable which comes from the var file
	actualText = terraform.Output(t, terraformOptions, "example2")
	assert.Equal(t, "test", actualText)

	// Test resources: the target specifies only the local_file.example resource
	//only the file  example.txt should be created, and example2.txt should NOT
	assert.True(t, files.FileExists(filepath.Join(terraformDir, "example.txt")))
	assert.False(t, files.FileExists(filepath.Join(terraformDir, "example2.txt")))

}
