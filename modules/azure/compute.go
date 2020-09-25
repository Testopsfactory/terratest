package azure

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// VirtualMachineExists indicates whether the specifcied Azure Virtual Machine exists.
// This function would fail the test if there is an error.
func VirtualMachineExists(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) bool {
	exists, err := VirtualMachineExistsE(t, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return exists
}

// VirtualMachineExistsE indicates whether the specifcied Azure Virtual Machine exists.
func VirtualMachineExistsE(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) (bool, error) {
	// Get VM Object
	_, err := GetVirtualMachineE(t, vmName, resGroupName, subscriptionID)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetVirtualMachineNics gets a list of Network Interface names for a specifcied Azure Virtual Machine.
// This function would fail the test if there is an error.
func GetVirtualMachineNics(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) []string {
	nicList, err := GetVirtualMachineNicsE(t, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return nicList
}

// GetVirtualMachineNicsE gets a list of Network Interface names for a specified Azure Virtual Machine.
func GetVirtualMachineNicsE(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) ([]string, error) {
	nics := []string{}

	// Get VM Object
	vm, err := GetVirtualMachineE(t, vmName, resGroupName, subscriptionID)
	if err != nil {
		return nics, err
	}

	vmNICs := *vm.NetworkProfile.NetworkInterfaces
	if len(vmNICs) == 0 {
		// No VM NICs attached is still valid but returning a meaningful error
		return nics, errors.New("No network interface attached to this Virtual Machine")
	}

	// Get the attached NIC names
	for _, nic := range vmNICs {
		nics = append(nics, GetNameFromResourceID(*nic.ID))
	}
	return nics, nil
}

// GetVirtualMachineNicCount gets the Network Interface count of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
func GetVirtualMachineNicCount(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) int {
	nicCount, err := GetVirtualMachineNicCountE(t, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return nicCount
}

// GetVirtualMachineNicCountE gets the Network Interface count of the specified Azure Virtual Machine.
func GetVirtualMachineNicCountE(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) (int, error) {
	nicCount := 0

	// Get VM Object
	vm, err := GetVirtualMachineE(t, vmName, resGroupName, subscriptionID)
	if err != nil {
		return nicCount, err
	}

	return len(*vm.NetworkProfile.NetworkInterfaces), nil
}

// GetVirtualMachineManagedDisks gets the list of Managed Disk names of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
func GetVirtualMachineManagedDisks(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) []string {
	diskNames, err := GetVirtualMachineManagedDisksE(t, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return diskNames
}

// GetVirtualMachineManagedDisksE gets the list of Managed Disk names of the specified Azure Virtual Machine.
func GetVirtualMachineManagedDisksE(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) ([]string, error) {
	diskNames := []string{}

	// Get VM Object
	vm, err := GetVirtualMachineE(t, vmName, resGroupName, subscriptionID)
	if err != nil {
		return diskNames, err
	}

	// Get VM Attached Disks
	vmDisks := *vm.StorageProfile.DataDisks
	for _, v := range vmDisks {
		diskNames = append(diskNames, *v.Name)
	}

	return diskNames, nil
}

// GetVirtualMachineManagedDiskCount gets the Managed Disk count of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
func GetVirtualMachineManagedDiskCount(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) int {
	mngDiskCount, err := GetVirtualMachineManagedDiskCountE(t, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return mngDiskCount
}

// GetVirtualMachineManagedDiskCountE gets the Managed Disk count of the specified Azure Virtual Machine.
func GetVirtualMachineManagedDiskCountE(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) (int, error) {
	mngDiskCount := -1

	// Get VM Object
	vm, err := GetVirtualMachineE(t, vmName, resGroupName, subscriptionID)
	if err != nil {
		return mngDiskCount, err
	}

	return len(*vm.StorageProfile.DataDisks), nil
}

// GetVirtualMachineOsDiskName gets the OS Disk name of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
func GetVirtualMachineOsDiskName(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) string {
	osDiskName, err := GetVirtualMachineOsDiskNameE(t, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return osDiskName
}

// GetVirtualMachineOsDiskNameE gets the OS Disk name of the specified Azure Virtual Machine.
func GetVirtualMachineOsDiskNameE(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) (string, error) {
	// Get VM Object
	vm, err := GetVirtualMachineE(t, vmName, resGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return *vm.StorageProfile.OsDisk.Name, nil
}

// GetVirtualMachineState gets the State of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
func GetVirtualMachineState(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) string {
	vmState, err := GetVirtualMachineStateE(t, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return vmState
}

// GetVirtualMachineStateE gets the State of the specified Azure Virtual Machine.
func GetVirtualMachineStateE(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) (string, error) {
	// Get VM Object
	vm, err := GetVirtualMachineE(t, vmName, resGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return vm.Status, nil
}

// GetVirtualMachineAvailabilitySetID gets the Availability Set ID of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
func GetVirtualMachineAvailabilitySetID(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) string {
	adminUser, err := GetVirtualMachineAvailabilitySetIDE(t, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return adminUser
}

// GetVirtualMachineAvailabilitySetIDE gets the Availability Set ID of the specified Azure Virtual Machine.
func GetVirtualMachineAvailabilitySetIDE(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) (string, error) {
	// Get VM Object
	vm, err := GetVirtualMachineE(t, vmName, resGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return GetNameFromResourceID(*vm.AvailabilitySet.ID), nil
}

// VMImage represents the storage image for the specified Azure Virtual Machine.
type VMImage struct {
	Publisher string
	Offer     string
	SKU       string
	Version   string
}

// GetVirtualMachineImage gets the Image of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
func GetVirtualMachineImage(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) VMImage {
	adminUser, err := GetVirtualMachineImageE(t, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return adminUser
}

// GetVirtualMachineImageE gets the Image  of the specified Azure Virtual Machine.
func GetVirtualMachineImageE(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) (VMImage, error) {
	vmImage := VMImage{}

	// Get VM Object
	vm, err := GetVirtualMachineE(t, vmName, resGroupName, subscriptionID)
	if err != nil {
		return vmImage, err
	}

	// Populate VM Image
	vmImage.Publisher = *vm.StorageProfile.ImageReference.Publisher
	vmImage.Offer = *vm.StorageProfile.ImageReference.Offer
	vmImage.SKU = *vm.StorageProfile.ImageReference.Sku
	vmImage.Version = *vm.StorageProfile.ImageReference.Version

	return vmImage, nil
}

// GetVirtualMachineAdminUser gets the Admin Username of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
func GetVirtualMachineAdminUser(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) string {
	adminUser, err := GetVirtualMachineAdminUserE(t, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return adminUser
}

// GetVirtualMachineAdminUserE gets the Admin Username of the specified Azure Virtual Machine.
func GetVirtualMachineAdminUserE(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) (string, error) {
	// Get VM Object
	vm, err := GetVirtualMachineE(t, vmName, resGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return string(*vm.OsProfile.AdminUsername), nil
}

// GetVirtualMachineSize gets the Size Type of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
func GetVirtualMachineSize(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) compute.VirtualMachineSizeTypes {
	size, err := GetVirtualMachineSizeE(t, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return size
}

// GetVirtualMachineSizeE gets the Size Type of the specified Azure Virtual Machine.
func GetVirtualMachineSizeE(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) (compute.VirtualMachineSizeTypes, error) {
	// Get VM Object
	vm, err := GetVirtualMachineE(t, vmName, resGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return vm.VirtualMachineProperties.HardwareProfile.VMSize, nil
}

// GetVirtualMachineTags gets the Tags of the specified Virtual Machine as a map.
// This function would fail the test if there is an error.
func GetVirtualMachineTags(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) map[string]string {
	tags, err := GetVirtualMachineTagsE(t, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return tags
}

// GetVirtualMachineTagsE gets the Tags of the specified Virtual Machine as a map.
func GetVirtualMachineTagsE(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) (map[string]string, error) {
	// Setup a blank map to populate and return
	tags := make(map[string]string)

	// Get VM Object
	vm, err := GetVirtualMachineE(t, vmName, resGroupName, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Range through existing tags and populate above map accordingly
	for k, v := range vm.Tags {
		tags[k] = *v
	}

	return tags, nil
}

// ***************************************************** //
// Get multiple Virtual Machines from a Resource Group
// ***************************************************** //

// GetResourceGroupVirtualMachines gets a list of all Virtual Machine names in the specified Resource Group.
// This function would fail the test if there is an error.
func GetResourceGroupVirtualMachines(t testing.TestingT, resGroupName string, subscriptionID string) []string {
	vms, err := GetResourceGroupVirtualMachinesE(t, resGroupName, subscriptionID)
	require.NoError(t, err)
	return vms
}

// GetResourceGroupVirtualMachinesE gets a list of all Virtual Machine names in the specified Resource Group.
func GetResourceGroupVirtualMachinesE(t testing.TestingT, resourceGroupName string, subscriptionID string) ([]string, error) {
	vmDetails := []string{}

	vmClient, err := GetVirtualMachineClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	vms, err := vmClient.List(context.Background(), resourceGroupName)
	if err != nil {
		return nil, err
	}

	for _, v := range vms.Values() {
		vmDetails = append(vmDetails, *v.Name)
	}
	return vmDetails, nil
}

// GetResourceGroupVirtualMachinesObjects gets all Virtual Machine objects in the specified Resource Group.
// This function would fail the test if there is an error.
func GetResourceGroupVirtualMachinesObjects(t testing.TestingT, resGroupName string, subscriptionID string) *map[string]compute.VirtualMachineProperties {
	vms, err := GetResourceGroupVirtualMachinesObjectsE(t, resGroupName, subscriptionID)
	require.NoError(t, err)
	return vms
}

// GetResourceGroupVirtualMachinesObjectsE gets all Virtual Machine objects in the specified Resource Group.
func GetResourceGroupVirtualMachinesObjectsE(t testing.TestingT, resourceGroupName string, subscriptionID string) (*map[string]compute.VirtualMachineProperties, error) {
	vmClient, err := GetVirtualMachineClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	vms, err := vmClient.List(context.Background(), resourceGroupName)
	if err != nil {
		return nil, err
	}

	vmDetails := make(map[string]compute.VirtualMachineProperties)
	for _, v := range vms.Values() {
		machineName := v.Name
		vmProperties := v.VirtualMachineProperties
		vmDetails[string(*machineName)] = *vmProperties
	}
	return &vmDetails, nil
}

// ******************************************************************** //
// Get VM using Instance and Instance property get, reducing SKD calls
// ******************************************************************** //

// Instance of the VM
type Instance struct {
	*compute.VirtualMachine
}

// GetVirtualMachineInstance gets a local Virtual Machine instance in the specified Resource Group.
// This function would fail the test if there is an error.
func GetVirtualMachineInstance(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) *Instance {
	vm, err := GetVirtualMachineInstanceE(t, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return vm
}

// GetVirtualMachineInstanceE gets a local Virtual Machine instance in the specified Resource Group.
func GetVirtualMachineInstanceE(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) (*Instance, error) {
	// Get VM Object
	vm, err := GetVirtualMachineE(t, vmName, resGroupName, subscriptionID)
	if err != nil {
		return nil, err
	}

	return &Instance{vm}, nil
}

// GetVirtualMachineInstanceSize gets the size of the Virtual Machine.
func (vm *Instance) GetVirtualMachineInstanceSize() compute.VirtualMachineSizeTypes {
	return vm.VirtualMachineProperties.HardwareProfile.VMSize
}

// *********************** //
// Get the base VM Object
// *********************** //

// GetVirtualMachine gets a Virtual Machine in the specified Azure Resource Group.
// This function would fail the test if there is an error.
func GetVirtualMachine(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) *compute.VirtualMachine {
	vm, err := GetVirtualMachineE(t, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return vm
}

// GetVirtualMachineE gets a Virtual Machine in the specified Azure Resource Group.
func GetVirtualMachineE(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) (*compute.VirtualMachine, error) {
	// Validate resource group name and subscription ID
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	// Get the client reference
	client, err := GetVirtualMachineClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	vm, err := client.Get(context.Background(), resGroupName, vmName, compute.InstanceView)
	if err != nil {
		return nil, err
	}

	return &vm, nil
}

// GetVirtualMachineClientE creates a Azure Virtual Machine client in the specified Azure Subscription.
func GetVirtualMachineClientE(subscriptionID string) (*compute.VirtualMachinesClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the VM client
	client := compute.NewVirtualMachinesClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer

	return &client, nil
}
