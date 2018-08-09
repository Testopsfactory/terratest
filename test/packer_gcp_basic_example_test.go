package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/packer"
)

// An example of how to test the Packer template in examples/packer-basic-example using Terratest.
func TestPackerGCPBasicExample(t *testing.T) {
	t.Parallel()

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	//awsRegion := aws.GetRandomRegion(t, nil, nil)

	packerOptions := &packer.Options{
		// The path to where the Packer template is located
		Template: "../examples/packer-basic-example/build.json",

		// Variables to pass to our Packer build using -var options
		Vars: map[string]string{
			"gcp_project_id": projectId,
			"gcp_zone":       "us-central1-a",
		},

		// Only build the AWS AMI
		Only: "googlecompute",
	}

	// Make sure the Packer build completes successfully
	imageID := packer.BuildArtifact(t, packerOptions)

	// Delete the Image after we're done
	defer gcp.DeleteImage(t, projectId, imageID)
}
