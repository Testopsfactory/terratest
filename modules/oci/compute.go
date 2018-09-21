package oci

import (
	"context"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/core"
	"testing"
)

// DeleteImage deletes a custom image with given OCID.
func DeleteImage(t *testing.T, ocid string) {
	err := DeleteImageE(t, ocid)
	if err != nil {
		t.Fatal(err)
	}
}

// DeleteImageE deletes a custom image with given OCID.
func DeleteImageE(t *testing.T, ocid string)  error {
	logger.Logf(t, "Deleting image with OCID %s", ocid)

	configProvider := common.DefaultConfigProvider()
	computeClient, err := core.NewComputeClientWithConfigurationProvider(configProvider)
	if err != nil {
		return err
	}

	_, err = computeClient.DeleteImage(context.Background(), core.DeleteImageRequest{ImageId: &ocid})
	return err
}
