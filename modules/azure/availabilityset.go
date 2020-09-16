package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/gruntwork-io/terratest/modules/collections"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// AvailabilitySetExists indicates whether the speficied Azure Availability Set exists
func AvailabilitySetExists(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) bool {
	exists, err := AvailabilitySetExistsE(t, avsName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return exists
}

// AvailabilitySetExistsE indicates whether the speficied Azure Availability Set exists
func AvailabilitySetExistsE(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) (bool, error) {
	_, err := GetAvailabilitySetE(t, avsName, resGroupName, subscriptionID)
	if err != nil {
		return false, err
	}
	return true, nil
}

// CheckAvailabilitySetContainsVM checks if the Virtual Machine is contained in the Availability Set VMs
func CheckAvailabilitySetContainsVM(t testing.TestingT, vmName string, avsName string, resGroupName string, subscriptionID string) bool {
	avsVMs, err := GetAvailabilitySetVMsE(t, avsName, resGroupName, subscriptionID)
	if err != nil {
		return false
	}

	return collections.ListContains(avsVMs, strings.ToLower(vmName))
}

// GetAvailabilitySetVMs gets a list of VM names in the specified Azure Availability Set
func GetAvailabilitySetVMs(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) []string {
	vms, err := GetAvailabilitySetVMsE(t, avsName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return vms
}

// GetAvailabilitySetVMsE gets a list of VM names in the specified Azure Availability Set
func GetAvailabilitySetVMsE(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) ([]string, error) {
	client, err := GetAvailabilitySetClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	avs, err := client.Get(context.Background(), resGroupName, avsName)
	if err != nil {
		return nil, err
	}

	vms := []string{}

	for _, vm := range *avs.VirtualMachines {
		vms = append(vms, strings.ToLower(GetNameFromResourceID(*vm.ID)))
	}

	return vms, nil
}

// GetAvailabilitySetFaultDomainCount gets the Fault Domain Count for the specified Azure Availability Set
func GetAvailabilitySetFaultDomainCount(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) int32 {
	avsFaultDomainCount, err := GetAvailabilitySetFaultDomainCountE(t, avsName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return avsFaultDomainCount
}

// GetAvailabilitySetFaultDomainCountE gets the Fault Domain Count for the specified Azure Availability Set
func GetAvailabilitySetFaultDomainCountE(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) (int32, error) {
	client, err := GetAvailabilitySetClientE(subscriptionID)
	if err != nil {
		return -1, err
	}

	avs, err := client.Get(context.Background(), resGroupName, avsName)
	if err != nil {
		return -1, err
	}

	return *avs.PlatformFaultDomainCount, nil
}

// GetAvailabilitySetE gets an Availability Set in the specified Azure Resource Group
func GetAvailabilitySetE(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) (*compute.AvailabilitySet, error) {
	// Validate resource group name and subscription ID
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	// Get the client refrence
	client, err := GetAvailabilitySetClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Availability Set
	avs, err := client.Get(context.Background(), resGroupName, avsName)
	if err != nil {
		return nil, err
	}

	return &avs, nil
}

// GetAvailabilitySetClientE gets a new Availability Set client in the specified Azure Subscription
func GetAvailabilitySetClientE(subscriptionID string) (*compute.AvailabilitySetsClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Availability Set client
	client := compute.NewAvailabilitySetsClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer

	return &client, nil
}
