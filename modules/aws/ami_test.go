package aws

import "testing"

func TestGetUbuntuAmiReturnsSomeAmi(t *testing.T) {
	t.Parallel()

	amiId := GetUbuntu1404Ami(t, "us-east-1")
	if amiId[:4] != "ami-" {
		t.Fatalf("Expected: string formatted like ami-*******. Received: %s", amiId)
	}
}
