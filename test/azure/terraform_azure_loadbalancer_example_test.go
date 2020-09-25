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

	// Initialize resource names with random unique suffixes.
	resourceGroupName := fmt.Sprintf("terratest-loadbalancer-rg-%s", random.UniqueId())
	loadBalancer01Name := fmt.Sprintf("lb-public-%s", random.UniqueId())
	loadBalancer02Name := fmt.Sprintf("lb-private-%s", random.UniqueId())

	frontendIPConfigForLB01 := fmt.Sprintf("cfg-%s", random.UniqueId())
	publicIPAddressForLB01 := fmt.Sprintf("pip-%s", random.UniqueId())

	vnetForLB02 := fmt.Sprintf("vnet-%s", random.UniqueId())
	frontendSubnetID := fmt.Sprintf("snt-%s", random.UniqueId())

	// Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-loadbalancer-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"resource_group_name": resourceGroupName,
			"loadbalancer01_name": loadBalancer01Name,
			"loadbalancer02_name": loadBalancer02Name,
			"vnet_name":           vnetForLB02,
			"lb01_feconfig":       frontendIPConfigForLB01,
			"pip_forlb01":         publicIPAddressForLB01,
			"feSubnet_forlb02":    frontendSubnetID,
		},
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	t.Run("Load Balancer 01", func(t *testing.T) {
		// load balancer 01 (with Public IP) exists
		lb01Exists := azure.LoadBalancerExists(t, loadBalancer01Name, resourceGroupName, "")
		assert.True(t, lb01Exists)
	})

	t.Run("Frontend Config for LB01", func(t *testing.T) {
		// Read the LB information
		lb01 := azure.GetLoadBalancer(t, loadBalancer01Name, resourceGroupName, "")
		lb01Props := lb01.LoadBalancerPropertiesFormat
		fe01Config := (*lb01Props.FrontendIPConfigurations)[0]

		// Verify settings
		assert.Equal(t, frontendIPConfigForLB01, *fe01Config.Name, "LB01 Frontend IP config name")
	})

	t.Run("IP Checks for LB01", func(t *testing.T) {
		// Get config from LB01, including IP Address and verify Public IP
		ipAddress, ipType := azure.GetLoadBalancerFrontendConfig(t, loadBalancer01Name, resourceGroupName, "")
		assert.NotEmpty(t, ipAddress)
		assert.Equal(t, string(azure.PublicIP), ipType)
	})

	t.Run("Load Balancer 02", func(t *testing.T) {
		// load balancer 02 (with Private IP on vnet/subnet) exists
		lb02Exists := azure.LoadBalancerExists(t, loadBalancer02Name, resourceGroupName, "")
		assert.True(t, lb02Exists)
	})

	t.Run("IP Check for Load Balancer 02", func(t *testing.T) {
		// Get config from LB02, including IP Address and verify Private IP

		ipAddress, ipType := azure.GetLoadBalancerFrontendConfig(t, loadBalancer02Name, resourceGroupName, "")
		assert.NotEmpty(t, ipAddress)
		assert.Equal(t, string(azure.PrivateIP), ipType)
	})
}
