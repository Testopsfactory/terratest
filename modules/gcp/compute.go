package gcp

import (
	"context"
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
)

const (
	cloudScope = "https://www.googleapis.com/auth/cloud-platform"
)

// DeleteImage deletes the given Compute Image.
func DeleteImage(t *testing.T, projectID string, imageID string) {
	err := DeleteImageE(t, projectID, imageID)
	if err != nil {
		t.Fatal(err)
	}
}

// DeleteImageE deletes the given Compute Image.
func DeleteImageE(t *testing.T, projectID string, imageID string) error {
	logger.Logf(t, "Destroying Image %s", imageID)

	ctx := context.Background()

	service, err := NewComputeServiceE(t)
	if err != nil {
		return err
	}

	if _, err := service.Images.Delete(projectID, imageID).Context(ctx).Do(); err != nil {
		return fmt.Errorf("Images.Delete(%s) got error: %v", imageID, err)
	}

	return err
}

// NewComputeService creates a Compute client.
func NewComputeService(t *testing.T) *compute.Service {
	client, err := NewComputeServiceE(t)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// NewComputeServiceE creates a Compute client.
func NewComputeServiceE(t *testing.T) (*compute.Service, error) {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, cloudScope)
	if err != nil {
		t.Fatalf("failed to get default client: %v", err)
	}

	service, err := compute.New(client)
	if err != nil {
		return nil, err
	}

	return service, nil
}
