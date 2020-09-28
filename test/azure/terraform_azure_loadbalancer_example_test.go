// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureLoadBalancerExample(t *testing.T) {
	t.Parallel()

	// subscriptionID is overridden by the environment variable "ARM_SUBSCRIPTION_ID"
	subscriptionID := ""
	rgName := fmt.Sprintf("terratest-loadbalancer-rg-%s", random.UniqueId())
	vnetName := fmt.Sprintf("vnet-%s", random.UniqueId())
	subnetName := fmt.Sprintf("subnet-%s", random.UniqueId())
	loadBalancer01Name := fmt.Sprintf("lb-public-%s", random.UniqueId())
	loadBalancer02Name := fmt.Sprintf("lb-private-%s", random.UniqueId())
	frontendForLB01 := fmt.Sprintf("cfg-%s", random.UniqueId())
	frontendForLB02 := fmt.Sprintf("cfg-%s", random.UniqueId())
	publicIPAddressForLB01 := fmt.Sprintf("pip-%s", random.UniqueId())
	privateIPForLB02 := "10.200.2.10"

	// Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-loadbalancer-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"resource_group_name":  rgName,
			"vnet_name":            vnetName,
			"subnet_name":          subnetName,
			"loadbalancer01_name":  loadBalancer01Name,
			"loadbalancer02_name":  loadBalancer02Name,
			"config_name_for_lb01": frontendForLB01,
			"config_name_for_lb02": frontendForLB02,
			"pip_for_lb01":         publicIPAddressForLB01,
			"privateip_for_lb02":   privateIPForLB02,
		},
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the values of output variables
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	expectedLB01Name := terraform.Output(t, terraformOptions, "loadbalancer_01_name")
	expectedLB02Name := terraform.Output(t, terraformOptions, "loadbalancer_02_name")
	expectedLB01FeConfigName := terraform.Output(t, terraformOptions, "config_name_for_lb01")
	expectedLB02FeConfigName := terraform.Output(t, terraformOptions, "config_name_for_lb02")
	expectedLB02PrivateIP := terraform.Output(t, terraformOptions, "privateip_for_lb02")

	t.Run("Public Load Balancer 01", func(t *testing.T) {
		// Check Public Load Balancer 01 exists.
		actualLB01Exists := azure.LoadBalancerExists(t, expectedLB01Name, resourceGroupName, subscriptionID)
		assert.True(t, actualLB01Exists)

		// Check Frontend Configuration for Load Balancer.
		actualLB01FeConfigNames := azure.GetLoadBalancerConfigNames(t, expectedLB01Name, resourceGroupName, subscriptionID)
		assert.Contains(t, actualLB01FeConfigNames, expectedLB01FeConfigName)

		// Check Frontend Configuration Public Address and Public IP assignment
		actualLB01IPAddress, actualLB01IPType := azure.GetLoadBalancerFrontendConfig(t, expectedLB01FeConfigName, expectedLB01Name, resourceGroupName, subscriptionID)
		assert.NotEmpty(t, actualLB01IPAddress)
		assert.Equal(t, azure.PublicIP, actualLB01IPType)
	})

	t.Run("Private Load Balancer 02", func(t *testing.T) {
		// Check Private Load Balancer 02 exists.
		actualLB02Exists := azure.LoadBalancerExists(t, expectedLB02Name, resourceGroupName, subscriptionID)
		assert.True(t, actualLB02Exists)

		// Check Frontend Configuration for Load Balancer.
		actualLB02FeConfigNames := azure.GetLoadBalancerConfigNames(t, expectedLB02Name, resourceGroupName, subscriptionID)
		assert.Contains(t, actualLB02FeConfigNames, expectedLB02FeConfigName)

		// Check Frontend Configuration Private IP Type and Address.
		actualLB02IPAddress, actualLB02IPType := azure.GetLoadBalancerFrontendConfig(t, expectedLB02FeConfigName, expectedLB02Name, resourceGroupName, subscriptionID)
		assert.NotEmpty(t, actualLB02IPAddress)
		assert.Equal(t, expectedLB02PrivateIP, actualLB02IPAddress)
		assert.Equal(t, azure.PrivateIP, actualLB02IPType)
	})
}
