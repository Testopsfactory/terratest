// /*

// This file contains unit tests for the client factory implementation(s).

// */

package azure

// import (
// 	"os"
// 	"testing"

// 	autorest "github.com/Azure/go-autorest/autorest/azure"
// 	"github.com/stretchr/testify/assert"
// )

// // Local consts for this file only
// var GovCloudEnvName = "AzureUSGovernmentCloud"
// var PublicCloudEnvName = "AzurePublicCloud"
// var ChinaCloudEnvName = "AzureChinaCloud"
// var GermanyCloudEnvName = "AzureGermanCloud"

// func TestDefaultEnvIsPublicWhenNotSet(t *testing.T) {
// 	// Get a client factory
// 	factory := multiEnvClientFactory{}

// 	// save any current env value and restore on exit
// 	currentEnv := os.Getenv(AzureEnvironmentEnvName)
// 	defer os.Setenv(AzureEnvironmentEnvName, currentEnv)

// 	// Set env var to gov
// 	os.Setenv(AzureEnvironmentEnvName, "")

// 	// get the default
// 	env := factory.getDefaultEnvironmentName()

// 	// Make sure it's public cloud
// 	assert.Equal(t, autorest.PublicCloud.Name, env)
// }

// func TestDefaultEnvSetToGov(t *testing.T) {
// 	// Get a client factory
// 	factory := multiEnvClientFactory{}

// 	// save any current env value and restore on exit
// 	currentEnv := os.Getenv(AzureEnvironmentEnvName)
// 	defer os.Setenv(AzureEnvironmentEnvName, currentEnv)

// 	// Set env var to gov
// 	os.Setenv(AzureEnvironmentEnvName, GovCloudEnvName)

// 	// get the default
// 	env := factory.getDefaultEnvironmentName()

// 	// Make sure it's public cloud
// 	assert.Equal(t, autorest.USGovernmentCloud.Name, env)
// }

// func TestVirtualMachineClientBaseURISetCorrectly(t *testing.T) {
// 	var cases = []struct {
// 		EnvironmentName string
// 		ExpectedBaseURI string
// 	}{
// 		{GovCloudEnvName, autorest.USGovernmentCloud.ResourceManagerEndpoint},
// 		{PublicCloudEnvName, autorest.PublicCloud.ResourceManagerEndpoint},
// 		{ChinaCloudEnvName, autorest.ChinaCloud.ResourceManagerEndpoint},
// 		{GermanyCloudEnvName, autorest.GermanCloud.ResourceManagerEndpoint},
// 	}

// 	// Get a client factory
// 	factory := multiEnvClientFactory{}

// 	// save any current env value and restore on exit
// 	currentEnv := os.Getenv(AzureEnvironmentEnvName)
// 	defer os.Setenv(AzureEnvironmentEnvName, currentEnv)

// 	for _, tt := range cases {
// 		t.Run(tt.EnvironmentName, func(t *testing.T) {
// 			// Override env setting
// 			os.Setenv(AzureEnvironmentEnvName, tt.EnvironmentName)

// 			// Get a VM client
// 			client, _ := factory.GetVirtualMachinesClientE("")

// 			// Check for correct ARM URI
// 			assert.Equal(t, tt.ExpectedBaseURI, client.BaseURI)
// 		})
// 	}
// }

// func TestSubscriptionClientBaseURISetCorrectly(t *testing.T) {
// 	var cases = []struct {
// 		EnvironmentName string
// 		ExpectedBaseURI string
// 	}{
// 		{GovCloudEnvName, autorest.USGovernmentCloud.ResourceManagerEndpoint},
// 		{PublicCloudEnvName, autorest.PublicCloud.ResourceManagerEndpoint},
// 		{ChinaCloudEnvName, autorest.ChinaCloud.ResourceManagerEndpoint},
// 		{GermanyCloudEnvName, autorest.GermanCloud.ResourceManagerEndpoint},
// 	}

// 	// Get a client factory
// 	factory := multiEnvClientFactory{}

// 	// save any current env value and restore on exit
// 	currentEnv := os.Getenv(AzureEnvironmentEnvName)
// 	defer os.Setenv(AzureEnvironmentEnvName, currentEnv)

// 	for _, tt := range cases {
// 		t.Run(tt.EnvironmentName, func(t *testing.T) {
// 			// Override env setting
// 			os.Setenv(AzureEnvironmentEnvName, tt.EnvironmentName)

// 			// Get a VM client
// 			client, _ := factory.GetSubscriptionClientE()

// 			// Check for correct ARM URI
// 			assert.Equal(t, tt.ExpectedBaseURI, client.BaseURI)
// 		})
// 	}
// }

// func TestManagedClustersClientBaseURISetCorrectly(t *testing.T) {
// 	var cases = []struct {
// 		EnvironmentName string
// 		ExpectedBaseURI string
// 	}{
// 		{GovCloudEnvName, autorest.USGovernmentCloud.ResourceManagerEndpoint},
// 		{PublicCloudEnvName, autorest.PublicCloud.ResourceManagerEndpoint},
// 		{ChinaCloudEnvName, autorest.ChinaCloud.ResourceManagerEndpoint},
// 		{GermanyCloudEnvName, autorest.GermanCloud.ResourceManagerEndpoint},
// 	}

// 	// Get a client factory
// 	factory := multiEnvClientFactory{}

// 	// save any current env value and restore on exit
// 	currentEnv := os.Getenv(AzureEnvironmentEnvName)
// 	defer os.Setenv(AzureEnvironmentEnvName, currentEnv)

// 	for _, tt := range cases {
// 		t.Run(tt.EnvironmentName, func(t *testing.T) {
// 			// Override env setting
// 			os.Setenv(AzureEnvironmentEnvName, tt.EnvironmentName)

// 			// Get a VM client
// 			client, _ := factory.GetManagedClustersClientE("")

// 			// Check for correct ARM URI
// 			assert.Equal(t, tt.ExpectedBaseURI, client.BaseURI)
// 		})
// 	}
// }
