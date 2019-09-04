package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestGetDefaultVpc(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	vpc := GetDefaultVpc(t, region)

	assert.NotEmpty(t, vpc.Name)
	assert.True(t, len(vpc.Subnets) > 0)
	assert.Regexp(t, "^vpc-[[:alnum:]]+$", vpc.Id)
}

func TestGetVpcById(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	vpc := createVpc(t, region)
	defer deleteVpc(t, *vpc.VpcId, region)

	vpcTest := GetVpcById(t, *vpc.VpcId, region)
	assert.Equal(t, *vpc.VpcId, vpcTest.Id)
}

func TestGetVpcsE(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	azs := GetAvailabilityZones(t, region)

	isDefaultFilterName := "isDefault"
	isDefaultFilterValue := "true"

	defaultVpcFilter := ec2.Filter{Name: &isDefaultFilterName, Values: []*string{&isDefaultFilterValue}}
	vpcs, _ := GetVpcsE(t, []*ec2.Filter{&defaultVpcFilter}, region)

	require.Equal(t, len(vpcs), 1)
	assert.NotEmpty(t, vpcs[0].Name)

	// the default VPC has by default one subnet per availability zone
	// https://docs.aws.amazon.com/vpc/latest/userguide/default-vpc.html
	assert.Equal(t, len(vpcs[0].Subnets), len(azs))
}

func TestGetFirstTwoOctets(t *testing.T) {
	t.Parallel()

	firstTwo := GetFirstTwoOctets("10.100.0.0/28")
	if firstTwo != "10.100" {
		t.Errorf("Received: %s, Expected: 10.100", firstTwo)
	}
}

func createVpc(t *testing.T, region string) ec2.Vpc {
	ec2Client := NewEc2Client(t, region)

	createVpcOutput, err := ec2Client.CreateVpc(&ec2.CreateVpcInput{
		CidrBlock: aws.String("10.10.0.0/16"),
	})

	require.NoError(t, err)
	return *createVpcOutput.Vpc
}

func deleteVpc(t *testing.T, vpcId string, region string) {
	ec2Client := NewEc2Client(t, region)

	_, err := ec2Client.DeleteVpc(&ec2.DeleteVpcInput{
		VpcId: aws.String(vpcId),
	})
	require.NoError(t, err)
}
