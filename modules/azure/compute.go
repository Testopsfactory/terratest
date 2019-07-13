package azure

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-06-01/compute"
	"github.com/stretchr/testify/require"
)

// GetSizeOfVirtualMachine gets the size type of the given Azure Virtual Machine
func GetSizeOfVirtualMachine(t *testing.T, vmName string, resGroupName string, subscriptionID string) string {
	size, err := GetSizeOfVirtualMachineE(t, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return size
}

// GetSizeOfVirtualMachineE gets the size type of the given Azure Virtual Machine
func GetSizeOfVirtualMachineE(t *testing.T, vmName string, resGroupName string, subscriptionID string) (string, error) {
	// Validate resource group name and subscription ID
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return "", err
	}

	subscriptionID, err = getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return "", err
	}

	// Create a VM client
	vmClient := compute.NewVirtualMachinesClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return "", err
	}

	// Attach authorizer to the client
	vmClient.Authorizer = *authorizer

	// Get the details of the target virtual machine
	ctx := context.TODO()
	vm, err := vmClient.Get(ctx, resGroupName, vmName, compute.InstanceView)
	if err != nil {
		return "", err
	}

	return string(vm.VirtualMachineProperties.HardwareProfile.VMSize), nil
}

// GetTagsForVirtualMachine gets the tags of the given Virtual Machine as a map
func GetTagsForVirtualMachine(t *testing.T, vmName string, resGroupName string, subscriptionID string) map[string]string {
	tags, err := GetTagsForVirtualMachineE(t, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return tags
}

// GetTagsForVirtualMachineE gets the tags of the given Virtual Machine as a map
func GetTagsForVirtualMachineE(t *testing.T, vmName string, resGroupName string, subscriptionID string) (map[string]string, error) {
	// Setup a blank map to populate and return
	tags := make(map[string]string)

	// Validate resource group name and subscription ID
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return tags, err
	}

	subscriptionID, err = getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return tags, err
	}

	// Create a VM client
	vmClient := compute.NewVirtualMachinesClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return tags, err
	}

	// Attach authorizer to the client
	vmClient.Authorizer = *authorizer

	// Get the details of the target virtual machine
	ctx := context.TODO()
	vm, err := vmClient.Get(ctx, resGroupName, vmName, compute.InstanceView)
	if err != nil {
		return tags, err
	}

	// Range through existing tags and populate above map accordingly
	for k, v := range vm.Tags {
		tags[k] = *v
	}

	return tags, nil
}
