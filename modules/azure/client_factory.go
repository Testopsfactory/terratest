/*

This file implements an Azure client factory that automatically handles setting up Base URI
values for sovereign cloud support. Note the list of clients below is not initially exhaustive;
rather, additional clients will me added as-needed.

*/

package azure

import (
	"os"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2019-11-01/containerservice"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-06-01/subscriptions"
	autorestAzure "github.com/Azure/go-autorest/autorest/azure"
)

const (
	// AzureEnvironmentEnvName is the name of the Azure environment to use. Set to one of the following:
	//
	// "AzureChinaCloud":        ChinaCloud
	// "AzureGermanCloud":       GermanCloud
	// "AzurePublicCloud":       PublicCloud
	// "AzureUSGovernmentCloud": USGovernmentCloud
	// "AzureStackCloud":		 Azure stack
	AzureEnvironmentEnvName = "AZURE_ENVIRONMENT"
)

// ClientType describes the type of client a module can create.
type ClientType int

// GetClientForSubscriptionsE returns a client instance based on the ClientType passed, or optionally an error.
func GetClientForSubscriptionsE() (subscriptions.Client, error) {
	// Lookup environment URI
	baseURI, err := getEnvironmentBaseURIE()
	if err != nil {
		return subscriptions.Client{}, err
	}

	// Create correct client based on type passed
	return subscriptions.NewClientWithBaseURI(baseURI), nil
}

// GetClientForVirtualMachinesE returns a client instance based on the ClientType passed, or optionally an error.
func GetClientForVirtualMachinesE(subscriptionID string) (compute.VirtualMachinesClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return compute.VirtualMachinesClient{}, err
	}

	// Lookup environment URI
	baseURI, err := getEnvironmentBaseURIE()
	if err != nil {
		return compute.VirtualMachinesClient{}, err
	}

	// Create correct client based on type passed
	return compute.NewVirtualMachinesClientWithBaseURI(baseURI, subscriptionID), nil
}

// GetClientForManagedClustersE returns a client instance based on the ClientType passed, or optionally an error.
func GetClientForManagedClustersE(subscriptionID string) (containerservice.ManagedClustersClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return containerservice.ManagedClustersClient{}, err
	}

	// Lookup environment URI
	baseURI, err := getEnvironmentBaseURIE()
	if err != nil {
		return containerservice.ManagedClustersClient{}, err
	}

	// Create correct client based on type passed
	return containerservice.NewManagedClustersClientWithBaseURI(baseURI, subscriptionID), nil
}

// getDefaultEnvironmentName returns either a configured Azure environment name, or the public default
func getDefaultEnvironmentName() string {
	envName, exists := os.LookupEnv(AzureEnvironmentEnvName)

	if exists && len(envName) > 0 {
		return envName
	}

	return autorestAzure.PublicCloud.Name
}

// getEnvironmentBaseUriE returns the ARM management URI for the configured Azure environment.
func getEnvironmentBaseURIE() (string, error) {
	envName := getDefaultEnvironmentName()
	env, err := autorestAzure.EnvironmentFromName(envName)
	if err != nil {
		return "", err
	}
	return env.ResourceManagerEndpoint, nil
}

// getKeyVaultURISuffixE returns the proper KeyVault URI suffix for the configured Azure environment.
func getKeyVaultURISuffixE() (string, error) {
	envName := getDefaultEnvironmentName()
	env, err := autorestAzure.EnvironmentFromName(envName)
	if err != nil {
		return "", err
	}
	return env.KeyVaultDNSSuffix, nil
}

// getStorageURISuffixE returns the proper storage URI suffix for the configured Azure environment
func getStorageURISuffixE() (string, error) {
	envName := getDefaultEnvironmentName()
	env, err := autorestAzure.EnvironmentFromName(envName)
	if err != nil {
		return "", err
	}
	return env.StorageEndpointSuffix, nil
}
