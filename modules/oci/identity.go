package oci

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/identity"
)

// You can set this environment variable to force Terratest to use a specific compartment.
const compartmentIDOverrideEnvVarName = "TF_VAR_compartment_ocid"

// You can set this environment variable to force Terratest to use a specific availability domain
// rather than a random one. This is convenient when iterating locally.
const availabilityDomainOverrideEnvVarName = "TF_VAR_AD"

// TODO Guido: document this
func GetRandomAvailabilityDomain(t *testing.T, compartmentID string) string {
	ad, err := GetRandomAvailabilityDomainE(t, compartmentID)
	if err != nil {
		t.Fatal(err)
	}
	return ad
}

// TODO Guido: document this
func GetRandomAvailabilityDomainE(t *testing.T, compartmentID string) (string, error) {
	adFromEnvVar := os.Getenv(availabilityDomainOverrideEnvVarName)
	if adFromEnvVar != "" {
		logger.Logf(t, "Using availability domain %s from environment variable %s", adFromEnvVar, availabilityDomainOverrideEnvVarName)
		return adFromEnvVar, nil
	}

	allADs, err := GetAllADsE(t, compartmentID)
	if err != nil {
		return "", err
	}

	ad := random.RandomString(allADs)

	logger.Logf(t, "Using availability domain %s", ad)
	return ad, nil
}

// TODO Guido: document this
func GetAllADs(t *testing.T, compartmentID string) []string {
	ads, err := GetAllADsE(t, compartmentID)
	if err != nil {
		t.Fatal(err)
	}
	return ads
}

// TODO Guido: document this
func GetAllADsE(t *testing.T, compartmentID string) ([]string, error) {
	configProvider := common.DefaultConfigProvider()
	client, err := identity.NewIdentityClientWithConfigurationProvider(configProvider)
	if err != nil {
		return nil, err
	}

	request := identity.ListAvailabilityDomainsRequest{CompartmentId: &compartmentID}
	response, err := client.ListAvailabilityDomains(context.Background(), request)
	if err != nil {
		return nil, err
	}

	if len(response.Items) == 0 {
		return nil, fmt.Errorf("No availability domains found in the %s compartment", compartmentID)
	}

	return adNames(response.Items), nil
}

// GetCompartmentIDFromEnvVar returns the Compartment for use with testing.
func GetCompartmentIDFromEnvVar() string {
	if compartmentID := os.Getenv(compartmentIDOverrideEnvVarName); compartmentID != "" {
		return compartmentID
	}
	return ""
}

func adNames(ads []identity.AvailabilityDomain) []string {
	names := []string{}
	for _, ad := range ads {
		names = append(names, *ad.Name)
	}
	return names
}
