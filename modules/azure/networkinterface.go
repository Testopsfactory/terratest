package azure

import (
	"context"
	"reflect"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// NetworkInterfaceExists indicates whether the specified Azure Network Interface exists.
// This function would fail the test if there is an error.
func NetworkInterfaceExists(t testing.TestingT, nicName string, resGroupName string, subscriptionID string) bool {
	exists, err := NetworkInterfaceExistsE(t, nicName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return exists
}

// NetworkInterfaceExistsE indicates whether the specified Azure Network Interface exists.
func NetworkInterfaceExistsE(t testing.TestingT, nicName string, resGroupName string, subscriptionID string) (bool, error) {
	// Get the Network Interface
	_, err := GetNetworkInterfaceE(t, nicName, resGroupName, subscriptionID)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetNetworkInterfacePrivateIPs gets a list of the Private IPs of a Network Interface configs.
// This function would fail the test if there is an error.
func GetNetworkInterfacePrivateIPs(t testing.TestingT, nicName string, resGroupName string, subscriptionID string) []string {
	IPs, err := GetNetworkInterfacePrivateIPsE(t, nicName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return IPs
}

// GetNetworkInterfacePrivateIPsE gets a list of the Private IPs of a Network Interface configs.
func GetNetworkInterfacePrivateIPsE(t testing.TestingT, nicName string, resGroupName string, subscriptionID string) ([]string, error) {
	var privateIPs []string

	// Get the Network Interface client
	nic, err := GetNetworkInterfaceE(t, nicName, resGroupName, subscriptionID)
	if err != nil {
		return privateIPs, err
	}

	// Get the Private IPs from each configuration
	for _, IPConfiguration := range *nic.IPConfigurations {
		privateIPs = append(privateIPs, *IPConfiguration.PrivateIPAddress)
	}

	return privateIPs, nil
}

// GetNetworkInterfacePublicIPs returns a list of all the Public IPs found in the Network Interface configurations.
// This function would fail the test if there is an error.
func GetNetworkInterfacePublicIPs(t testing.TestingT, nicName string, resGroupName string, subscriptionID string) []string {
	IPs, err := GetNetworkInterfacePublicIPsE(t, nicName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return IPs
}

// GetNetworkInterfacePublicIPsE returns a list of all the Public IPs found in the Network Interface configurations.
func GetNetworkInterfacePublicIPsE(t testing.TestingT, nicName string, resGroupName string, subscriptionID string) ([]string, error) {
	var publicIPs []string

	// Get the Network Interface client
	nic, err := GetNetworkInterfaceE(t, nicName, resGroupName, subscriptionID)
	if err != nil {
		return publicIPs, err
	}

	// Get the Public IPs from each configuration available
	for _, IPConfiguration := range *nic.IPConfigurations {
		// Iterate each config, for successful configurations check for a Public Address reference.
		// Not failing on errors as as this is an optimistic accumulator.
		nicConfig, err := GetNetworkInterfaceConfigurationE(t, nicName, *IPConfiguration.Name, resGroupName, subscriptionID)
		if err == nil {
			if publicAddressID := nicConfig.PublicIPAddress; !reflect.ValueOf(nicConfig.PublicIPAddress).IsNil() {
				// Get Public Address IP for configs with using Public Address client
				publicIP, err := GetPublicAddressIPE(t, GetNameFromResourceID(*publicAddressID.ID), resGroupName, subscriptionID)
				if err == nil {
					publicIPs = append(publicIPs, publicIP)
				}
			}
		}
	}

	return publicIPs, nil
}

// GetNetworkInterfaceConfigurationE gets a Network Interface Configuration in the specified Azure Resource Group.
func GetNetworkInterfaceConfigurationE(t testing.TestingT, nicName string, nicConfigName string, resGroupName string, subscriptionID string) (*network.InterfaceIPConfiguration, error) {
	// Validate Azure Resource Group
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	// Get the client reference
	client, err := GetNetworkInterfaceConfigurationClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Network Interface
	nicConfig, err := client.Get(context.Background(), resGroupName, nicName, nicConfigName)
	if err != nil {
		return nil, err
	}

	return &nicConfig, nil
}

// GetNetworkInterfaceConfigurationClientE creates a new Network Interface Configuration client in the specified Azure Subscription.
func GetNetworkInterfaceConfigurationClientE(subscriptionID string) (*network.InterfaceIPConfigurationsClient, error) {
	// Validate Azure Subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the NIC client
	client := network.NewInterfaceIPConfigurationsClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer

	return &client, nil
}

// GetNetworkInterfaceE gets a Network Interface in the specified Azure Resource Group.
func GetNetworkInterfaceE(t testing.TestingT, nicName string, resGroupName string, subscriptionID string) (*network.Interface, error) {
	// Validate Azure Resource Group
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	// Get the client reference
	client, err := GetNetworkInterfaceClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Network Interface
	nic, err := client.Get(context.Background(), resGroupName, nicName, "")
	if err != nil {
		return nil, err
	}

	return &nic, nil
}

// GetNetworkInterfaceClientE creates a new Network Interface client in the specified Azure Subscription.
func GetNetworkInterfaceClientE(subscriptionID string) (*network.InterfacesClient, error) {
	// Validate Azure Subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the NIC client
	client := network.NewInterfacesClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer

	return &client, nil
}
