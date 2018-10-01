package oci

import "os"

// You can set this environment variable to force Terratest to use a specific compartment.
const compartmentIDEnvVar = "TF_VAR_compartment_ocid"

// You can set this environment variable to force Terratest to use a specific availability domain
// rather than a random one. This is convenient when iterating locally.
const availabilityDomainEnvVar = "TF_VAR_availability_domain"

// You can set this environment variable to force Terratest to use a specific subnet.
const subnetIDEnvVar = "TF_VAR_subnet_ocid"

// You can set this environment variable to force Terratest to use a pass phrase.
const passPhraseEnvVar = "TF_VAR_pass_phrase"

// GetCompartmentIDFromEnvVar returns the compartment OCID for use with testing.
func GetCompartmentIDFromEnvVar() string {
	return os.Getenv(compartmentIDEnvVar)
}

// GetSubnetIDFromEnvVar returns the subnet OCID for use with testing.
func GetSubnetIDFromEnvVar() string {
	return os.Getenv(subnetIDEnvVar)
}

// GetPassPhraseFromEnvVar returns the pass phrase for use with testing.
func GetPassPhraseFromEnvVar() string {
	return os.Getenv(passPhraseEnvVar)
}
