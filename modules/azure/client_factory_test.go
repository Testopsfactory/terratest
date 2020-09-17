/*

This file contains unit tests for the client factory implementation(s).

*/

package azure

import (
	"os"
	"testing"

	autorest "github.com/Azure/go-autorest/autorest/azure"
	"github.com/stretchr/testify/assert"
)

// Local consts for this file only
var GovCloudEnvName = "AzureUSGovernmentCloud"
var PublicCloudEnvName = "AzurePublicCloud"

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

func TestVirtualMachineClientBaseURISetToGov(t *testing.T) {
	// Get a client factory
	factory := multiEnvClientFactory{}

	// save any current env value and restore on exit
	currentEnv := os.Getenv(AzureEnvironmentEnvName)
	defer os.Setenv(AzureEnvironmentEnvName, currentEnv)

	// Set env var to gov
	os.Setenv(AzureEnvironmentEnvName, GovCloudEnvName)

	// Get a VM client
	client, _ := factory.GetVirtualMachinesClientE("")

	// Check for correct ARM URI
	assert.Equal(t, autorest.USGovernmentCloud.ResourceManagerEndpoint, client.BaseURI)
}

func TestVirtualMachineClientBaseURISetToPublic(t *testing.T) {
	// Get a client factory
	factory := multiEnvClientFactory{}

	// save any current env value and restore on exit
	currentEnv := os.Getenv(AzureEnvironmentEnvName)
	defer os.Setenv(AzureEnvironmentEnvName, currentEnv)

	// Set env var to gov
	os.Setenv(AzureEnvironmentEnvName, PublicCloudEnvName)

	// Get a VM client
	client, _ := factory.GetVirtualMachinesClientE("")

	// Check for correct ARM URI
	assert.Equal(t, autorest.PublicCloud.ResourceManagerEndpoint, client.BaseURI)
}

func TestSubscriptionClientBaseURISetToGov(t *testing.T) {
	// Get a client factory
	factory := multiEnvClientFactory{}

	// save any current env value and restore on exit
	currentEnv := os.Getenv(AzureEnvironmentEnvName)
	defer os.Setenv(AzureEnvironmentEnvName, currentEnv)

	// Set env var to gov
	os.Setenv(AzureEnvironmentEnvName, GovCloudEnvName)

	// Get a VM client
	client, _ := factory.GetSubscriptionClientE()

	// Check for correct ARM URI
	assert.Equal(t, autorest.USGovernmentCloud.ResourceManagerEndpoint, client.BaseURI)
}

func TestSubscriptionClientBaseURISetToPublic(t *testing.T) {
	// Get a client factory
	factory := multiEnvClientFactory{}

	// save any current env value and restore on exit
	currentEnv := os.Getenv(AzureEnvironmentEnvName)
	defer os.Setenv(AzureEnvironmentEnvName, currentEnv)

	// Set env var to gov
	os.Setenv(AzureEnvironmentEnvName, PublicCloudEnvName)

	// Get a VM client
	client, _ := factory.GetSubscriptionClientE()

	// Check for correct ARM URI
	assert.Equal(t, autorest.PublicCloud.ResourceManagerEndpoint, client.BaseURI)
}

func TestManagedClustersClientBaseURISetToGov(t *testing.T) {
	// Get a client factory
	factory := multiEnvClientFactory{}

	// save any current env value and restore on exit
	currentEnv := os.Getenv(AzureEnvironmentEnvName)
	defer os.Setenv(AzureEnvironmentEnvName, currentEnv)

	// Set env var to gov
	os.Setenv(AzureEnvironmentEnvName, GovCloudEnvName)

	// Get a VM client
	client, _ := factory.GetManagedClustersClientE("")

	// Check for correct ARM URI
	assert.Equal(t, autorest.USGovernmentCloud.ResourceManagerEndpoint, client.BaseURI)
}

func TestManagedClustersClientBaseURISetToPublic(t *testing.T) {
	// Get a client factory
	factory := multiEnvClientFactory{}

	// save any current env value and restore on exit
	currentEnv := os.Getenv(AzureEnvironmentEnvName)
	defer os.Setenv(AzureEnvironmentEnvName, currentEnv)

	// Set env var to gov
	os.Setenv(AzureEnvironmentEnvName, PublicCloudEnvName)

	// Get a VM client
	client, _ := factory.GetManagedClustersClientE("")

	// Check for correct ARM URI
	assert.Equal(t, autorest.PublicCloud.ResourceManagerEndpoint, client.BaseURI)
}
