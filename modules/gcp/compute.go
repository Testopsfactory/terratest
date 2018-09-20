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

// Corresponds to a GCP Compute Instance (https://cloud.google.com/compute/docs/instances/)
type Instance struct {
	*compute.Instance
}

// Corresponds to a GCP Image (https://cloud.google.com/compute/docs/images)
type Image struct {
	*compute.Image
}

// Zonal and Regional Instance Groups both use the same type of underlying InstanceGroup resource, but they use different
// GCP API calls, so it makes sense to maintain separate types for each of them.

// Corresponds to a GCP Zonal Instance Group (https://cloud.google.com/compute/docs/instance-groups/)
type ZonalInstanceGroup struct {
	*compute.InstanceGroup
}

// Corresponds to a GCP Regional Instance Group (https://cloud.google.com/compute/docs/instance-groups/)
type RegionalInstanceGroup struct {
	*compute.InstanceGroup
}

func newComputeInstance(instance *compute.Instance) *Instance {
	return &Instance{instance}
}

func newComputeImage(name string) *Image {
	image := &compute.Image{
		Name: name,
	}
	return &Image{image}
}

// GetInstance gets the given Instance in the given region.
func GetInstance(t *testing.T, projectID string, instanceName string) *Instance {
	instance, err := GetInstanceE(t, projectID, instanceName)
	if err != nil {
		t.Fatal(err)
	}

	return instance
}

// GetInstanceE gets the given Compute Instance for the given Comput Instance name.
func GetInstanceE(t *testing.T, projectID string, instanceName string) (*Instance, error) {
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
func (c *Instance) GetPublicIp(t *testing.T) string {
	ip, err := c.GetPublicIpE(t)
	if err != nil {
		t.Fatal(err)
	}
	return ip
}

// GetPublicIpE gets the public IP address of the given Compute Instance.
func (c *Instance) GetPublicIpE(t *testing.T) (string, error) {
	// If there are no accessConfigs specified, then this instance will have no external internet access:
	// https://cloud.google.com/compute/docs/reference/rest/v1/instances.
	if len(c.NetworkInterfaces[0].AccessConfigs) == 0 {
		return "", fmt.Errorf("Attempted to get public IP of Compute Instance %s, but that Compute Instance does not have a public IP address", c.instance.Name)
	}

	ip := c.NetworkInterfaces[0].AccessConfigs[0].NatIP

	return ip, nil
}

// GetLabelsForComputeInstanceE returns all the tags for the given Compute Instance.
func (c *Instance) GetLabels(t *testing.T) map[string]string {
	return c.Labels
}

// AddLabelsToInstance adds the tags to the given taggable instance.
func (c *Instance) AddLabelsToInstance(t *testing.T, projectID string, labels map[string]string) {
	err := c.AddLabelsToInstanceE(t, projectID, labels)
	if err != nil {
		t.Fatal(err)
	}
}

// AddLabelsToInstanceE adds the tags to the given taggable instance.
func (c *Instance) AddLabelsToInstanceE(t *testing.T, projectID string, labels map[string]string) error {
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
func (i *Image) DeleteImage(t *testing.T, projectID string) {
	err := i.DeleteImageE(t, projectID)
	if err != nil {
		t.Fatal(err)
	}
}

// DeleteImageE deletes the given Compute Image.
func (i *Image) DeleteImageE(t *testing.T, projectID string) error {
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
func (ig *ZonalInstanceGroup) GetInstanceIds(t *testing.T, projectID string) []string {
	ids, err := ig.GetInstanceIdsE(t, projectID)
	if err != nil {
		t.Fatal(err)
	}
	return ids
}

// GetInstanceIdsForZonalInstanceGroupE gets the IDs of Instances in the given Zonal Instance Group.
func (ig *ZonalInstanceGroup) GetInstanceIdsE(t *testing.T, projectID string) ([]string, error) {
	logger.Logf(t, "Get instances for Zonal Instance Group %s in zone %s", ig.Name, ig.Zone)

	ctx := context.Background()

	service, err := NewComputeServiceE(t)
	if err != nil {
		return nil, err
	}

	requestBody := &compute.InstanceGroupsListInstancesRequest{
		InstanceState: "ALL",
	}

	instanceIDs := []string{}
	req := service.InstanceGroups.ListInstances(projectID, ig.Zone, ig.Name, requestBody)
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
		return nil, fmt.Errorf("InstanceGroups.ListInstances(%s) got error: %v", ig.Name, err)
	}

	return instanceIDs, nil
}

// GetInstanceIdsForInstanceGroup gets the IDs of Instances in the given Regional Instance Group.
func (ig *RegionalInstanceGroup) GetInstanceIds(t *testing.T, projectID string) []string {
	ids, err := ig.GetInstanceIdsE(t, projectID)
	if err != nil {
		t.Fatal(err)
	}
	return ids
}

// GetInstanceIdsForRegionalInstanceGroupE gets the IDs of Instances in the given Regional Instance Group.
func (ig *RegionalInstanceGroup) GetInstanceIdsE(t *testing.T, projectID string) ([]string, error) {
	logger.Logf(t, "Get instances for Regional Instance Group %s in Region %s", ig.Name, ig.Region)

	ctx := context.Background()

	service, err := NewComputeServiceE(t)
	if err != nil {
		return nil, err
	}

	requestBody := &compute.RegionInstanceGroupsListInstancesRequest{
		InstanceState: "ALL",
	}

	instanceIDs := []string{}
	req := service.RegionInstanceGroups.ListInstances(projectID, ig.Region, ig.Name, requestBody)
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
		return nil, fmt.Errorf("InstanceGroups.ListInstances(%s) got error: %v", ig.Name, err)
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
