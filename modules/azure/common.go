package azure

import (
	"fmt"
	"os"

	autorest "github.com/Azure/go-autorest/autorest/azure"
)

const (
	// AzureSubscriptionID is an optional env variable supported by the `azurerm` Terraform provider to
	// designate a target Azure subscription ID
	AzureSubscriptionID = "ARM_SUBSCRIPTION_ID"

	// AzureResGroupName is an optional env variable custom to Terratest to designate a target Azure resource group
	AzureResGroupName = "AZURE_RES_GROUP_NAME"

	// AzureEnvironmentEnvName is the name of the Azure environment to use
	// Set to AzureUSGovernmentCloud or AzurePublicCloud
	AzureEnvironmentEnvName = "AZURE_ENVIRONMENT"
)

// GetTargetAzureSubscription is a helper function to find the correct target Azure Subscription ID,
// with provided arguments taking precedence over environment variables
func GetTargetAzureSubscription(subscriptionID string) (string, error) {
	return getTargetAzureSubscription(subscriptionID)
}

func getTargetAzureSubscription(subscriptionID string) (string, error) {
	fmt.Printf("Initial subscription ID is %s\n", subscriptionID)
	if subscriptionID == "" {
		if id, exists := os.LookupEnv(AzureSubscriptionID); exists {
			return id, nil
		}

		return "", SubscriptionIDNotFound{}
	}

	fmt.Printf("Final subscription ID is %s\n", subscriptionID)

	return subscriptionID, nil
}

// GetTargetAzureResourceGroupName is a helper function to find the correct target Azure Resource Group name,
// with provided arguments taking precedence over environment variables
func GetTargetAzureResourceGroupName(resourceGroupName string) (string, error) {
	return getTargetAzureResourceGroupName(resourceGroupName)
}

func getTargetAzureResourceGroupName(resourceGroupName string) (string, error) {
	if resourceGroupName == "" {
		if name, exists := os.LookupEnv(AzureResGroupName); exists {
			return name, nil
		}

		return "", ResourceGroupNameNotFound{}
	}

	return resourceGroupName, nil
}

// getDefaultEnvironmentName returns either a configured Azure environment name, or the public default
func getDefaultEnvironmentName() string {
	envName, exists := os.LookupEnv(AzureEnvironmentEnvName)

	if !exists || envName == "" {
		envName = autorest.PublicCloud.Name
	}

	return envName
}

// getEnvironmentBaseUri returns the ARM management URI for the configured Azure environment.
func getEnvironmentBaseURI() (string, error) {
	envName := getDefaultEnvironmentName()
	env, err := autorest.EnvironmentFromName(envName)
	if err != nil {
		return "", err
	}
	return env.ResourceManagerEndpoint, nil
}

// getKeyVaultURISuffix returns the proper KeyVault URI suffix for the configured Azure environment.
func getKeyVaultURISuffix() (string, error) {
	envName := getDefaultEnvironmentName()
	env, err := autorest.EnvironmentFromName(envName)
	if err != nil {
		return "", err
	}
	return env.KeyVaultDNSSuffix, nil
}

// getStorageURISuffix returns the proper storage URI suffix for the configured Azure environment
func getStorageURISuffix() (string, error) {
	envName := getDefaultEnvironmentName()
	env, err := autorest.EnvironmentFromName(envName)
	if err != nil {
		return "", err
	}
	return env.StorageEndpointSuffix, nil
}
