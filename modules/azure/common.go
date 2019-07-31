package azure

import "os"

const (
	// AzureSubscriptionID is an optional env variable custom to Terratest to designate a target Azure subscription ID
	AzureSubscriptionID = "AZURE_SUB_ID"

	// AzureResGroupName is an optional env variable custom to Terratest to designate a target Azure resource group
	AzureResGroupName = "AZURE_RES_GROUP_NAME"
)

// getTargetAzureSubscription is a helper function to find the correct target Azure Subscription ID,
// with provided arguments taking precedence over environment variables
func getTargetAzureSubscription(subscriptionID string) (string, error) {
	if subscriptionID == "" {
		if id, exists := os.LookupEnv(AzureSubscriptionID); exists {
			return id, nil
		}

		return "", SubscriptionIDNotFound{}
	}

	return subscriptionID, nil
}

// getTargetAzureResourceGroupName is a helper function to find the correct target Azure Resource Group name,
// with provided arguments taking precedence over environment variables
func getTargetAzureResourceGroupName(resourceGroupName string) (string, error) {
	if resourceGroupName == "" {
		if name, exists := os.LookupEnv(AzureResGroupName); exists {
			return name, nil
		}

		return "", ResourceGroupNameNotFound{}
	}

	return resourceGroupName, nil
}
