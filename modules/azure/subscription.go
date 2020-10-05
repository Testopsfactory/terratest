package azure

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-06-01/subscriptions"
)

// GetSubscriptionClient is a helper function that will setup an Azure Subscription client on your behalf
func GetSubscriptionClient() (*subscriptions.Client, error) {
	// Create a Subscription client
	factory := NewClientFactory()
	client, err := factory.GetClientE(SubscriptionsClientType, "")
	if err != nil {
		return nil, err
	}

	// type cast and verify
	subscriptionClient, ok := client.(subscriptions.Client)
	if !ok {
		return nil, fmt.Errorf("Unable to convert client to subscriptions.Client")
	}

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	subscriptionClient.Authorizer = *authorizer
	return &subscriptionClient, nil
}
