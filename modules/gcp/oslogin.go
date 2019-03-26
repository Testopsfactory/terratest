package gcp

import (
	"context"
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/oslogin/v1"
)

// Add an OS Login SSH Key
func ImportSSHKey(t *testing.T, user, key string) {
	err := ImportSSHKeyE(t, user, key)
	if err != nil {
		t.Fatalf("Could not add SSH Key to user %s: %s", user, err)
	}
}

func ImportSSHKeyE(t *testing.T, user, key string) error {
	logger.Logf(t, "Importing SSH key for user %s", user)

	ctx := context.Background()
	service, err := NewOSLoginServiceE(t)
	if err != nil {
		return err
	}

	parent := fmt.Sprintf("users/%s", user)

	sshPublicKey := &oslogin.SshPublicKey{
		Key: key,
	}

	_, err = service.Users.ImportSshPublicKey(parent, sshPublicKey).Context(ctx).Do()
	if err != nil {
		return err
	}

	return nil
}

// Retrieve the login profile; OS Login + ephemeral gcloud keys + identities for
// the user.
func GetLoginProfile(t *testing.T, user string) *oslogin.LoginProfile {
	profile, err := GetLoginProfileE(t, user)
	if err != nil {
		t.Fatalf("Could not get login profile for user %s: %s", user, err)
	}

	return profile
}

func GetLoginProfileE(t *testing.T, user string) (*oslogin.LoginProfile, error) {
	logger.Logf(t, "Getting login profile for user %s", user)

	ctx := context.Background()
	service, err := NewOSLoginServiceE(t)
	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("users/%s", user)

	profile, err := service.Users.GetLoginProfile(name).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// NewOSLoginServiceE creates a new OS Login service, which is used to make OS Login API calls.
func NewOSLoginServiceE(t *testing.T) (*oslogin.Service, error) {
	ctx := context.Background()

	client, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("Failed to get default client: %v", err)
	}

	service, err := oslogin.New(client)
	if err != nil {
		return nil, err
	}

	return service, nil
}
