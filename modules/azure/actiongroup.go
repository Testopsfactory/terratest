package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/resources/mgmt/insights"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetActionGroupResource gets the AppServicePlan.
// planName - required to find the AppServicePlan.
// resGroupName - use an empty string if you have the AZURE_RES_GROUP_NAME environment variable set
// subscriptionId - use an empty string if you have the ARM_SUBSCRIPTION_ID environment variable set
func GetActionGroupResource(t *testing.TestingT, ruleName string, resGroupName string, subscriptionID string) *insights.ActionGroupResource {
	actionGroupResource, err := getActionGroupResourceE(ruleName)
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

	actionGroup, err := client.Get(context.Background(), rgName, planName)
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

	authorizer, err := azure.NewAuthorizer()
	if err != nil {
		return nil, err
	}

	metricAlertsClient.Authorizer = *authorizer

	return &metricAlertsClient, nil
}
