package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/gruntwork-io/terratest/modules/collections"
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

// LoadBalancerExistsE returns true if the load balancer exists, else returns false with err
func LoadBalancerExistsE(loadBalancerName string, resourceGroupName string, subscriptionID string) (bool, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return false, err
	}
	client, err := GetLoadBalancerClientE(subscriptionID)
	if err != nil {
		return false, err
	}
	lb, err := client.Get(context.Background(), resourceGroupName, loadBalancerName, "")
	if err != nil {
		return false, err
	}

	return *lb.Name == loadBalancerName, nil
}

// GetLoadBalancerE returns a load balancer resource as specified by name, else returns nil with err
func GetLoadBalancerE(loadBalancerName string, resourceGroupName string, subscriptionID string) (*network.LoadBalancer, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}
	client, err := GetLoadBalancerClientE(subscriptionID)
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
func GetLoadBalancerClientE(subscriptionID string) (*network.LoadBalancersClient, error) {
	loadBalancerClient := network.NewLoadBalancersClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	loadBalancerClient.Authorizer = *authorizer
	return &loadBalancerClient, nil
}

// GetLoadBalancerFrontendConfig returns an IP address and specifies public or private
// This function would fail the test if there is an error.
func GetLoadBalancerFrontendConfig(loadBalancer01Name string, resourceGroupName string, subscriptionID string) (ipAddress string, publicOrPrivate string, err1 error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return "", "", err
	}

	// Read the LB information
	lb01, err := GetLoadBalancerE(loadBalancer01Name, resourceGroupName, "")
	if err != nil {
		return "", "", err
	}
	lb01Props := lb01.LoadBalancerPropertiesFormat
	fe01Config := (*lb01Props.FrontendIPConfigurations)[0]
	fe01Props := *fe01Config.FrontendIPConfigurationPropertiesFormat

	if fe01Props.PrivateIPAddress == nil {
		// Get PublicIPAddressResource name for Load Balancer
		pipResourceName, err := collections.GetSliceLastValueE(*fe01Props.PublicIPAddress.ID, "/")
		if err != nil {
			return "", "", err
		}

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
