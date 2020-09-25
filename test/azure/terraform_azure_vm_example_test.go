// +build azure azureslim,compute

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
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

	subscriptionID := "" // Subscription ID, leave blank if available as an Environment Var
	uniquePostfix := random.UniqueId()
	expectedVmAdminUser := "testadmin"
	expectedVMSize := "Standard_DS1_v2"
	expectedImageSKU := "2019-Datacenter-Core-smalldisk"
	expectedImageVersion := "latest"
	expectedDiskType := "Standard_LRS"
	expectedSubnetAddressRange := "10.0.17.0/24"
	expectedPrivateIPAddress := "10.0.17.4"
	var expectedAvsFaultDomainCount int32 = 2
	expectedManagedDiskCount := 1
	expectedNicCount := 1

	// Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-vm-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"postfix":          uniquePostfix,
			"user_name":        expectedVmAdminUser,
			"vm_size":          expectedVMSize,
			"vm_image_sku":     expectedImageSKU,
			"vm_image_version": expectedImageVersion,
			"disk_type":        expectedDiskType,
			"private_ip":       expectedPrivateIPAddress,
			"subnet_prefix":    expectedSubnetAddressRange,
		},
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the values of output variables
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	expectedVMName := terraform.Output(t, terraformOptions, "vm_name")
	expectedVNetName := terraform.Output(t, terraformOptions, "virtual_network_name")
	expectedSubnetName := terraform.Output(t, terraformOptions, "subnet_name")
	expectedPublicAddressName := terraform.Output(t, terraformOptions, "public_ip_name")
	expectedNicName := terraform.Output(t, terraformOptions, "network_interface_name")
	expectedAvsName := terraform.Output(t, terraformOptions, "availability_set_name")
	expectedOSDiskName := terraform.Output(t, terraformOptions, "os_disk_name")
	expectedDiskName := terraform.Output(t, terraformOptions, "managed_disk_name")

	t.Run("Strategies", func(t *testing.T) {
		// Check the VM Size directly
		actualVMSize := azure.GetVirtualMachineSize(t, expectedVMName, resourceGroupName, subscriptionID)
		assert.Equal(t, expectedVMSize, string(actualVMSize))

		// Check the VM Size by object ref
		vmRef := azure.GetVirtualMachine(t, expectedVMName, resourceGroupName, subscriptionID)
		actualVMSize = vmRef.HardwareProfile.VMSize
		assert.Equal(t, expectedVMSize, string(actualVMSize))

		// Check the VM Size by instance getter
		vmInstance := azure.GetVirtualMachineInstance(t, expectedVMName, resourceGroupName, subscriptionID)
		actualVMSize = vmInstance.GetVirtualMachineInstanceSize()
		assert.Equal(t, expectedVMSize, string(actualVMSize))
	})

	t.Run("MultipleVMs", func(t *testing.T) {
		// Get a list of all VMs and confirm one (or more) VMs exist
		vmList := azure.GetResourceGroupVirtualMachines(t, resourceGroupName, subscriptionID)
		assert.True(t, len(vmList) > 0)
		assert.Contains(t, vmList, expectedVMName)

		// Get all VMs by ref (warning: pointer de-ref panic if vm is not in list!)
		vmsByRef := azure.GetResourceGroupVirtualMachinesObjects(t, resourceGroupName, subscriptionID)
		assert.True(t, len(*vmsByRef) > 0)
		thisVm := (*vmsByRef)[expectedVMName]
		assert.Equal(t, expectedVMSize, string(thisVm.HardwareProfile.VMSize))
	})

	t.Run("Information", func(t *testing.T) {
		// Check the Virtual Machine exists
		assert.True(t, azure.VirtualMachineExists(t, expectedVMName, resourceGroupName, subscriptionID))

		// Check the Admin User
		actualVmAdminUser := azure.GetVirtualMachineAdminUser(t, expectedVMName, resourceGroupName, subscriptionID)
		assert.Equal(t, expectedVmAdminUser, actualVmAdminUser)

		// Check the Storage Image reference
		actualImage := azure.GetVirtualMachineImage(t, expectedVMName, resourceGroupName, subscriptionID)
		assert.Equal(t, expectedImageSKU, actualImage.SKU)
		assert.Equal(t, expectedImageVersion, actualImage.Version)
	})

	t.Run("AvailabilitySet", func(t *testing.T) {
		// Check the Availability Set
		actualexpectedAvsName := azure.GetVirtualMachineAvailabilitySetID(t, expectedVMName, resourceGroupName, subscriptionID)
		assert.True(t, strings.EqualFold(expectedAvsName, actualexpectedAvsName))

		// Check the Availability set fault domain counts
		actualAvsFaultDomainCount := azure.GetAvailabilitySetFaultDomainCount(t, expectedAvsName, resourceGroupName, subscriptionID)
		assert.Equal(t, expectedAvsFaultDomainCount, actualAvsFaultDomainCount)

		actualVMsInAvs := azure.GetAvailabilitySetVMNamesInCaps(t, expectedAvsName, resourceGroupName, subscriptionID)
		assert.Contains(t, actualVMsInAvs, strings.ToUpper(expectedVMName))
	})

	t.Run("Disk", func(t *testing.T) {
		// Check the OS Disk name
		actualOSDiskName := azure.GetVirtualMachineOsDiskName(t, expectedVMName, resourceGroupName, subscriptionID)
		assert.Equal(t, expectedOSDiskName, actualOSDiskName)

		// Check the Managed Disk count
		actualManagedDiskCount := azure.GetVirtualMachineManagedDiskCount(t, expectedVMName, resourceGroupName, subscriptionID)
		assert.Equal(t, expectedManagedDiskCount, actualManagedDiskCount)

		// Check the VM Managed Disk exists in the list of all VM Managed Disks
		actualManagedDiskNames := azure.GetVirtualMachineManagedDisks(t, expectedVMName, resourceGroupName, subscriptionID)
		assert.Contains(t, actualManagedDiskNames, expectedDiskName)

		// Check the Disk Type
		actualDiskType := azure.GetDiskType(t, expectedDiskName, resourceGroupName, subscriptionID)
		assert.Equal(t, compute.DiskStorageAccountTypes(expectedDiskType), actualDiskType)
	})

	t.Run("NetworkInterface", func(t *testing.T) {
		// Check the Network Interface count
		actualNicCount := azure.GetVirtualMachineNicCount(t, expectedVMName, resourceGroupName, subscriptionID)
		assert.Equal(t, expectedNicCount, actualNicCount)

		// Check the VM Network Interface exists in the list of all VM Network Interfaces
		actualNics := azure.GetVirtualMachineNics(t, expectedVMName, resourceGroupName, subscriptionID)
		assert.Contains(t, actualNics, expectedNicName)

		// Check the Private IP
		actualNicIPs := azure.GetNetworkInterfacePrivateIPs(t, expectedNicName, resourceGroupName, subscriptionID)
		assert.Contains(t, actualNicIPs, expectedPrivateIPAddress)

		// Check the Public IP exists
		actualPublicIP := azure.GetPublicAddressIP(t, expectedPublicAddressName, resourceGroupName, subscriptionID)
		assert.NotNil(t, actualPublicIP)
	})

	t.Run("Vnet&Subnet", func(t *testing.T) {
		// Check the Subnet exists in the Virtual Network Subnets
		actualVnetSubnets := azure.GetVirtualNetworkSubnets(t, expectedVNetName, resourceGroupName, subscriptionID)
		assert.NotNil(t, actualVnetSubnets[expectedVNetName])

		// Check the Private IP is in the Subnet Range
		actualVMNicIPInSubnet := azure.CheckSubnetContainsIP(t, expectedPrivateIPAddress, expectedSubnetName, expectedVNetName, resourceGroupName, subscriptionID)
		assert.True(t, actualVMNicIPInSubnet)
	})
}
