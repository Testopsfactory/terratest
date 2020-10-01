package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// VirtualMachineExists indicates whether the specifcied Azure Virtual Machine exists.
// This function would fail the test if there is an error.
func VirtualMachineExists(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) bool {
	exists, err := VirtualMachineExistsE(vmName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return exists
}

// VirtualMachineExistsE indicates whether the specifcied Azure Virtual Machine exists.
func VirtualMachineExistsE(vmName string, resGroupName string, subscriptionID string) (bool, error) {
	// Get VM Object
	_, err := GetVirtualMachineE(vmName, resGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetVirtualMachineNics gets a list of Network Interface names for a specifcied Azure Virtual Machine.
// This function would fail the test if there is an error.
func GetVirtualMachineNics(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) []string {
	nicList, err := GetVirtualMachineNicsE(vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return nicList
}

// GetVirtualMachineNicsE gets a list of Network Interface names for a specified Azure Virtual Machine.
func GetVirtualMachineNicsE(vmName string, resGroupName string, subscriptionID string) ([]string, error) {

	// Get VM Object
	vm, err := GetVirtualMachineE(vmName, resGroupName, subscriptionID)
	if err != nil {
		return nil, err
	}

	vmNICs := *vm.NetworkProfile.NetworkInterfaces
	if len(vmNICs) == 0 {
		// No NIC present
		return nil, nil
	}

	// Get the Names of the attached NICs
	nics := make([]string, len(vmNICs))

	for i, nic := range vmNICs {
		nics[i] = GetNameFromResourceID(*nic.ID)
	}
	return nics, nil
}

// GetVirtualMachineNicCount gets the Network Interface count of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
func GetVirtualMachineNicCount(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) int {
	nicCount, err := GetVirtualMachineNicCountE(vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return nicCount
}

// GetVirtualMachineNicCountE gets the Network Interface count of the specified Azure Virtual Machine.
func GetVirtualMachineNicCountE(vmName string, resGroupName string, subscriptionID string) (int, error) {
	// Get VM Object
	vm, err := GetVirtualMachineE(vmName, resGroupName, subscriptionID)
	if err != nil {
		return -1, err
	}

	return len(*vm.NetworkProfile.NetworkInterfaces), nil
}

// GetVirtualMachineManagedDisks gets the list of Managed Disk names of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
func GetVirtualMachineManagedDisks(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) []string {
	diskNames, err := GetVirtualMachineManagedDisksE(vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return diskNames
}

// GetVirtualMachineManagedDisksE gets the list of Managed Disk names of the specified Azure Virtual Machine.
func GetVirtualMachineManagedDisksE(vmName string, resGroupName string, subscriptionID string) ([]string, error) {

	// Get VM Object
	vm, err := GetVirtualMachineE(vmName, resGroupName, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get VM attached Disks
	vmDisks := *vm.StorageProfile.DataDisks

	// No Attached Disks present
	if len(vmDisks) == 0 {
		return nil, nil
	}

	// Get the Names of the attached Managed Disks
	diskNames := make([]string, len(vmDisks))
	for _, v := range vmDisks {
		diskNames = append(diskNames, *v.Name)
	}

	return diskNames, nil
}

// GetVirtualMachineManagedDiskCount gets the Managed Disk count of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
func GetVirtualMachineManagedDiskCount(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) int {
	mngDiskCount, err := GetVirtualMachineManagedDiskCountE(vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return mngDiskCount
}

// GetVirtualMachineManagedDiskCountE gets the Managed Disk count of the specified Azure Virtual Machine.
func GetVirtualMachineManagedDiskCountE(vmName string, resGroupName string, subscriptionID string) (int, error) {
	// Get VM Object
	vm, err := GetVirtualMachineE(vmName, resGroupName, subscriptionID)
	if err != nil {
		return -1, err
	}

	return len(*vm.StorageProfile.DataDisks), nil
}

// GetVirtualMachineOsDiskName gets the OS Disk name of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
func GetVirtualMachineOsDiskName(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) string {
	osDiskName, err := GetVirtualMachineOsDiskNameE(vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return osDiskName
}

// GetVirtualMachineOsDiskNameE gets the OS Disk name of the specified Azure Virtual Machine.
func GetVirtualMachineOsDiskNameE(vmName string, resGroupName string, subscriptionID string) (string, error) {
	// Get VM Object
	vm, err := GetVirtualMachineE(vmName, resGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return *vm.StorageProfile.OsDisk.Name, nil
}

// GetVirtualMachineState gets the State of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
func GetVirtualMachineState(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) string {
	vmState, err := GetVirtualMachineStateE(vmName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return vmState
}

// GetVirtualMachineStateE gets the State of the specified Azure Virtual Machine.
func GetVirtualMachineStateE(vmName string, resGroupName string, subscriptionID string) (string, error) {
	// Get VM Object
	vm, err := GetVirtualMachineE(vmName, resGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return vm.Status, nil
}

// GetVirtualMachineAvailabilitySetID gets the Availability Set ID of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
func GetVirtualMachineAvailabilitySetID(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) string {
	avsID, err := GetVirtualMachineAvailabilitySetIDE(vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return avsID
}

// GetVirtualMachineAvailabilitySetIDE gets the Availability Set ID of the specified Azure Virtual Machine.
func GetVirtualMachineAvailabilitySetIDE(vmName string, resGroupName string, subscriptionID string) (string, error) {
	// Get VM Object
	vm, err := GetVirtualMachineE(vmName, resGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	// Virtual Machine has no associated Availability Set
	if vm.AvailabilitySet == nil {
		return "", nil
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
	vmImage, err := GetVirtualMachineImageE(vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return vmImage
}

// GetVirtualMachineImageE gets the Image  of the specified Azure Virtual Machine.
func GetVirtualMachineImageE(vmName string, resGroupName string, subscriptionID string) (VMImage, error) {
	var vmImage VMImage

	// Get VM Object
	vm, err := GetVirtualMachineE(vmName, resGroupName, subscriptionID)
	if err != nil {
		return vmImage, err
	}

	if vm.StorageProfile == nil {
		return vmImage, NewNotFoundError("Image Reference", "Any", vmName)
	}

	// Populate VM Image; values always present, no nil checks needed
	vmImage.Publisher = *vm.StorageProfile.ImageReference.Publisher
	vmImage.Offer = *vm.StorageProfile.ImageReference.Offer
	vmImage.SKU = *vm.StorageProfile.ImageReference.Sku
	vmImage.Version = *vm.StorageProfile.ImageReference.Version

	return vmImage, nil
}

// GetVirtualMachineAdminUser gets the Admin Username of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
func GetVirtualMachineAdminUser(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) string {
	adminUser, err := GetVirtualMachineAdminUserE(vmName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return adminUser
}

// GetVirtualMachineAdminUserE gets the Admin Username of the specified Azure Virtual Machine.
func GetVirtualMachineAdminUserE(vmName string, resGroupName string, subscriptionID string) (string, error) {
	// Get VM Object
	vm, err := GetVirtualMachineE(vmName, resGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return string(*vm.OsProfile.AdminUsername), nil
}

// GetVirtualMachineSize gets the Size Type of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
func GetVirtualMachineSize(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) compute.VirtualMachineSizeTypes {
	size, err := GetVirtualMachineSizeE(vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return size
}

// GetVirtualMachineSizeE gets the Size Type of the specified Azure Virtual Machine.
func GetVirtualMachineSizeE(vmName string, resGroupName string, subscriptionID string) (compute.VirtualMachineSizeTypes, error) {
	// Get VM Object
	vm, err := GetVirtualMachineE(vmName, resGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return vm.VirtualMachineProperties.HardwareProfile.VMSize, nil
}

// GetVirtualMachineTags gets the Tags of the specified Virtual Machine as a map.
// This function would fail the test if there is an error.
func GetVirtualMachineTags(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) map[string]string {
	tags, err := GetVirtualMachineTagsE(vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return tags
}

// GetVirtualMachineTagsE gets the Tags of the specified Virtual Machine as a map.
func GetVirtualMachineTagsE(vmName string, resGroupName string, subscriptionID string) (map[string]string, error) {
	// Setup a blank map to populate and return
	var tags map[string]string

	// Get VM Object
	vm, err := GetVirtualMachineE(vmName, resGroupName, subscriptionID)
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

// ListVirtualMachinesForResourceGroup gets a list of all Virtual Machine names in the specified Resource Group.
// This function would fail the test if there is an error.
func ListVirtualMachinesForResourceGroup(t testing.TestingT, resGroupName string, subscriptionID string) []string {
	vms, err := ListVirtualMachinesForResourceGroupE(resGroupName, subscriptionID)
	require.NoError(t, err)
	return vms
}

// ListVirtualMachinesForResourceGroupE gets a list of all Virtual Machine names in the specified Resource Group.
func ListVirtualMachinesForResourceGroupE(resourceGroupName string, subscriptionID string) ([]string, error) {
	var vmDetails []string

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

// GetVirtualMachinesForResourceGroup gets all Virtual Machine objects in the specified Resource Group.
// This function would fail the test if there is an error.
func GetVirtualMachinesForResourceGroup(t testing.TestingT, resGroupName string, subscriptionID string) *map[string]compute.VirtualMachineProperties {
	vms, err := GetVirtualMachinesForResourceGroupE(resGroupName, subscriptionID)
	require.NoError(t, err)
	return vms
}

// GetVirtualMachinesForResourceGroupE gets all Virtual Machine objects in the specified Resource Group.
func GetVirtualMachinesForResourceGroupE(resourceGroupName string, subscriptionID string) (*map[string]compute.VirtualMachineProperties, error) {
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
	vm, err := GetVirtualMachineInstanceE(vmName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return vm
}

// GetVirtualMachineInstanceE gets a local Virtual Machine instance in the specified Resource Group.
func GetVirtualMachineInstanceE(vmName string, resGroupName string, subscriptionID string) (*Instance, error) {
	// Get VM Object
	vm, err := GetVirtualMachineE(vmName, resGroupName, subscriptionID)
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
	vm, err := GetVirtualMachineE(vmName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return vm
}

// GetVirtualMachineE gets a Virtual Machine in the specified Azure Resource Group.
func GetVirtualMachineE(vmName string, resGroupName string, subscriptionID string) (*compute.VirtualMachine, error) {
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
