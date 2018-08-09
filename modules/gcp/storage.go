package gcp

import (
	"context"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/gruntwork-io/terratest/modules/logger"
)

// CreateStorageBucket creates a Google Cloud bucket with the given BucketAttrs. Note that Google Storage bucket names must be globally unique.
func CreateStorageBucket(t *testing.T, projectID string, name string, attr *storage.BucketAttrs) {
	err := CreateStorageBucketE(t, projectID, name, attr)
	if err != nil {
		t.Fatal(err)
	}
}

// CreateStorageBucketE creates a Google Cloud bucket with the given BucketAttrs. Note that Google Storage bucket names must be globally unique.
func CreateStorageBucketE(t *testing.T, projectID string, name string, attr *storage.BucketAttrs) error {
	logger.Logf(t, "Creating bucket %s", name)

	ctx := context.Background()

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	// Creates a Bucket instance.
	bucket := client.Bucket(name)

	// Creates the new bucket.
	return bucket.Create(ctx, projectID, attr)
}

// DeleteStorageBucket destroys the Google Storage bucket.
func DeleteStorageBucket(t *testing.T, name string) {
	err := DeleteStorageBucketE(t, name)
	if err != nil {
		t.Fatal(err)
	}
}

// DeleteStorageBucketE destroys the S3 bucket in the given region with the given name.
func DeleteStorageBucketE(t *testing.T, name string) error {
	logger.Logf(t, "Deleting bucket %s", name)

	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	return client.Bucket(name).Delete(ctx)
}

// AssertStorageBucketExists checks if the given storage bucket exists and fails the test if it does not.
func AssertStorageBucketExists(t *testing.T, name string) {
	err := AssertStorageBucketExistsE(t, name)
	if err != nil {
		t.Fatal(err)
	}
}

// AssertStorageBucketExistsE checks if the given storage bucket exists and returns an error if it does not.
func AssertStorageBucketExistsE(t *testing.T, name string) error {
	logger.Logf(t, "Finding bucket %s", name)

	ctx := context.Background()

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	// Creates a Bucket instance.
	bucket := client.Bucket(name)

	if _, err := bucket.Attrs(ctx); err != nil {
		// ErrBucketNotExist
		return err
	}

	it := bucket.Objects(ctx, nil)
	if _, err := it.Next(); err == storage.ErrBucketNotExist {
		return err
	}

	return nil
}
