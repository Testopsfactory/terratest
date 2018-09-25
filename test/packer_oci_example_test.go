package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/packer"
	"github.com/sw-samuraj/terratest/modules/oci"
)

// An example of how to test the Packer template in examples/packer-basic-example using Terratest.
func TestPackerOciExample(t *testing.T) {
	t.Parallel()

	compartmentID := oci.GetCompartmentIDFromEnvVar()
	baseImageID := oci.GetMostRecentImageID(t, compartmentID, "Canonical Ubuntu", "18.04")
	availabilityDomain := oci.GetRandomAvailabilityDomain(t, compartmentID)

	packerOptions := &packer.Options{
		// The path to where the Packer template is located
		Template: "../examples/packer-basic-example/build.json",

		// Variables to pass to our Packer build using -var options
		Vars: map[string]string{
			"oci_compartment_ocid":    compartmentID,
			"oci_base_image_ocid":     baseImageID,
			"oci_availability_domain": availabilityDomain,
			"oci_subnet_ocid":         "ocid1.subnet.oc1.phx.aaaaaaaa",
			"oci_key_file":            fmt.Sprintf("%s/.oci/oci_api_key.pem", os.Getenv("HOME")),
			"oci_pass_phrase":         "my-secret-pass-phrase",
		},

		// Only build an OCI image
		Only: "oracle-oci",
	}

	// Make sure the Packer build completes successfully
	ocid := packer.BuildArtifact(t, packerOptions)

	// Delete the OCI image after we're done
	defer oci.DeleteImage(t, ocid)
}
