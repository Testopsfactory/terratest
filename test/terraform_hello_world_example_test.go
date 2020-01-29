package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// website::tag::2::An example of how to test the simple Terraform module in examples/terraform-hello-world-example using Terratest.
func TestTerraformHelloWorldExample(t *testing.T) {
	t.Parallel()

	terraformOptions := &terraform.Options{
		// website::tag::2::Set the path to the Terraform code that will be tested.
		TerraformDir: "../examples/terraform-hello-world-example",
	}

	// website::tag::6::Clean up resources with "terraform destroy". Using "defer" runs the command at the end of the test, whether the test succeeds or fails.
	defer terraform.Destroy(t, terraformOptions)

	// website::tag::3::Run "terraform init" and "terraform apply" and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// website::tag::4::Run `terraform output` to get the values of output variables
	output := terraform.Output(t, terraformOptions, "hello_world")

	// website::tag::5::Check the output against expected values.
	assert.Equal(t, "Hello, World!", output)
}
