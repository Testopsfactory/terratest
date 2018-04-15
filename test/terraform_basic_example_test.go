package test

import (
	"testing"
	"github.com/gruntwork-io/terratest/terraform"
	"github.com/stretchr/testify/assert"
)

// An example of how to test the simple Terraform module in examples/terraform-basic-example using Terratest.
func TerraformBasicExampleTest(t *testing.T) {
	t.Parallel()

	expectedText := "foo"

	terraformOptions := terraform.Options {
		TestName: t.Name(),
		TerraformDir: "../examples/terraform-basic-example",
		Vars: map[string]string {
			"example": expectedText,
		},
	}

	terraform.Apply(t, terraformOptions)
	defer terraform.Destroy(t, terraformOptions)

	actualText := terraform.Output(t, terraformOptions, "example")
	assert.Equal(t, expectedText, actualText)
}
