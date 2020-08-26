package azure

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/resources/mgmt/insights"
	"github.com/stretchr/testify/require"
)

// GetActionGroupResource gets the AppServicePlan.
// ruleName - required to find the AppServicePlan.
// resGroupName - use an empty string if you have the AZURE_RES_GROUP_NAME environment variable set
// subscriptionId - use an empty string if you have the ARM_SUBSCRIPTION_ID environment variable set
func GetActionGroupResource(t *testing.T, ruleName string, resGroupName string, subscriptionID string) *insights.ActionGroupResource {
	actionGroupResource, err := getActionGroupResourceE(ruleName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return actionGroupResource
}

func getActionGroupResourceE(ruleName string, resGroupName string, subscriptionID string) (*insights.ActionGroupResource, error) {
	rgName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := getActionGroupClient(subscriptionID)
	if err != nil {
		return nil, err
	}

	actionGroup, err := client.Get(context.Background(), rgName, ruleName)
	if err != nil {
		return nil, err
	}

	return &actionGroup, nil
}

func getActionGroupClient(subscriptionID string) (*insights.ActionGroupsClient, error) {
	subID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	metricAlertsClient := insights.NewActionGroupsClient(subID)

	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	metricAlertsClient.Authorizer = *authorizer

	return &metricAlertsClient, nil
}
