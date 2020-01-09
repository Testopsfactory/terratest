package gcp

import (
	"context"
	"fmt"
	"testing"

	cloudbuild "cloud.google.com/go/cloudbuild/apiv1/v2"
	"google.golang.org/api/iterator"
	cloudbuildpb "google.golang.org/genproto/googleapis/devtools/cloudbuild/v1"
)

// NewCloudBuildService creates a new Cloud Build service, which is used to make Cloud Build API calls.
func NewCloudBuildService(t *testing.T) *cloudbuild.Client {
	service, err := NewCloudBuildServiceE(t)
	if err != nil {
		t.Fatal(err)
	}
	return service
}

// NewCloudBuildServiceE creates a new Cloud Build service, which is used to make Cloud Build API calls.
func NewCloudBuildServiceE(t *testing.T) (*cloudbuild.Client, error) {
	ctx := context.Background()

	service, err := cloudbuild.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return service, nil
}

// CreateBuild creates a new build blocking until the operation is complete.
func CreateBuild(t *testing.T, projectID string, build *cloudbuildpb.Build) *cloudbuildpb.Build {
	out, err := CreateBuildE(t, projectID, build)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// CreateBuildE creates a new build blocking until the operation is complete.
func CreateBuildE(t *testing.T, projectID string, build *cloudbuildpb.Build) (*cloudbuildpb.Build, error) {
	ctx := context.Background()

	service, err := NewCloudBuildServiceE(t)
	if err != nil {
		return nil, err
	}

	req := &cloudbuildpb.CreateBuildRequest{
		ProjectId: projectID,
		Build:     build,
	}

	op, err := service.CreateBuild(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("CreateBuildE.CreateBuild(%s) got error: %v", projectID, err)
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		return nil, fmt.Errorf("CreateBuildE.Wait(%s) got error: %v", projectID, err)
	}

	return resp, nil
}

// GetBuild gets the given build.
func GetBuild(t *testing.T, projectID string, buildID string) *cloudbuildpb.Build {
	out, err := GetBuildE(t, projectID, buildID)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// GetBuildE gets the given build.
func GetBuildE(t *testing.T, projectID string, buildID string) (*cloudbuildpb.Build, error) {
	ctx := context.Background()

	service, err := NewCloudBuildServiceE(t)
	if err != nil {
		return nil, err
	}

	req := &cloudbuildpb.GetBuildRequest{
		ProjectId: projectID,
		Id:        buildID,
	}

	resp, err := service.GetBuild(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetBuildE.GetBuild(%s, %s) got error: %v", projectID, buildID, err)
	}

	return resp, nil
}

// GetBuilds gets the list of builds for a given project.
func GetBuilds(t *testing.T, projectID string) []*cloudbuildpb.Build {
	out, err := GetBuildsE(t, projectID)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// GetBuildsE gets the list of builds for a given project.
func GetBuildsE(t *testing.T, projectID string) ([]*cloudbuildpb.Build, error) {
	ctx := context.Background()

	service, err := NewCloudBuildServiceE(t)
	if err != nil {
		return nil, err
	}

	req := &cloudbuildpb.ListBuildsRequest{
		ProjectId: projectID,
	}

	it := service.ListBuilds(ctx, req)
	builds := []*cloudbuildpb.Build{}

	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("GetBuildsE.ListBuilds(%s) got error: %v", projectID, err)
		}

		builds = append(builds, resp)
	}

	return builds, nil
}

// GetBuildsForTrigger gets a list of builds for a specific cloud build trigger.
func GetBuildsForTrigger(t *testing.T, projectID string, triggerID string) []*cloudbuildpb.Build {
	out, err := GetBuildsForTriggerE(t, projectID, triggerID)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// GetBuildsForTriggerE gets a list of builds for a specific cloud build trigger.
func GetBuildsForTriggerE(t *testing.T, projectID string, triggerID string) ([]*cloudbuildpb.Build, error) {
	builds, err := GetBuildsE(t, projectID)
	if err != nil {
		return nil, fmt.Errorf("GetBuildsE.ListBuilds(%s) got error: %v", projectID, err)
	}

	filteredBuilds := []*cloudbuildpb.Build{}
	for _, build := range builds {
		if build.GetBuildTriggerId() == triggerID {
			filteredBuilds = append(filteredBuilds, build)
		}
	}

	return filteredBuilds, nil
}
