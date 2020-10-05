// /*

// This file contains unit tests for the client factory implementation(s).

// */

package azure

import (
	"os"
	"testing"

	autorest "github.com/Azure/go-autorest/autorest/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Local consts for this file only
var GovCloudEnvName = "AzureUSGovernmentCloud"
var PublicCloudEnvName = "AzurePublicCloud"
var ChinaCloudEnvName = "AzureChinaCloud"
var GermanyCloudEnvName = "AzureGermanCloud"

func TestDefaultEnvIsPublicWhenNotSet(t *testing.T) {
	// Get a client factory
	factory := multiEnvClientFactory{}

	// save any current env value and restore on exit
	currentEnv := os.Getenv(AzureEnvironmentEnvName)
	defer os.Setenv(AzureEnvironmentEnvName, currentEnv)

	// Set env var to gov
	os.Setenv(AzureEnvironmentEnvName, "")

	// get the default
	env := factory.getDefaultEnvironmentName()

	// Make sure it's public cloud
	assert.Equal(t, autorest.PublicCloud.Name, env)
}

func TestDefaultEnvSetToGov(t *testing.T) {
	// Get a client factory
	factory := multiEnvClientFactory{}

	// save any current env value and restore on exit
	currentEnv := os.Getenv(AzureEnvironmentEnvName)
	defer os.Setenv(AzureEnvironmentEnvName, currentEnv)

	// Set env var to gov
	os.Setenv(AzureEnvironmentEnvName, GovCloudEnvName)

	// get the default
	env := factory.getDefaultEnvironmentName()

	// Make sure it's public cloud
	assert.Equal(t, autorest.USGovernmentCloud.Name, env)
}

func TestClientsBaseURISetCorrectly(t *testing.T) {
	var cases = []struct {
		EnvironmentName string
		Client          ClientType
		ExpectedBaseURI string
	}{
		{GovCloudEnvName, SubscriptionsClientType, autorest.USGovernmentCloud.ResourceManagerEndpoint},
		{GovCloudEnvName, VirtualMachinesClientType, autorest.USGovernmentCloud.ResourceManagerEndpoint},
		{GovCloudEnvName, ManagedClustersClientType, autorest.USGovernmentCloud.ResourceManagerEndpoint},

		{PublicCloudEnvName, SubscriptionsClientType, autorest.PublicCloud.ResourceManagerEndpoint},
		{PublicCloudEnvName, VirtualMachinesClientType, autorest.PublicCloud.ResourceManagerEndpoint},
		{PublicCloudEnvName, ManagedClustersClientType, autorest.PublicCloud.ResourceManagerEndpoint},

		{ChinaCloudEnvName, SubscriptionsClientType, autorest.ChinaCloud.ResourceManagerEndpoint},
		{ChinaCloudEnvName, VirtualMachinesClientType, autorest.ChinaCloud.ResourceManagerEndpoint},
		{ChinaCloudEnvName, ManagedClustersClientType, autorest.ChinaCloud.ResourceManagerEndpoint},

		{GermanyCloudEnvName, SubscriptionsClientType, autorest.GermanCloud.ResourceManagerEndpoint},
		{GermanyCloudEnvName, VirtualMachinesClientType, autorest.GermanCloud.ResourceManagerEndpoint},
		{GermanyCloudEnvName, ManagedClustersClientType, autorest.GermanCloud.ResourceManagerEndpoint},
	}

	// Get a client factory
	factory := multiEnvClientFactory{}

	// save any current env value and restore on exit
	currentEnv := os.Getenv(AzureEnvironmentEnvName)
	defer os.Setenv(AzureEnvironmentEnvName, currentEnv)

	for _, tt := range cases {
		t.Run(tt.EnvironmentName, func(t *testing.T) {
			// Override env setting
			os.Setenv(AzureEnvironmentEnvName, tt.EnvironmentName)

			// Get a VM client
			client, _ := factory.GetClientE(tt.Client, "")

			// Check for correct ARM URI
			baseURI, err := factory.GetClientBaseURIE(tt.Client, client)
			require.NoError(t, err)
			assert.Equal(t, tt.ExpectedBaseURI, baseURI)
		})
	}
}
