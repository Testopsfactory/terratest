package test

import (
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureDataFactoryExample(t *testing.T) {
	t.Parallel()

	uniquePostfix := strings.ToLower(random.UniqueId())
	expectedDataFactoryProvisioningState := "Succeeded"
	expectedLocation := "West US2"

	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-datafactory-example",
		Vars: map[string]interface{}{
			"postfix":  uniquePostfix,
			"location": expectedLocation,
		},
	}

	// website::tag::4:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// website::tag::2:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// website::tag::3:: Run `terraform output` to get the values of output variables
	expectedResourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	expectedDataFactoryName := terraform.Output(t, terraformOptions, "datafactory_name")

	// // website::tag::4:: Get synapse details and assert them against the terraform output
	actualDataFactory := azure.GetDataFactory(t, expectedResourceGroupName, expectedDataFactoryName, "")

	assert.Equal(t, expectedDataFactoryName, *actualDataFactory.Name)
	assert.Equal(t, expectedDataFactoryProvisioningState, *actualDataFactory.FactoryProperties.ProvisioningState)

	actualDataFactoryExits := azure.DataFactoryExists(t, expectedDataFactoryName, expectedResourceGroupName, "")
	assert.True(t, actualDataFactoryExits)
}
