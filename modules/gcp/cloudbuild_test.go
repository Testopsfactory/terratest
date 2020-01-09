package gcp

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	cloudbuildpb "google.golang.org/genproto/googleapis/devtools/cloudbuild/v1"
)

func TestCreateBuild(t *testing.T) {
	t.Parallel()
	// TODO

	// 1. Create tarball with single Dockerfile
	// 2. Create GCS bucket
	// 3. Upload tarball to GCS Bucket
	// 4. Create new build
	// 4. Trigger a build using the Cloud Build API

	// Create and add some files to the archive.
	tarball := createSampleAppTarball(t)

	// Create GCS bucket
	projectID := GetGoogleProjectIDFromEnvVar(t)
	id := random.UniqueId()
	gsBucketName := "cloud-build-terratest-" + strings.ToLower(id)
	sampleAppPath := "docker-example.tar.gz"
	imagePath := fmt.Sprintf("gcr.io/%s/test-image-%s", projectID, strings.ToLower(id))

	logger.Logf(t, "Random values selected Bucket Name = %s\n", gsBucketName)

	CreateStorageBucket(t, projectID, gsBucketName, nil)
	defer DeleteStorageBucket(t, gsBucketName)

	// Write the compressed archive to the storage bucket
	objectURL := WriteBucketObject(t, gsBucketName, sampleAppPath, tarball, "application/gzip")
	logger.Logf(t, "Got URL: %s", objectURL)

	// Create a new build
	build := &cloudbuildpb.Build{
		Source: &cloudbuildpb.Source{
			Source: &cloudbuildpb.Source_StorageSource{
				StorageSource: &cloudbuildpb.StorageSource{
					Bucket: gsBucketName,
					Object: sampleAppPath,
				},
			},
		},
		Steps: []*cloudbuildpb.BuildStep{{
			Name: "gcr.io/cloud-builders/docker",
			Args: []string{"build", "-t", imagePath, "."},
		}},
		Images: []string{imagePath},
	}

	// CreateBuild blocks until the build is complete
	b := CreateBuild(t, projectID, build)

	// TODO - Delete the GCR image

	// Empty the storage bucket so we can delete it
	defer EmptyStorageBucket(t, gsBucketName)
}

func TestGetBuild(t *testing.T) {
	t.Parallel()
}

func TestGetBuilds(t *testing.T) {
	t.Parallel()
}

func TestGetBuildsForTrigger(t *testing.T) {
	t.Parallel()
}

func createCloudBuildTrigger(t *testing.T, projectID string) {

}

func destroyCloudBuildTrigger(t *testing.T, projectID string) {

}

func createSampleAppTarball(t *testing.T) *bytes.Reader {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	file := `FROM busybox:latest
MAINTAINER Rob Morgan (rob@gruntwork.io)
	`

	hdr := &tar.Header{
		Name: "Dockerfile",
		Mode: 0600,
		Size: int64(len(file)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write([]byte(file)); err != nil {
		t.Fatal(err)
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}

	// gzip the tar archive
	var zbuf bytes.Buffer
	gzw := gzip.NewWriter(&zbuf)
	if _, err := gzw.Write(buf.Bytes()); err != nil {
		t.Fatal(err)
	}
	if err := gzw.Close(); err != nil {
		t.Fatal(err)
	}

	// return the compressed buffer
	return bytes.NewReader(zbuf.Bytes())
}
