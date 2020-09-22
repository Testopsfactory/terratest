// +build azure azureslim,network

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureNetworkExample(t *testing.T) {
	t.Parallel()

	// Create values for Terraform
	subscriptionID := "" // subscriptionID is overridden by the environment variable "ARM_SUBSCRIPTION_ID"
	prefix := fmt.Sprintf("terratest-net-%s", random.UniqueId())
	expectedLocation := "eastus2"
	expectedSubnetRange := "10.0.20.0/24"
	expectedPrivateIP := "10.0.20.5"
	expectedDnsIp01 := "10.0.0.5"
	expectedDnsIp02 := "10.0.0.6"
	exectedDNSLabel := strings.ToLower(prefix) // only lowercase, numeric and hyphens chars allowed for DNS

	// Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// Relative path to the Terraform dir
		TerraformDir: "../../examples/azure/terraform-azure-network-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"prefix":            prefix,
			"subnet_prefix":     expectedSubnetRange,
			"private_ip":        expectedPrivateIP,
			"dns_ip_01":         expectedDnsIp01,
			"dns_ip_02":         expectedDnsIp02,
			"location":          expectedLocation,
			"domain_name_label": exectedDNSLabel,
		},
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the values of output variables
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	virtualNetworkName := terraform.Output(t, terraformOptions, "virtual_network_name")
	subnetName := terraform.Output(t, terraformOptions, "subnet_name")
	publicAddressName := terraform.Output(t, terraformOptions, "public_address_name")
	nicInternalName := terraform.Output(t, terraformOptions, "network_interface_internal")
	nicExternalName := terraform.Output(t, terraformOptions, "network_interface_external")

	// Tests are separated into subtests to differentiate integrated tests and pure resource tests

	// Integrated network resource tests
	t.Run("VirtualNetwork_Subnet", func(t *testing.T) {
		// Check the Subnet exists in the Virtual Network Subnets with the expected Address Prefix
		actualVnetSubnets := azure.GetVirtualNetworkSubnets(t, virtualNetworkName, resourceGroupName, subscriptionID)
		assert.NotNil(t, actualVnetSubnets[subnetName])
		assert.Equal(t, expectedSubnetRange, actualVnetSubnets[subnetName])
	})

	t.Run("NIC_PublicAddress", func(t *testing.T) {
		// Check the internal network interface does NOT have a public IP
		actualPrivateIPOnly := azure.GetNetworkInterfacePublicIPs(t, nicInternalName, resourceGroupName, subscriptionID)
		assert.True(t, len(actualPrivateIPOnly) == 0)

		// Check the external network interface has a public IP
		actualPublicIPs := azure.GetNetworkInterfacePublicIPs(t, nicExternalName, resourceGroupName, subscriptionID)
		assert.True(t, len(actualPublicIPs) == 1)
	})

	t.Run("Subnet_NIC", func(t *testing.T) {
		// Check the private IP is in the subnet range
		checkPrivateIpInSubnet := azure.CheckSubnetContainsIP(t, expectedPrivateIP, subnetName, virtualNetworkName, resourceGroupName, subscriptionID)
		assert.True(t, checkPrivateIpInSubnet)
	})

	// Tests for individual network resources
	t.Run("Network", func(t *testing.T) {
		// Check the Virtual Network exists
		assert.True(t, azure.VirtualNetworkExists(t, virtualNetworkName, resourceGroupName, subscriptionID))

		// Check the Subnet exists
		assert.True(t, azure.SubnetExists(t, subnetName, virtualNetworkName, resourceGroupName, subscriptionID))

		// Check the Network Interfaces exist
		assert.True(t, azure.NetworkInterfaceExists(t, nicInternalName, resourceGroupName, subscriptionID))
		assert.True(t, azure.NetworkInterfaceExists(t, nicExternalName, resourceGroupName, subscriptionID))

		// Check Public Address exists
		assert.True(t, azure.PublicAddressExists(t, publicAddressName, resourceGroupName, subscriptionID))

		// Check the Virtual Network DNS server IPs
		actualDNSIPs := azure.GetVirtualNetworkDNSServerIPs(t, virtualNetworkName, resourceGroupName, subscriptionID)
		assert.Contains(t, actualDNSIPs, expectedDnsIp01)
		assert.Contains(t, actualDNSIPs, expectedDnsIp02)

		// Check the Network Interface private IP
		actualPrivateIPs := azure.GetNetworkInterfacePrivateIPs(t, nicInternalName, resourceGroupName, subscriptionID)
		assert.Contains(t, actualPrivateIPs, expectedPrivateIP)

		// Check the Public Address public IP
		actualPublicIP := azure.GetPublicAddressIP(t, publicAddressName, resourceGroupName, subscriptionID)
		assert.True(t, len(actualPublicIP) > 0)

		// Check DNS created for this example is reserved
		actualDnsNotAvailable := azure.CheckPublicDNSNameAvailability(t, expectedLocation, exectedDNSLabel, subscriptionID)
		assert.False(t, actualDnsNotAvailable)

		// Check new randomized DNS is available
		newDNSLabel := strings.ToLower(fmt.Sprintf("%s-%s-%s", prefix, random.UniqueId(), random.UniqueId()))
		actualDnsAvailable := azure.CheckPublicDNSNameAvailability(t, expectedLocation, newDNSLabel, subscriptionID)
		assert.True(t, actualDnsAvailable)
	})

}
