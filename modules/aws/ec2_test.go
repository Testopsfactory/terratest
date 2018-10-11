package aws

import (
	"terratest/modules/packer"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"
)

// TestAmiIsPublic checks whether AMI is publicly accessible
func TestAmiIsPublic(t *testing.T) {
	t.Parallel()

	awsRegion := "eu-west-2"
	packerVars := map[string]string{
		"aws_region": awsRegion,
	}

	amiID := BuildAmi(t, packerVars)
	defer DeleteAmi(t, awsRegion, amiID)

	ec2Client := NewEc2Client(t, awsRegion)
	MakeAmiPublic(t, amiID, ec2Client)
	amiIsPublic := GetAmiPubliclyAccessible(t, awsRegion, amiID)
	assert.True(t, amiIsPublic)
}

// TestAmiIsPrivateAndNotSharedWithAccount checks whether AMI is private and not shared with any other account
func TestAmiIsPrivateAndNotSharedWithAccount(t *testing.T) {
	t.Parallel()

	requestingAccount := "123456789012"
	awsRegion := "eu-west-1"
	packerVars := map[string]string{
		"aws_region": awsRegion,
	}

	amiID := BuildAmi(t, packerVars)
	defer DeleteAmi(t, awsRegion, amiID)

	ec2Client := NewEc2Client(t, awsRegion)
	ShareAmi(t, amiID, "099720109477", ec2Client)
	accountsWithLaunchPermissions := GetAccountsWithLaunchPermissionsForAmi(t, awsRegion, amiID)
	assert.NotContains(t, accountsWithLaunchPermissions, requestingAccount)
}

// TestAmiIsPrivateAndSharedWithAccount checks whether AMI is private and shared with another account
func TestAmiIsPrivateAndSharedWithAccount(t *testing.T) {
	t.Parallel()

	requestingAccount := "099720109477"
	awsRegion := "eu-central-1"
	packerVars := map[string]string{
		"aws_region": awsRegion,
	}

	amiID := BuildAmi(t, packerVars)
	defer DeleteAmi(t, awsRegion, amiID)

	ec2Client := NewEc2Client(t, awsRegion)
	ShareAmi(t, amiID, requestingAccount, ec2Client)
	accountsWithLaunchPermissions := GetAccountsWithLaunchPermissionsForAmi(t, awsRegion, amiID)
	assert.Contains(t, accountsWithLaunchPermissions, requestingAccount)
}

func BuildAmi(t *testing.T, packerVars map[string]string) string {

	packerOptions := &packer.Options{
		// The path to where the Packer template is located
		Template: "../../examples/packer-basic-example/build.json",
		Vars:     packerVars,
		// Only build the AWS AMI
		Only: "amazon-ebs",
	}

	// Build the Docker image using Packer
	return packer.BuildArtifact(t, packerOptions)
}

func ShareAmi(t *testing.T, amiID string, accountID string, ec2Client *ec2.EC2) {
	input := &ec2.ModifyImageAttributeInput{
		ImageId: aws.String(amiID),
		LaunchPermission: &ec2.LaunchPermissionModifications{
			Add: []*ec2.LaunchPermission{
				{
					UserId: aws.String(accountID),
				},
			},
		},
	}
	_, err := ec2Client.ModifyImageAttribute(input)
	if err != nil {
		t.Fatal(err)
	}
}

func MakeAmiPublic(t *testing.T, amiID string, ec2Client *ec2.EC2) {
	input := &ec2.ModifyImageAttributeInput{
		ImageId: aws.String(amiID),
		LaunchPermission: &ec2.LaunchPermissionModifications{
			Add: []*ec2.LaunchPermission{
				{
					Group: aws.String("all"),
				},
			},
		},
	}
	_, err := ec2Client.ModifyImageAttribute(input)
	if err != nil {
		t.Fatal(err)
	}
}
