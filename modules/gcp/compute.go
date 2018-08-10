package gcp

import (
	"context"
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"golang.org/x/oauth2/google"
	compute "google.golang.org/api/compute/v1"
)

// GetLabelsForComputeInstance returns all the tags for the given Compute Instance.
func GetLabelsForComputeInstance(t *testing.T, projectID string, zone string, instanceID string) map[string]string {
	labels, err := GetLabelsForComputeInstanceE(t, projectID, zone, instanceID)
	if err != nil {
		t.Fatal(err)
	}
	return labels
}

// GetLabelsForComputeInstanceE returns all the tags for the given Compute Instance.
func GetLabelsForComputeInstanceE(t *testing.T, projectID string, zone string, instanceID string) (map[string]string, error) {
	logger.Logf(t, "Getting Labels for Compute Instance %s", instanceID)

	ctx := context.Background()

	service, err := NewComputeServiceE(t)
	if err != nil {
		return nil, err
	}

	instance, err := service.Instances.Get(projectID, zone, instanceID).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("Instances.Get(%s) got error: %v", instanceID, err)
	}

	return instance.Labels, err
}

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

// AddLabelsToInstance adds the tags to the given taggable instance.
func AddLabelsToInstance(t *testing.T, projectID string, zone string, instance string, labels map[string]string) {
	err := AddLabelsToInstanceE(t, projectID, zone, instance, labels)
	if err != nil {
		t.Fatal(err)
	}
}

// AddLabelsToInstanceE adds the tags to the given taggable instance.
func AddLabelsToInstanceE(t *testing.T, projectID string, zone string, instance string, labels map[string]string) error {
	logger.Logf(t, "Adding labels to instance %s in zone %s", instance, zone)

	ctx := context.Background()

	service, err := NewComputeServiceE(t)
	if err != nil {
		return err
	}

	// Get the fingerprint of the existing labels
	existingInstance, err := service.Instances.Get(projectID, zone, instance).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("Instances.Get(%s) got error: %v", instance, err)
	}
	req := compute.InstancesSetLabelsRequest{Labels: labels, LabelFingerprint: existingInstance.LabelFingerprint}

	// Perform the SetLabels request
	if _, err := service.Instances.SetLabels(projectID, zone, instance, &req).Context(ctx).Do(); err != nil {
		return fmt.Errorf("Instances.SetLabels(%s) got error: %v", instance, err)
	}

	return err
}

// NewComputeService creates a new Compute service.
func NewComputeService(t *testing.T) *compute.Service {
	client, err := NewComputeServiceE(t)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// NewComputeServiceE creates a new Compute service.
func NewComputeServiceE(t *testing.T) (*compute.Service, error) {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		t.Fatalf("Failed to get default client: %v", err)
	}

	service, err := compute.New(client)
	if err != nil {
		return nil, err
	}

	return service, nil
}
