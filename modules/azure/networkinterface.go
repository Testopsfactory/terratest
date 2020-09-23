package azure

import (
	"context"
	"encoding/json"
	"strings"

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

// GetNetworkInterfacePublicIPs returns a list of all the Public IPs found in the Network Interface configurations.
// This function would fail the test if there is an error.
func GetNetworkInterfacePublicIPs(t testing.TestingT, nicName string, resGroupName string, subscriptionID string) []string {
	IPs, err := GetNetworkInterfacePublicIPsE(t, nicName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return IPs
}

// GetNetworkInterfacePublicIPsE returns a list of all the Public IPs found in the Network Interface configurations.
func GetNetworkInterfacePublicIPsE(t testing.TestingT, nicName string, resGroupName string, subscriptionID string) ([]string, error) {
	publicIPs := []string{}

	// Get the Network Interface client
	nic, err := GetNetworkInterfaceE(t, nicName, resGroupName, subscriptionID)
	if err != nil {
		return publicIPs, err
	}

	// Get the Public IPs from each configuration available
	for _, IPConfiguration := range *nic.IPConfigurations {
		if ipConfigHasPublicIP(&IPConfiguration) {
			// Get the ID from the long string NIC representation
			publicAddressID := GetNameFromResourceID(*IPConfiguration.PublicIPAddress.ID)

			// Get the Public Ip from the Public Address client
			publicIP := GetPublicAddressIP(t, publicAddressID, resGroupName, subscriptionID)
			publicIPs = append(publicIPs, publicIP)
		}
	}

	return publicIPs, nil
}

// ipConfigHasPublicIP returns true if an IP Configuration has a Public IP Address.
// This helper method was created since a config without a public address causes a nil pointer panic
// and the string representation is searched for the publicIPAddress text to identify it's presence.
func ipConfigHasPublicIP(ipConfig *network.InterfaceIPConfiguration) bool {
	var byteIPConfig []byte

	byteIPConfig, err := json.Marshal(ipConfig)
	if err != nil {
		return false
	}

	return strings.Contains(string(byteIPConfig), "publicIPAddress")
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
	privateIPs := []string{}

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
