// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"strconv"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureAvailabilitySetExample(t *testing.T) {
	t.Parallel()

	// subscriptionID is overridden by the environment variable "ARM_SUBSCRIPTION_ID"
	subscriptionID := ""
	prefix := "terratest-avs"

	// Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// Relative path to the Terraform dir
		TerraformDir: "../../examples/azure/terraform-azure-availabilityset-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"prefix": prefix,
		},
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the values of output variables
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	availabilitySetName := terraform.Output(t, terraformOptions, "availability_set_name")
	fdc, _ := strconv.ParseInt(terraform.Output(t, terraformOptions, "availability_set_fdc"), 10, 32)
	avsSetFDC := int32(fdc)
	vmName := terraform.Output(t, terraformOptions, "vm_name")

	// Check the Availability Set Exists
	actualAvsExists := azure.AvailabilitySetExists(t, availabilitySetName, resourceGroupName, subscriptionID)
	assert.True(t, actualAvsExists)

	// Check the Availability set Fault Domain Count
	actualAvsFaultDomainCount := azure.GetAvailabilitySetFaultDomainCount(t, availabilitySetName, resourceGroupName, subscriptionID)
	assert.Equal(t, avsSetFDC, actualAvsFaultDomainCount)

	// Check the Availability Set for a VM
	assert.True(t, azure.CheckAvailabilitySetContainsVM(t, vmName, availabilitySetName, resourceGroupName, subscriptionID))
}
