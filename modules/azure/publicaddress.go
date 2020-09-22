package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// PublicAddressExists indicates whether the speficied AzurePublic Address exists.
// This function would fail the test if there is an error.
func PublicAddressExists(t testing.TestingT, publicAddressName string, resGroupName string, subscriptionID string) bool {
	exists, err := PublicAddressExistsE(t, publicAddressName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return exists
}

// PublicAddressExistsE indicates whether the speficied AzurePublic Address exists.
func PublicAddressExistsE(t testing.TestingT, publicAddressName string, resGroupName string, subscriptionID string) (bool, error) {
	// Get the Public Address
	_, err := GetPublicIPAddressE(t, publicAddressName, resGroupName, subscriptionID)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetPublicAddressIP gets the IP of a Public IP Address. This function would fail the test if there is an error.
func GetPublicAddressIP(t testing.TestingT, publicAddressName string, resGroupName string, subscriptionID string) string {
	IP, err := GetPublicAddressIPE(t, publicAddressName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return IP
}

// GetPublicAddressIPE gets the IP of a Public IP Address.
func GetPublicAddressIPE(t testing.TestingT, publicAddressName string, resGroupName string, subscriptionID string) (string, error) {
	// Create a NIC client
	pip, err := GetPublicIPAddressE(t, publicAddressName, resGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return *pip.IPAddress, nil
}

// CheckPublicDNSNameAvailability checks whether a Domain Name in the cloudapp.azure.com zone
// is available for use. This function would fail the test if there is an error.
func CheckPublicDNSNameAvailability(t testing.TestingT, location string, domainNameLabel string, subscriptionID string) bool {
	available, err := CheckPublicDNSNameAvailabilityE(t, location, domainNameLabel, subscriptionID)
	if err != nil {
		return false
	}
	return available
}

// CheckPublicDNSNameAvailabilityE checks whether a Domain Name in the cloudapp.azure.com zone is available for use.
func CheckPublicDNSNameAvailabilityE(t testing.TestingT, location string, domainNameLabel string, subscriptionID string) (bool, error) {
	client, err := GetPublicIPAddressClientE(subscriptionID)
	if err != nil {
		return false, err
	}

	res, err := client.CheckDNSNameAvailability(context.Background(), location, domainNameLabel)
	if err != nil {
		return false, err
	}

	return *res.Available, nil
}

// GetPublicIPAddressE gets a Public IP Addresses in the specified Azure Resource Group.
func GetPublicIPAddressE(t testing.TestingT, publicIPAddressName string, resGroupName string, subscriptionID string) (*network.PublicIPAddress, error) {
	// Validate resource group name and subscription ID
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	// Get the client refrence
	client, err := GetPublicIPAddressClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Public IP Address
	pip, err := client.Get(context.Background(), resGroupName, publicIPAddressName, "")
	if err != nil {
		return nil, err
	}
	return &pip, nil
}

// GetPublicIPAddressClientE creates a Public IP Addresses client in the specified Azure Subscription.
func GetPublicIPAddressClientE(subscriptionID string) (*network.PublicIPAddressesClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Public IP Address client
	client := network.NewPublicIPAddressesClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer

	return &client, nil
}
