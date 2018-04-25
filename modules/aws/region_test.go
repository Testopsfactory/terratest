package aws

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetRandomRegionExcludesGloballyForbiddenRegions(t *testing.T) {
	t.Parallel()

	approvedRegions := []string{"ca-central-1", "us-east-1", "us-east-2", "us-west-1", "us-west-2", "eu-west-1", "eu-west-2", "eu-central-1", "ap-southeast-1", "ap-northeast-1", "ap-northeast-2", "ap-south-1"}
	forbiddenRegions := []string{"us-west-2", "ap-northeast-2"}

	for i := 0; i < 1000; i++ {
		randomRegion := GetRandomRegion(t, approvedRegions, forbiddenRegions)
		assert.NotContains(t, forbiddenRegions, randomRegion)
	}

}

func TestGetAvailabilityZones(t *testing.T) {
	t.Parallel()

	azs := GetAvailabilityZones(t,"us-west-2")

	// Every AWS account has access to different AZs, so he best we can do is make sure we get at least one back
	assert.True(t, len(azs) > 1)
}

