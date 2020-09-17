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
	autorest "github.com/Azure/go-autorest/autorest/azure"
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

// ClientFactory describes the methods available on client factory implementatoins
type ClientFactory interface {
	// GetVirtualMachinesClientE returns a configured compute client, setup for proper cloud environment use.
	GetVirtualMachinesClientE(subscriptionID string) (compute.VirtualMachinesClient, error)

	// GetSubscriptionClientE returns a configured compute client, setup for proper cloud environment use.
	GetSubscriptionClientE() (subscriptions.Client, error)

	// GetManagedClustersClientE returns a configured compute client, setup for proper cloud environment use.
	GetManagedClustersClientE(subscriptionID string) (containerservice.ManagedClustersClient, error)
}

// multiEnvClientFactory is used to coordinate handing out properly configured Azure SDK clients
// that are properly setup for use with Public or Sovereign clouds (depending on configuration)
type multiEnvClientFactory struct{}

// NewClientFactory returns a new multi-environment client factory
func NewClientFactory() ClientFactory {
	return &multiEnvClientFactory{}
}

// GetVirtualMachinesClientE returns a configured compute client, setup for proper cloud environment use.
func (factory multiEnvClientFactory) GetVirtualMachinesClientE(subscriptionID string) (compute.VirtualMachinesClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return compute.VirtualMachinesClient{}, err
	}

	// Lookup environment URI
	baseURI, err := factory.getEnvironmentBaseURI()
	if err != nil {
		return compute.VirtualMachinesClient{}, err
	}

	// Create a VM client and return
	vmClient := compute.NewVirtualMachinesClientWithBaseURI(baseURI, subscriptionID)
	return vmClient, nil
}

// GetSubscriptionClientE returns a configured compute client, setup for proper cloud environment use.
func (factory multiEnvClientFactory) GetSubscriptionClientE() (subscriptions.Client, error) {
	// Lookup environment URI
	baseURI, err := factory.getEnvironmentBaseURI()
	if err != nil {
		return subscriptions.Client{}, err
	}

	// Create a Subscription client
	client := subscriptions.NewClientWithBaseURI(baseURI)
	return client, nil
}

// GetManagedClustersClientE returns a configured compute client, setup for proper cloud environment use.
func (factory multiEnvClientFactory) GetManagedClustersClientE(subscriptionID string) (containerservice.ManagedClustersClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return containerservice.ManagedClustersClient{}, err
	}

	// Lookup environment URI
	baseURI, err := factory.getEnvironmentBaseURI()
	if err != nil {
		return containerservice.ManagedClustersClient{}, err
	}

	client := containerservice.NewManagedClustersClientWithBaseURI(baseURI, subscriptionID)
	return client, nil
}

// getDefaultEnvironmentName returns either a configured Azure environment name, or the public default
func (factory multiEnvClientFactory) getDefaultEnvironmentName() string {
	envName, exists := os.LookupEnv(AzureEnvironmentEnvName)

	if !exists || envName == "" {
		envName = autorest.PublicCloud.Name
	}

	return envName
}

// getEnvironmentBaseUri returns the ARM management URI for the configured Azure environment.
func (factory multiEnvClientFactory) getEnvironmentBaseURI() (string, error) {
	envName := factory.getDefaultEnvironmentName()
	env, err := autorest.EnvironmentFromName(envName)
	if err != nil {
		return "", err
	}
	return env.ResourceManagerEndpoint, nil
}

// getKeyVaultURISuffix returns the proper KeyVault URI suffix for the configured Azure environment.
func (factory multiEnvClientFactory) getKeyVaultURISuffix() (string, error) {
	envName := factory.getDefaultEnvironmentName()
	env, err := autorest.EnvironmentFromName(envName)
	if err != nil {
		return "", err
	}
	return env.KeyVaultDNSSuffix, nil
}

// getStorageURISuffix returns the proper storage URI suffix for the configured Azure environment
func (factory multiEnvClientFactory) getStorageURISuffix() (string, error) {
	envName := factory.getDefaultEnvironmentName()
	env, err := autorest.EnvironmentFromName(envName)
	if err != nil {
		return "", err
	}
	return env.StorageEndpointSuffix, nil
}
