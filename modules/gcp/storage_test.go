package gcp

import (
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
)

var (
	projectId = os.Getenv("GOOGLE_STORAGE_PROJECT_ID")
)

func TestCreateAndDestroyStorageBucket(t *testing.T) {
	t.Parallel()

	id := random.UniqueId()
	logger.Logf(t, "Random values selected Id = %s\n", id)

	gsBucketName := "gruntwork-terratest-" + strings.ToLower(id)

	CreateStorageBucket(t, projectId, gsBucketName, nil)
	DeleteStorageBucket(t, gsBucketName)
}

func TestAssertStorageBucketExistsNoFalseNegative(t *testing.T) {
	t.Parallel()

	id := random.UniqueId()
	gsBucketName := "gruntwork-terratest-" + strings.ToLower(id)
	logger.Logf(t, "Random values selected Id = %s\n", id)

	CreateStorageBucket(t, projectId, gsBucketName, nil)
	defer DeleteStorageBucket(t, gsBucketName)

	AssertStorageBucketExists(t, gsBucketName)
}

func TestAssertStorageBucketExistsNoFalsePositive(t *testing.T) {
	t.Parallel()

	id := random.UniqueId()
	gsBucketName := "gruntwork-terratest-" + strings.ToLower(id)
	logger.Logf(t, "Random values selected Id = %s\n", id)

	// Don't create a new storage bucket so we can confirm that our function works as expected.

	err := AssertStorageBucketExistsE(t, gsBucketName)
	if err == nil {
		t.Fatalf("Function claimed that the Storage Bucket '%s' exists, but in fact it does not.", gsBucketName)
	}
}
