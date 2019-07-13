package azure

import "os"

const (
	// AzuSubID is an optional env variable custom to Terratest to designate a target Azure subscription ID
	AzuSubID = "AZURE_SUB_ID"

	// AzuResGroupName is an optional env variable custom to Terratest to designate a target Azure resource group
	AzuResGroupName = "AZURE_RG_NAME"
)

// getTargetAzureSubscription is a helper function to find the correct target Azure Subscription ID,
// with provided arguments taking precedence over environment variables
func getTargetAzureSubscription(subID string) (string, error) {
	if subID == "" {
		if id, exists := os.LookupEnv(AzuSubID); exists {
			return id, nil
		}

		return "", SubscriptionIDNotFound{}
	}

	return subID, nil
}

// getTargetAzureResourceGroupName is a helper function to find the correct target Azure Resource Group name,
// with provided arguments taking precedence over environment variables
func getTargetAzureResourceGroupName(rgName string) (string, error) {
	if rgName == "" {
		if name, exists := os.LookupEnv(AzuResGroupName); exists {
			return name, nil
		}

		return "", ResourceGroupNameNotFound{}
	}

	return rgName, nil
}
