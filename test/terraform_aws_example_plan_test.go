package test

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// An example of how to test the Terraform module in examples/terraform-aws-example using Terratest.
func TestTerraformAwsExamplePlan(t *testing.T) {
	t.Parallel()

	// Make a copy of the terraform module to a temporary directory. This allows running multiple tests in parallel
	// against the same terraform module.
	exampleFolder := test_structure.CopyTerraformFolderToTemp(t, "../", "examples/terraform-aws-example")

	// Give this EC2 Instance a unique ID for a name tag so we can distinguish it from any other EC2 Instance running
	// in your AWS account
	expectedName := fmt.Sprintf("terratest-aws-example-%s", random.UniqueId())

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := aws.GetRandomStableRegion(t, nil, nil)

	// website::tag::1::Configure Terraform setting path to Terraform code, EC2 instance name, and AWS Region. We also
	// configure the options with default retryable errors to handle the most common retryable errors encountered in
	// terraform testing.
	planFilePath := filepath.Join(exampleFolder, "plan.out")
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/terraform-aws-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"instance_name": expectedName,
		},

		// Environment variables to set when running Terraform
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": awsRegion,
		},

		// Configure a plan file path so we can introspect the plan and make assertions about it.
		PlanFilePath: planFilePath,
	})

	// website::tag::2::Run `terraform init`, `terraform plan`, and `terraform show` and fail the test if there are any errors
	jsonOut := terraform.InitAndPlanAndShow(t, terraformOptions)

	// website::tag::3::Parse out the plan json into a generic map structure so that we can introspect it. You can
	// alternatively use https://github.com/hashicorp/terraform-json to get a concrete struct with all the types
	// resolved.
	var plan map[string]interface{}
	require.NoError(
		t,
		json.Unmarshal([]byte(jsonOut), &plan),
	)

	// Assert that the instance that is planned to be created has the expected tag set.
	plannedValues, hasType := plan["planned_values"].(map[string]interface{})
	require.True(t, hasType, "planned_values key in plan object is not a map")
	rootModule, hasType := plannedValues["root_module"].(map[string]interface{})
	require.True(t, hasType, "root_module key in planned_values in plan object is not a map")
	rootModuleResources, hasType := rootModule["resources"].([]interface{})
	require.True(t, hasType, "resources key in root_module in planned_values in plan object is not a list")
	require.Equal(t, 1, len(rootModuleResources))
	ec2InstanceResourcePlan, hasType := rootModuleResources[0].(map[string]interface{})
	require.True(t, hasType, "EC2 instance resource in plan object is not a map")

	resourceAddress, hasType := ec2InstanceResourcePlan["address"].(string)
	require.True(t, hasType, "Address is not a string")
	assert.Equal(t, "aws_instance.example", resourceAddress)

	resourceValues, hasType := ec2InstanceResourcePlan["values"].(map[string]interface{})
	require.True(t, hasType, "Values is not a map")
	tags, hasType := resourceValues["tags"].(map[string]interface{})
	require.True(t, hasType, "Tags in resource values is not a map")
	nameTagValue, hasType := tags["Name"].(string)
	require.True(t, hasType, "Name tag in tags is not a map")
	assert.Equal(t, expectedName, nameTagValue)
}
