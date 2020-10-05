package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2019-11-01/containerservice"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// GetManagedClustersClientE is a helper function that will setup an Azure ManagedClusters client on your behalf
func GetManagedClustersClientE(subscriptionID string) (*containerservice.ManagedClustersClient, error) {
	// Create a cluster client
	factory := NewClientFactory()
	client, err := factory.GetClientE(ManagedClustersClientType, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Type cast and verify
	clustersClient, ok := client.(containerservice.ManagedClustersClient)
	if ok {
		return nil, fmt.Errorf("Unable to convert client to type containerservice.ManagedClustersClient")
	}

	// setup authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	clustersClient.Authorizer = *authorizer
	return &clustersClient, nil
}

// GetManagedClusterE will return ManagedCluster
func GetManagedClusterE(t testing.TestingT, resourceGroupName, clusterName, subscriptionID string) (*containerservice.ManagedCluster, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}
	client, err := GetManagedClustersClientE(subscriptionID)
	if err != nil {
		return nil, err
	}
	managedCluster, err := client.Get(context.Background(), resourceGroupName, clusterName)
	if err != nil {
		return nil, err
	}
	return &managedCluster, nil
}
