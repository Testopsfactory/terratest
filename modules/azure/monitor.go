package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/resources/mgmt/insights"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// ActionGroupExists indicates whether the speficied Azure Availability Set exists
func ActionGroupExists(t testing.TestingT, actionGroupName string, resGroupName string, subscriptionID string) bool {
	exists, err := ActionGroupExistsE(t, actionGroupName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return exists
}

// ActionGroupExistsE indicates whether the speficied Azure Availability Set exists
func ActionGroupExistsE(t testing.TestingT, actionGroupName string, resGroupName string, subscriptionID string) (bool, error) {
	_, err := GetActionGroupE(t, actionGroupName, resGroupName, subscriptionID)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetActionGroupE gets a Load Balancer in the specified Azure Resource Group
func GetActionGroupE(t testing.TestingT, actionGroupName string, resGroupName string, subscriptionID string) (*insights.ActionGroupResource, error) {
	// Validate resource group name and subscription ID
	_, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	// Get the client refrence
	client, err := GetActionGroupsClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Action Group
	actionGroup, err := client.Get(context.Background(), resGroupName, actionGroupName)
	if err != nil {
		return nil, err
	}

	return &actionGroup, nil
}

// GetActionGroupsClientE gets an Action Groups client in the specified Azure Subscription
func GetActionGroupsClientE(subscriptionID string) (*insights.ActionGroupsClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Load Balancer client
	client := insights.NewActionGroupsClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer

	return &client, nil
}
