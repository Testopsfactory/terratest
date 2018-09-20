package gcp

import (
	"context"
	"fmt"
	"path"
	"testing"

	"google.golang.org/api/compute/v1"

	"github.com/gruntwork-io/terratest/modules/logger"
	"golang.org/x/oauth2/google"
)

type ComputeInstance struct {
	*compute.Instance
}

type ComputeImage struct {
	*compute.Image
}

func newComputeInstance(instance *compute.Instance) *ComputeInstance {
	return &ComputeInstance{instance}
}

func newComputeImage(image *compute.Image) *ComputeImage {
	return &ComputeImage{image}
}

// GetInstance gets the given Instance in the given region.
func GetInstance(t *testing.T, projectID string, instanceName string) *ComputeInstance {
	instance, err := GetInstanceE(t, projectID, instanceName)
	if err != nil {
		t.Fatal(err)
	}

	return instance
}

// GetInstanceE gets the given Compute Instance for the given Comput Instance name.
func GetInstanceE(t *testing.T, projectID string, instanceName string) (*ComputeInstance, error) {
	logger.Logf(t, "Getting Compute Instance %s", instanceName)

	ctx := context.Background()

	service, err := NewComputeServiceE(t)
	if err != nil {
		return nil, err
	}

	instanceAggregatedList, err := service.Instances.AggregatedList(projectID).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("Instances.AggregatedList(%s) got error: %v", projectID, err)
	}

	for _, instanceList := range instanceAggregatedList.Items {
		for _, instance := range instanceList.Instances {
			if instanceName == instance.Name {
				return newComputeInstance(instance), nil
			}
		}
	}

	return nil, fmt.Errorf("Compute Instance %s could not be found in project %s", instanceName, projectID)
}

// GetPublicIPOfComputeInstance gets the public IP address of the given Compute Instance.
func (c *ComputeInstance) GetPublicIp(t *testing.T) string {
	ip, err := c.GetPublicIpE(t)
	if err != nil {
		t.Fatal(err)
	}
	return ip
}

// GetPublicIpE gets the public IP address of the given Compute Instance.
func (c *ComputeInstance) GetPublicIpE(t *testing.T) (string, error) {
	// If there are no accessConfigs specified, then this instance will have no external internet access:
	// https://cloud.google.com/compute/docs/reference/rest/v1/instances.
	if len(c.NetworkInterfaces[0].AccessConfigs) == 0 {
		return "", fmt.Errorf("Attempted to get public IP of Compute Instance %s, but that Compute Instance does not have a public IP address", c.instance.Name)
	}

	ip := c.NetworkInterfaces[0].AccessConfigs[0].NatIP

	return ip, nil
}

// GetLabelsForComputeInstanceE returns all the tags for the given Compute Instance.
func (c *ComputeInstance) GetLabels(t *testing.T) map[string]string {
	return c.Labels
}

// AddLabelsToInstance adds the tags to the given taggable instance.
func (c *ComputeInstance) AddLabelsToInstance(t *testing.T, projectID string, labels map[string]string) {
	err := c.AddLabelsToInstanceE(t, projectID, labels)
	if err != nil {
		t.Fatal(err)
	}
}

// AddLabelsToInstanceE adds the tags to the given taggable instance.
func (c *ComputeInstance) AddLabelsToInstanceE(t *testing.T, projectID string, labels map[string]string) error {
	logger.Logf(t, "Adding labels to instance %s in zone %s", c.Name, c.Zone)

	ctx := context.Background()

	service, err := NewComputeServiceE(t)
	if err != nil {
		return err
	}

	req := compute.InstancesSetLabelsRequest{Labels: labels, LabelFingerprint: c.LabelFingerprint}
	if _, err := service.Instances.SetLabels(projectID, c.Zone, c.Name, &req).Context(ctx).Do(); err != nil {
		return fmt.Errorf("Instances.SetLabels(%s) got error: %v", c.Name, err)
	}

	return err
}

// DeleteImage deletes the given Compute Image.
func (i *ComputeImage) DeleteImage(t *testing.T, projectID string) {
	err := i.DeleteImageE(t, projectID)
	if err != nil {
		t.Fatal(err)
	}
}

