package aws

import (
	"testing"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws"
	"sort"
	"time"
	"fmt"
)

const CanonicalAccountId = "099720109477"
const CentOsAccountId = "679593333241"
const AmazonAccountId = "amazon"

// Get the ID of the most recent AMI in the given region that has the given owner and matches the given filters. Each
// filter should correspond to the name and values of a filter supported by DescribeImagesInput:
// https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#DescribeImagesInput
func GetMostRecentAmiId(t *testing.T, region string, ownerId string, filters map[string][]string) string {
	amiId, err := GetMostRecentAmiIdE(t, region, ownerId, filters)
	if err != nil {
		t.Fatal(err)
	}
	return amiId
}

// Get the ID of the most recent AMI in the given region that has the given owner and matches the given filters. Each
// filter should correspond to the name and values of a filter supported by DescribeImagesInput:
// https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#DescribeImagesInput
func GetMostRecentAmiIdE(t *testing.T, region string, ownerId string, filters map[string][]string) (string, error) {
	ec2Client, err := NewEc2ClientE(t, region)
	if err != nil {
		return "", err
	}

	ec2Filters := []*ec2.Filter{}
	for name, values := range filters {
		ec2Filters = append(ec2Filters, &ec2.Filter{Name: aws.String(name), Values: aws.StringSlice(values)})
	}

	input := ec2.DescribeImagesInput{
		Filters: ec2Filters,
		Owners:  []*string{aws.String(ownerId)},
	}

	out, err := ec2Client.DescribeImages(&input)
	if err != nil {
		return "", err
	}

	if len(out.Images) == 0 {
		return "", NoImagesFound{Region: region, OwnerId: ownerId, Filters: filters}
	}

	mostRecentImage := mostRecentAmi(out.Images)
	return aws.StringValue(mostRecentImage.ImageId), nil
}

// Image sorting code borrowed from: https://github.com/hashicorp/packer/blob/7f4112ba229309cfc0ebaa10ded2abdfaf1b22c8/builder/amazon/common/step_source_ami_info.go
type imageSort []*ec2.Image

func (a imageSort) Len() int      { return len(a) }
func (a imageSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a imageSort) Less(i, j int) bool {
	itime, _ := time.Parse(time.RFC3339, *a[i].CreationDate)
	jtime, _ := time.Parse(time.RFC3339, *a[j].CreationDate)
	return itime.Unix() < jtime.Unix()
}

// Returns the most recent AMI out of a slice of images.
func mostRecentAmi(images []*ec2.Image) *ec2.Image {
	sortedImages := images
	sort.Sort(imageSort(sortedImages))
	return sortedImages[len(sortedImages)-1]
}

// Get the ID of the most recent Ubuntu 14.04 HVM x86_64 EBS GP2 AMI in the given region
func GetUbuntu1404Ami(t *testing.T, region string) string {
	amiId, err := GetUbuntu1404AmiE(t, region)
	if err != nil {
		t.Fatal(err)
	}
	return amiId
}

// Get the ID of the most recent Ubuntu 14.04 HVM x86_64 EBS GP2 AMI in the given region
func GetUbuntu1404AmiE(t *testing.T, region string) (string, error) {
	filters := map[string][]string{
		"name":                             {"*ubuntu-trusty-14.04-amd64-server-*"},
		"virtualization-type":              {"hvm"},
		"architecture":                     {"x86_64"},
		"root-device-type":                 {"ebs"},
		"block-device-mapping.volume-type": {"gp2"},
	}

	return GetMostRecentAmiIdE(t, region, CanonicalAccountId, filters)
}

// Get the ID of the most recent Ubuntu 16.04 HVM x86_64 EBS GP2 AMI in the given region
func GetUbuntu1604Ami(t *testing.T, region string) string {
	amiId, err := GetUbuntu1604AmiE(t, region)
	if err != nil {
		t.Fatal(err)
	}
	return amiId
}

// Get the ID of the most recent Ubuntu 16.04 HVM x86_64 EBS GP2 AMI in the given region
func GetUbuntu1604AmiE(t *testing.T, region string) (string, error) {
	filters := map[string][]string{
		"name":                             {"*ubuntu-xenial-16.04-amd64-server-*"},
		"virtualization-type":              {"hvm"},
		"architecture":                     {"x86_64"},
		"root-device-type":                 {"ebs"},
		"block-device-mapping.volume-type": {"gp2"},
	}

	return GetMostRecentAmiIdE(t, region, CanonicalAccountId, filters)
}

// Return a CentOS 7 public AMI from the given region.
// WARNING: you may have to accept the terms & conditions of this AMI in AWS MarketPlace for your AWS Account before
// you can successfully launch the AMI.
func GetCentos7Ami(t *testing.T, region string) string {
	amiId, err := GetCentos7AmiE(t, region)
	if err != nil {
		t.Fatal(err)
	}
	return amiId
}

// Return a CentOS 7 public AMI from the given region.
// WARNING: you may have to accept the terms & conditions of this AMI in AWS MarketPlace for your AWS Account before
// you can successfully launch the AMI.
func GetCentos7AmiE(t *testing.T, region string) (string, error) {
	filters := map[string][]string{
		"name":                             {"*CentOS Linux 7 x86_64 HVM EBS*"},
		"virtualization-type":              {"hvm"},
		"architecture":                     {"x86_64"},
		"root-device-type":                 {"ebs"},
		"block-device-mapping.volume-type": {"gp2"},
	}

	return GetMostRecentAmiIdE(t, region, CentOsAccountId, filters)
}

// Return an Amazon Linux AMI HVM, SSD Volume Type public AMI for the given region.
func GetAmazonLinuxAmi(t *testing.T, region string) string {
	amiId, err := GetAmazonLinuxAmiE(t, region)
	if err != nil {
		t.Fatal(err)
	}
	return amiId
}

// Return an Amazon Linux AMI HVM, SSD Volume Type public AMI for the given region.
func GetAmazonLinuxAmiE(t *testing.T, region string) (string, error) {
	filters := map[string][]string{
		"name":                             {"*amzn-ami-hvm-*-x86_64*"},
		"virtualization-type":              {"hvm"},
		"architecture":                     {"x86_64"},
		"root-device-type":                 {"ebs"},
		"block-device-mapping.volume-type": {"gp2"},
	}

	return GetMostRecentAmiIdE(t, region, AmazonAccountId, filters)
}

// Return an Amazon ECS-Optimized Amazon Linux AMI for the given region. This AMI is useful for running an ECS cluster.
func GetEcsOptimizedAmazonLinuxAmi(t *testing.T, region string) string {
	amiId, err := GetEcsOptimizedAmazonLinuxAmiE(t, region)
	if err != nil {
		t.Fatal(err)
	}
	return amiId
}

// Return an Amazon ECS-Optimized Amazon Linux AMI for the given region. This AMI is useful for running an ECS cluster.
func GetEcsOptimizedAmazonLinuxAmiE(t *testing.T, region string) (string, error) {
	filters := map[string][]string{
		"name":                             {"*amzn-ami*amazon-ecs-optimized*"},
		"virtualization-type":              {"hvm"},
		"architecture":                     {"x86_64"},
		"root-device-type":                 {"ebs"},
		"block-device-mapping.volume-type": {"gp2"},
	}

	return GetMostRecentAmiIdE(t, region, AmazonAccountId, filters)
}

type NoImagesFound struct {
	Region  string
	OwnerId string
	Filters map[string][]string
}

func (err NoImagesFound) Error() string {
	return fmt.Sprintf("No AMIs found in %s for owner ID %s and filters: %v", err.Region, err.OwnerId, err.Filters)
}
