package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// IPType enumerator for IP Types, Public or Private.
type IPType int

// IPType values listed using iota
const (
	PublicIP IPType = iota
	PrivateIP
)

// String values for IPType
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
	_, err := GetLoadBalancerE(t, loadBalancerName, resourceGroupName, subscriptionID)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetLoadBalancerConfigNames gets a list of the Frontend Configuration Names for the Load Balancer.
// This function would fail the test if there is an error.
func GetLoadBalancerConfigNames(t testing.TestingT, loadBalancerName string, resourceGroupName string, subscriptionID string) []string {
	configName, err := GetLoadBalancerConfigNamesE(t, loadBalancerName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return configName
}

// GetLoadBalancerConfigNamesE gets a list of the Frontend Configuration Names for the Load Balancer.
func GetLoadBalancerConfigNamesE(t testing.TestingT, loadBalancerName string, resourceGroupName string, subscriptionID string) ([]string, error) {
	var configNames []string

	lb, err := GetLoadBalancerE(t, loadBalancerName, resourceGroupName, subscriptionID)
	if err != nil {
		return configNames, err
	}

	lbProps := lb.LoadBalancerPropertiesFormat
	feConfigs := *lbProps.FrontendIPConfigurations
	if len(feConfigs) == 0 {
		// No Frontend Configurations present
		return configNames, NewNotFoundError("Frontend Config", "Any", loadBalancerName)
	}

	// Get the names of the Frontend IP Configurations present
	for _, config := range feConfigs {
		configNames = append(configNames, *config.Name)
	}

	return configNames, nil
}

// GetLoadBalancerFrontendConfig gets the IP address and specifies public or private for the specified Load Balancer.
// This function would fail the test if there is an error.
func GetLoadBalancerFrontendConfig(t testing.TestingT, feConfigName string, loadBalancerName string, resourceGroupName string, subscriptionID string) (ipAddress string, publicOrPrivate IPType) {
	ipAddress, ipType, err := GetLoadBalancerFrontendConfigE(t, feConfigName, loadBalancerName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return ipAddress, ipType
}

// GetLoadBalancerFrontendConfigE gets the IP address and specifies public or private for the specified Load Balancer.
func GetLoadBalancerFrontendConfigE(t testing.TestingT, feConfigName string, loadBalancer01Name string, resourceGroupName string, subscriptionID string) (ipAddress string, publicOrPrivate IPType, err1 error) {
	// Get the specified Load Balancer Frontend Config
	feConfig, err := GetLoadBalancerFrontendIPConfigurationE(t, feConfigName, loadBalancer01Name, resourceGroupName, subscriptionID)
	if err != nil {
		return "", -1, err
	}

	// Get the Properties of the Frontend Configuration
	feProps := *feConfig.FrontendIPConfigurationPropertiesFormat

	// Check Public Type
	if feProps.PublicIPAddress != nil {
		// Get PublicIPAddressResource name for Load Balancer Frontend Configuration
		pipResourceName := GetNameFromResourceID(*feProps.PublicIPAddress.ID)

		client, err := GetPublicIPAddressClientE(subscriptionID)
		if err != nil {
			return "", -1, err
		}
		publicIPAddress, err := client.Get(context.Background(), resourceGroupName, pipResourceName, subscriptionID)
		if err != nil {
			return "", -1, err
		}

		pipProps := *publicIPAddress.PublicIPAddressPropertiesFormat
		ipValue := (pipProps.IPAddress)

		// Public IP
		return *ipValue, IPType(PublicIP), nil
	}

	// Private IP
	return *feProps.PrivateIPAddress, IPType(PrivateIP), nil

}

// GetLoadBalancerFrontendIPConfigurationE gets a Load Balancer Configuration in the specified Azure Resource Group
func GetLoadBalancerFrontendIPConfigurationE(t testing.TestingT, feConfigName string, loadBalancerName string, resourceGroupName string, subscriptionID string) (*network.FrontendIPConfiguration, error) {
	// Validate Azure Resource Group Name
	resourceGroupName, err := getTargetAzureResourceGroupName(resourceGroupName)
	if err != nil {
		return nil, err
	}

	// Get the client reference
	client, err := GetLoadBalancerFrontendIPConfigurationClientE(t, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Load Balancer Frontend Configuration
	lbc, err := client.Get(context.Background(), resourceGroupName, loadBalancerName, feConfigName)
	if err != nil {
		return nil, err
	}

	return &lbc, nil
}

// GetLoadBalancerFrontendIPConfigurationClientE creates a Load Balancer Configuration client.
func GetLoadBalancerFrontendIPConfigurationClientE(t testing.TestingT, subscriptionID string) (*network.LoadBalancerFrontendIPConfigurationsClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Load Balancer Frontend Configuration client
	client := network.NewLoadBalancerFrontendIPConfigurationsClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer

	return &client, nil
}

// GetLoadBalancerE gets a Load Balancer in the specified Azure Resource Group
func GetLoadBalancerE(t testing.TestingT, loadBalancerName string, resourceGroupName string, subscriptionID string) (*network.LoadBalancer, error) {
	// Validate Azure Resource Group Name
	resourceGroupName, err := getTargetAzureResourceGroupName(resourceGroupName)
	if err != nil {
		return nil, err
	}

	// Get the client reference
	client, err := GetLoadBalancerClientE(t, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Load Balancer
	lb, err := client.Get(context.Background(), resourceGroupName, loadBalancerName, subscriptionID)
	if err != nil {
		return nil, err
	}

	return &lb, nil
}

// GetLoadBalancerClientE gets a new Load Balancer client in the specified Azure Subscription.
func GetLoadBalancerClientE(t testing.TestingT, subscriptionID string) (*network.LoadBalancersClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Load Balancer client
	client := network.NewLoadBalancersClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer

	return &client, nil
}