// DeleteImageE deletes the given Compute Image.
func (i *ComputeImage) DeleteImageE(t *testing.T, projectID string) error {
	logger.Logf(t, "Destroying Image %s", i.Name)

	ctx := context.Background()

	service, err := NewComputeServiceE(t)
	if err != nil {
		return err
	}

	if _, err := service.Images.Delete(projectID, i.Name).Context(ctx).Do(); err != nil {
		return fmt.Errorf("Images.Delete(%s) got error: %v", i.Name, err)
	}

	return err
}

// GetInstanceIdsForInstanceGroup gets the IDs of Instances in the given Instance Group.
func GetInstanceIdsForZonalInstanceGroup(t *testing.T, projectID string, zone string, groupName string) []string {
	ids, err := GetInstanceIdsForZonalInstanceGroupE(t, projectID, zone, groupName)
	if err != nil {
		t.Fatal(err)
	}
	return ids
}

// GetInstanceIdsForZonalInstanceGroupE gets the IDs of Instances in the given Zonal Instance Group.
func GetInstanceIdsForZonalInstanceGroupE(t *testing.T, projectID string, zone string, groupName string) ([]string, error) {
	logger.Logf(t, "Get instances for Zonal Instance Group %s in zone %s", groupName, zone)

	ctx := context.Background()

	service, err := NewComputeServiceE(t)
	if err != nil {
		return nil, err
	}

	requestBody := &compute.InstanceGroupsListInstancesRequest{
		InstanceState: "ALL",
	}

	instanceIDs := []string{}
	req := service.InstanceGroups.ListInstances(projectID, zone, groupName, requestBody)
	if err := req.Pages(ctx, func(page *compute.InstanceGroupsListInstances) error {
		for _, instance := range page.Items {
			// For some reason service.InstanceGroups.ListInstances returns us a collection
			// with Instance URLs and we need only the Instance ID for the next call. Use
			// the path functions to chop the Instance ID off the end of the URL.
			instanceID := path.Base(instance.Instance)
			instanceIDs = append(instanceIDs, instanceID)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("InstanceGroups.ListInstances(%s) got error: %v", groupName, err)
	}

	return instanceIDs, nil
}

// GetInstanceIdsForInstanceGroup gets the IDs of Instances in the given Regional Instance Group.
func GetInstanceIdsForRegionalInstanceGroup(t *testing.T, projectID string, region string, groupName string) []string {
	ids, err := GetInstanceIdsForRegionalInstanceGroupE(t, projectID, region, groupName)
	if err != nil {
		t.Fatal(err)
	}
	return ids
}

// GetInstanceIdsForRegionalInstanceGroupE gets the IDs of Instances in the given Regional Instance Group.
func GetInstanceIdsForRegionalInstanceGroupE(t *testing.T, projectID string, region string, groupName string) ([]string, error) {
	logger.Logf(t, "Get instances for Regional Instance Group %s in Region %s", groupName, region)

	ctx := context.Background()

	service, err := NewComputeServiceE(t)
	if err != nil {
		return nil, err
	}

	requestBody := &compute.RegionInstanceGroupsListInstancesRequest{
		InstanceState: "ALL",
	}

	instanceIDs := []string{}
	req := service.RegionInstanceGroups.ListInstances(projectID, region, groupName, requestBody)
	if err := req.Pages(ctx, func(page *compute.RegionInstanceGroupsListInstances) error {
		for _, instance := range page.Items {
			// For some reason service.InstanceGroups.ListInstances returns us a collection
			// with Instance URLs and we need only the Instance ID for the next call. Use
			// the path functions to chop the Instance ID off the end of the URL.
			instanceID := path.Base(instance.Instance)
			instanceIDs = append(instanceIDs, instanceID)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("InstanceGroups.ListInstances(%s) got error: %v", groupName, err)
	}

	return instanceIDs, nil
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
		return nil, fmt.Errorf("Failed to get default client: %v", err)
	}

	service, err := compute.New(client)
	if err != nil {
		return nil, err
	}

	return service, nil
}
