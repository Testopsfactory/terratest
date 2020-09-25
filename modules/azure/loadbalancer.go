package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// IPType is enumerator for IP Types, e.g. Public/Private
type IPType int

// IPType values listed using iota
const (
	PublicIP IPType = iota
	PrivateIP
)

// string values for IPType
func (ipType IPType) String() string {
	return [...]string{"PublicIP", "PrivateIP"}[ipType]
}

// LoadBalancerExists indicates whether the specified Load Balancer exists.
// This function would fail the test if there is an error.
func LoadBalancerExists(t testing.TestingT, loadBalancerName string, resourceGroupName string, subscriptionID string) bool {
	exists, err := LoadBalancerExistsE(t, loadBalancerName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return exists
}

// LoadBalancerExistsE indicates whether the specified Load Balancer exists.
func LoadBalancerExistsE(t testing.TestingT, loadBalancerName string, resourceGroupName string, subscriptionID string) (bool, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return false, err
	}
	client, err := GetLoadBalancerClientE(t, subscriptionID)
	if err != nil {
		return false, err
	}
	lb, err := client.Get(context.Background(), resourceGroupName, loadBalancerName, "")
	if err != nil {
		return false, err
	}

	return *lb.Name == loadBalancerName, nil
}

// GetLoadBalancer returns a load balancer resource as specified by name, else returns nil with err
// This function would fail the test if there is an error.
func GetLoadBalancer(t testing.TestingT, loadBalancerName string, resourceGroupName string, subscriptionID string) *network.LoadBalancer {
	lb, err := GetLoadBalancerE(t, loadBalancerName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return lb
}

// GetLoadBalancerE returns a load balancer resource as specified by name, else returns nil with err
func GetLoadBalancerE(t testing.TestingT, loadBalancerName string, resourceGroupName string, subscriptionID string) (*network.LoadBalancer, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}
	client, err := GetLoadBalancerClientE(t, subscriptionID)
	if err != nil {
		return nil, err
	}
	lb, err := client.Get(context.Background(), resourceGroupName, loadBalancerName, "")
	if err != nil {
		return nil, err
	}

	return &lb, nil
}

// GetLoadBalancerClientE creates a load balancer client.
func GetLoadBalancerClientE(t testing.TestingT, subscriptionID string) (*network.LoadBalancersClient, error) {
	loadBalancerClient := network.NewLoadBalancersClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	loadBalancerClient.Authorizer = *authorizer
	return &loadBalancerClient, nil
}

// GetLoadBalancerFrontendConfig indicates whether the specified Load Balancer exists.
// This function would fail the test if there is an error.
func GetLoadBalancerFrontendConfig(t testing.TestingT, loadBalancerName string, resourceGroupName string, subscriptionID string) (ipAddress string, publicOrPrivate string) {
	ipAddress, ipType, err := GetLoadBalancerFrontendConfigE(t, loadBalancerName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return ipAddress, ipType
}

// GetLoadBalancerFrontendConfigE returns an IP address and specifies public or private
func GetLoadBalancerFrontendConfigE(t testing.TestingT, loadBalancer01Name string, resourceGroupName string, subscriptionID string) (ipAddress string, publicOrPrivate string, err1 error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return "", "", err
	}

	// Read the LB information
	lb01, err := GetLoadBalancerE(t, loadBalancer01Name, resourceGroupName, "")
	if err != nil {
		return "", "", err
	}
	lb01Props := lb01.LoadBalancerPropertiesFormat
	fe01Config := (*lb01Props.FrontendIPConfigurations)[0]
	fe01Props := *fe01Config.FrontendIPConfigurationPropertiesFormat

	if fe01Props.PrivateIPAddress == nil {
		// Get PublicIPAddressResource name for Load Balancer
		pipResourceName := GetNameFromResourceID(*fe01Props.PublicIPAddress.ID)

		client, err := GetPublicIPAddressClientE(subscriptionID)
		if err != nil {
			return "", "", err
		}
		publicIPAddress, err := client.Get(context.Background(), resourceGroupName, pipResourceName, "")
		if err != nil {
			return "", "", err
		}

		pipProps := *publicIPAddress.PublicIPAddressPropertiesFormat
		ipValue := (pipProps.IPAddress)

		// return public IP
		return *ipValue, string(PublicIP), nil
	} else {
		// return private IP
		return *fe01Props.PrivateIPAddress, string(PrivateIP), nil
	}

}

// GetPublicIPAddressClientE creates a PublicIPAddresses client
func GetPublicIPAddressClientE(subscriptionID string) (*network.PublicIPAddressesClient, error) {
	client := network.NewPublicIPAddressesClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	return &client, nil
}
