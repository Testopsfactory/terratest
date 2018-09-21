package gcp

import (
	"context"
	"fmt"
	"path"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"

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

// Corresponds to a GCP Zonal Instance Group (https://cloud.google.com/compute/docs/instance-groups/)
type ZonalInstanceGroup struct {
	*compute.InstanceGroup
}

// Corresponds to a GCP Regional Instance Group (https://cloud.google.com/compute/docs/instance-groups/)
type RegionalInstanceGroup struct {
	*compute.InstanceGroup
}

// NewInstance creates a new instance of the (GCP Compute) Instance type
func NewInstance(t *testing.T, projectID string, name string) *Instance {
	logger.Logf(t, "Getting Compute Instance %s", name)

	ctx := context.Background()

	service, err := NewComputeServiceE(t)
	if err != nil {
		t.Fatal(err)
	}

	instanceAggregatedList, err := service.Instances.AggregatedList(projectID).Context(ctx).Do()
	if err != nil {
		t.Fatalf("Instances.AggregatedList(%s) got error: %v", projectID, err)
	}

	for _, instanceList := range instanceAggregatedList.Items {
		for _, instance := range instanceList.Instances {
			if name == instance.Name {
				return &Instance{instance}
			}
		}
	}

	t.Fatalf("Compute Instance %s could not be found in project %s", name, projectID)
	return nil
}

// NewImage creates a new instance of the (GCP Compute) Image type
func NewImage(t *testing.T, projectID string, name string) *Image {
	logger.Logf(t, "Getting Image %s", name)

	service, err := NewComputeServiceE(t)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	req := service.Images.Get(projectID, name)
	image, err := req.Context(ctx).Do()
	if err != nil {
		t.Fatal(err)
	}

	return &Image{image}
}

// NewRegionalInstanceGroup creates a new instance of the Regional Instance Group type
func NewRegionalInstanceGroup(t *testing.T, projectID string, region string, name string) *RegionalInstanceGroup {
	logger.Logf(t, "Getting Regional Instance Group %s", name)

	service, err := NewComputeServiceE(t)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	req := service.RegionInstanceGroups.Get(projectID, region, name)
	instanceGroup, err := req.Context(ctx).Do()
	if err != nil {
		t.Fatal(err)
	}

	return &RegionalInstanceGroup{instanceGroup}
}

// NewZonalInstanceGroup creates a new instance of the Regional Instance Group type
func NewZonalInstanceGroup(t *testing.T, projectID string, zone string, name string) *ZonalInstanceGroup {
	logger.Logf(t, "Getting Zonal Instance Group %s", name)

	service, err := NewComputeServiceE(t)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	req := service.InstanceGroups.Get(projectID, zone, name)
	instanceGroup, err := req.Context(ctx).Do()
	if err != nil {
		t.Fatal(err)
	}

	return &ZonalInstanceGroup{instanceGroup}
}

// GetPublicIP gets the public IP address of the given Compute Instance.
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
		return "", fmt.Errorf("Attempted to get public IP of Compute Instance %s, but that Compute Instance does not have a public IP address", c.Name)
	}

	ip := c.NetworkInterfaces[0].AccessConfigs[0].NatIP

	return ip, nil
}

// GetLabels returns all the tags for the given Compute Instance.
func (c *Instance) GetLabels(t *testing.T) map[string]string {
	return c.Labels
}

// GetZone returns the Zone in which the Compute Instance is located.
func (c *Instance) GetZone(t *testing.T) string {
	return ZoneUrlToZone(c.Zone)
}

// SetLabels adds the tags to the given Compute Instance.
func (c *Instance) SetLabels(t *testing.T, projectID string, labels map[string]string) {
	err := c.SetLabelsE(t, projectID, labels)
	if err != nil {
		t.Fatal(err)
	}
}

