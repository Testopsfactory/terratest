// +build azure azureslim,compute

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureVmExample(t *testing.T) {
	t.Parallel()

	// subscriptionID is overridden by the environment variable "ARM_SUBSCRIPTION_ID"
	subscriptionID := ""
	uniquePostfix := random.UniqueId() // "resource" //
	expectedVmAdminUser := "testadmin"
	expectedVMSize := compute.VirtualMachineSizeTypes("Standard_B1s")
	expectedImageSKU := "2019-Datacenter-Core-smalldisk"
	expectedImageVersion := "latest"
	expectedDiskType := "Standard_LRS"
	expectedSubnetAddressRange := "10.0.17.0/24"
	expectedPrivateIPAddress := "10.0.17.4"
	expectedManagedDiskCount := 1
	expectedNicCount := 1

	// Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located.
		TerraformDir: "../../examples/azure/terraform-azure-vm-example",

		// Variables to pass to our Terraform code using -var options.
		Vars: map[string]interface{}{
			"postfix":          uniquePostfix,
			"user_name":        expectedVmAdminUser,
			"vm_size":          string(expectedVMSize),
			"vm_image_sku":     expectedImageSKU,
			"vm_image_version": expectedImageVersion,
			"disk_type":        expectedDiskType,
			"private_ip":       expectedPrivateIPAddress,
			"subnet_prefix":    expectedSubnetAddressRange,
		},
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the values of output variables.
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	virtualMachineName := terraform.Output(t, terraformOptions, "vm_name")
	expectedVNetName := terraform.Output(t, terraformOptions, "virtual_network_name")
	expectedSubnetName := terraform.Output(t, terraformOptions, "subnet_name")
	expectedPublicAddressName := terraform.Output(t, terraformOptions, "public_ip_name")
	expectedNicName := terraform.Output(t, terraformOptions, "network_interface_name")
	expectedAvsName := terraform.Output(t, terraformOptions, "availability_set_name")
	expectedOSDiskName := terraform.Output(t, terraformOptions, "os_disk_name")
	expectedDiskName := terraform.Output(t, terraformOptions, "managed_disk_name")

	// Comment for PR, to be removed:
	// Please let me know if there are too many tests or alternate examples, happy to reduce the
	// complexity and amount of code to be maintained. I tried to illustrate different approaches
	// we have used in various scenarios to illustrate the flexability of Terratest.

	t.Run("Strategies", func(t *testing.T) {
		// Check the VM Size directly.
		actualVMSize := azure.GetVirtualMachineSize(t, virtualMachineName, resourceGroupName, subscriptionID)
		assert.Equal(t, expectedVMSize, actualVMSize)

		// Check the VM Size by reference alternate example.
		// This strategy is beneficial when checking multiple properties by using one VM reference, avoiding
		// multiple SDK calls.
		vmRef := azure.GetVirtualMachine(t, virtualMachineName, resourceGroupName, subscriptionID)
		actualVMSize = vmRef.HardwareProfile.VMSize
		assert.Equal(t, expectedVMSize, actualVMSize)

		// Check the VM Size by instance alternate example.
		// This strategy is beneficial when checking multiple properties by using one VM instance and making
		// calls against it with the added benefit of property checks and abstraction.
		vmInstance := azure.GetVirtualMachineInstance(t, virtualMachineName, resourceGroupName, subscriptionID)
		actualVMSize = vmInstance.GetVirtualMachineInstanceSize()
		assert.Equal(t, expectedVMSize, actualVMSize)
	})

	t.Run("MultipleVMs", func(t *testing.T) {
		// This is beneficial when multiple VMs in a Resource Group need to be tested at once.

		// Check against all VM names in a Resource Group.
		vmList := azure.ListVirtualMachinesForResourceGroup(t, resourceGroupName, subscriptionID)
		assert.True(t, len(vmList) > 0)
		assert.Contains(t, vmList, virtualMachineName)

		// Get all VMs in a Resource Group by reference alternate example.
		// This strategy is beneficial when checking multiple VMs & their properties by avoiding
		// multiple SDK calls. The penalty for this approach is introducing direct references
		// which need to be checked for nil for optional configurations.
		vmsByRef := azure.GetVirtualMachinesForResourceGroup(t, resourceGroupName, subscriptionID)
		assert.True(t, len(*vmsByRef) > 0)

		// Check for the VM.
		thisVM := (*vmsByRef)[virtualMachineName]
		assert.Equal(t, expectedVMSize, thisVM.HardwareProfile.VMSize)

		// Check for the VM negative test.
		fakeVM := fmt.Sprintf("vm-%s", random.UniqueId())
		assert.Nil(t, (*vmsByRef)[fakeVM].VMID)
	})

	t.Run("Information", func(t *testing.T) {
		// Check if the Virtual Machine exists.
		assert.True(t, azure.VirtualMachineExists(t, virtualMachineName, resourceGroupName, subscriptionID))

		// Check the Admin User of the VM.
		actualVmAdminUser := azure.GetVirtualMachineAdminUser(t, virtualMachineName, resourceGroupName, subscriptionID)
		assert.Equal(t, expectedVmAdminUser, actualVmAdminUser)

		// Check the Storage Image properties of the VM.
		actualImage := azure.GetVirtualMachineImage(t, virtualMachineName, resourceGroupName, subscriptionID)
		assert.Equal(t, expectedImageSKU, actualImage.SKU)
		assert.Equal(t, expectedImageVersion, actualImage.Version)
	})

	t.Run("AvailabilitySet", func(t *testing.T) {
		// Check the Availability Set of the VM.
		// The AVS ID returned from the VM is always CAPS so ignoring case in the assertion.
		actualexpectedAvsName := azure.GetVirtualMachineAvailabilitySetID(t, virtualMachineName, resourceGroupName, subscriptionID)
		assert.True(t, strings.EqualFold(expectedAvsName, actualexpectedAvsName))

		// Check AVS for multiple VMs at a time alternate example.
		actualVMsInAvs := azure.GetAvailabilitySetVMNamesInCaps(t, expectedAvsName, resourceGroupName, subscriptionID)
		assert.Contains(t, actualVMsInAvs, strings.ToUpper(virtualMachineName))
	})

	t.Run("Disk", func(t *testing.T) {
		// Check the OS Disk name of the VM.
		actualOSDiskName := azure.GetVirtualMachineOsDiskName(t, virtualMachineName, resourceGroupName, subscriptionID)
		assert.Equal(t, expectedOSDiskName, actualOSDiskName)

		// Check the Managed Disk count of the VM.
		actualManagedDiskCount := azure.GetVirtualMachineManagedDiskCount(t, virtualMachineName, resourceGroupName, subscriptionID)
		assert.Equal(t, expectedManagedDiskCount, actualManagedDiskCount)

		// Check the VM Managed Disk exists in the list of all VM Managed Disks.
		actualManagedDiskNames := azure.GetVirtualMachineManagedDisks(t, virtualMachineName, resourceGroupName, subscriptionID)
		assert.Contains(t, actualManagedDiskNames, expectedDiskName)

		// Check the Disk Type of the Managed Disk of the VM.
		actualDiskType := azure.GetDiskType(t, expectedDiskName, resourceGroupName, subscriptionID)
		assert.Equal(t, compute.DiskStorageAccountTypes(expectedDiskType), actualDiskType)
	})

	// See the Terratest Azure Network Example for other related tests.
	t.Run("NetworkInterface", func(t *testing.T) {
		// Check the Network Interface count of the VM.
		actualNicCount := azure.GetVirtualMachineNicCount(t, virtualMachineName, resourceGroupName, subscriptionID)
		assert.Equal(t, expectedNicCount, actualNicCount)

		// Check the VM Network Interface exists in the list of all VM Network Interfaces.
		actualNics := azure.GetVirtualMachineNics(t, virtualMachineName, resourceGroupName, subscriptionID)
		assert.Contains(t, actualNics, expectedNicName)

		// Check for the Private IP in the NICs IP list.
		actualPrivateIPAddress := azure.GetNetworkInterfacePrivateIPs(t, expectedNicName, resourceGroupName, subscriptionID)
		assert.Contains(t, actualPrivateIPAddress, expectedPrivateIPAddress)

		// Check for the Public IP for the NIC. No expected value since it is assigned runtime.
		actualPublicIP := azure.GetIPOfPublicIPAddressByName(t, expectedPublicAddressName, resourceGroupName, subscriptionID)
		assert.NotNil(t, actualPublicIP)
	})

	t.Run("Vnet&Subnet", func(t *testing.T) {
		// Check the Subnet exists in the Virtual Network.
		actualVnetSubnets := azure.GetVirtualNetworkSubnets(t, expectedVNetName, resourceGroupName, subscriptionID)
		assert.NotNil(t, actualVnetSubnets[expectedVNetName])

		// Check the Private IP is in the Subnet Range.
		actualVMNicIPInSubnet := azure.CheckSubnetContainsIP(t, expectedPrivateIPAddress, expectedSubnetName, expectedVNetName, resourceGroupName, subscriptionID)
		assert.True(t, actualVMNicIPInSubnet)
	})
}