// SetLabelsE adds the tags to the given Compute Instance.
func (c *Instance) SetLabelsE(t *testing.T, projectID string, labels map[string]string) error {
	logger.Logf(t, "Adding labels to instance %s in zone %s", c.Name, c.Zone)

	ctx := context.Background()

	service, err := NewComputeServiceE(t)
	if err != nil {
		return err
	}

	req := compute.InstancesSetLabelsRequest{Labels: labels, LabelFingerprint: c.LabelFingerprint}
	if _, err := service.Instances.SetLabels(projectID, c.GetZone(t), c.Name, &req).Context(ctx).Do(); err != nil {
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

// GetInstanceIds gets the IDs of Instances in the given Instance Group.
func (ig *ZonalInstanceGroup) GetInstanceIds(t *testing.T, projectID string) []string {
	ids, err := ig.GetInstanceIdsE(t, projectID)
	if err != nil {
		t.Fatal(err)
	}
	return ids
}

// GetInstanceIdsE gets the IDs of Instances in the given Zonal Instance Group.
func (ig *ZonalInstanceGroup) GetInstanceIdsE(t *testing.T, projectID string) ([]string, error) {
	logger.Logf(t, "Get instances for Zonal Instance Group %s", ig.Name)

	ctx := context.Background()

	service, err := NewComputeServiceE(t)
	if err != nil {
		return nil, err
	}

	requestBody := &compute.InstanceGroupsListInstancesRequest{
		InstanceState: "ALL",
	}

	instanceIDs := []string{}
	zone := ZoneUrlToZone(ig.Zone)

	req := service.InstanceGroups.ListInstances(projectID, zone, ig.Name, requestBody)
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

// GetInstanceIds gets the IDs of Instances in the given Regional Instance Group.
func (ig *RegionalInstanceGroup) GetInstanceIds(t *testing.T, projectID string) []string {
	ids, err := ig.GetInstanceIdsE(t, projectID)
	if err != nil {
		t.Fatal(err)
	}
	return ids
}

// GetInstanceIdsE gets the IDs of Instances in the given Regional Instance Group.
func (ig *RegionalInstanceGroup) GetInstanceIdsE(t *testing.T, projectID string) ([]string, error) {
	logger.Logf(t, "Get instances for Regional Instance Group %s", ig.Name)

	ctx := context.Background()

	service, err := NewComputeServiceE(t)
	if err != nil {
		return nil, err
	}

	requestBody := &compute.RegionInstanceGroupsListInstancesRequest{
		InstanceState: "ALL",
	}

	instanceIDs := []string{}
	region := RegionUrlToRegion(ig.Region)

	req := service.RegionInstanceGroups.ListInstances(projectID, region, ig.Name, requestBody)
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

// GetRandomInstance returns a randomly selected Instance from the Regional Instance Group
func (ig *RegionalInstanceGroup) GetRandomInstance(t *testing.T, projectID string) *Instance {
	instance, err := ig.GetRandomInstanceE(t, projectID)
	if err != nil {
		t.Fatal(err)
	}

	return instance
}

// GetRandomInstanceE returns a randomly selected Instance from the Regional Instance Group
func (ig *RegionalInstanceGroup) GetRandomInstanceE(t *testing.T, projectID string) (*Instance, error) {
	instanceIDs := ig.GetInstanceIds(t, projectID)

	clusterSize := int(ig.Size)
	randIndex := random.Random(1, clusterSize)

	if randIndex > len(instanceIDs) {
		return nil, fmt.Errorf("Could not find any instances in Regional Instance Group %s in Region %s", ig.Name, ig.Region)
	}

	instanceID := instanceIDs[randIndex-1]
	instance := NewInstance(t, projectID, instanceID)

	return instance, nil
}

// TODO: Is there some way to avoid the total duplication of this function and the Regional Instance Group equivalent in Go?
// GetRandomInstance returns a randomly selected Instance from the Zonal Instance Group
func (ig *ZonalInstanceGroup) GetRandomInstance(t *testing.T, projectID string) *Instance {
	instance, err := ig.GetRandomInstanceE(t, projectID)
	if err != nil {
		t.Fatal(err)
	}

	return instance
}

// GetRandomInstanceE returns a randomly selected Instance from the Zonal Instance Group
func (ig *ZonalInstanceGroup) GetRandomInstanceE(t *testing.T, projectID string) (*Instance, error) {
	instanceIDs := ig.GetInstanceIds(t, projectID)

	clusterSize := int(ig.Size)
	randIndex := random.Random(1, clusterSize)

	if randIndex > len(instanceIDs) {
		return nil, fmt.Errorf("Could not find any instances in Regional Instance Group %s in Region %s", ig.Name, RegionUrlToRegion(ig.Region))
	}

	instanceID := instanceIDs[randIndex-1]
	instance := NewInstance(t, projectID, instanceID)

	return instance, nil
}

// NewComputeService creates a new Compute service, which is used to make GCP API calls.
func NewComputeService(t *testing.T) *compute.Service {
	client, err := NewComputeServiceE(t)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// NewComputeServiceE creates a new Compute service, which is used to make GCP API calls.
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
