package aws

import (
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"testing"
	"github.com/aws/aws-sdk-go/aws"
)

// Get the IDs of EC2 Instances in the given ASG
func GetInstanceIdsForAsg(t *testing.T, asgName string, awsRegion string) []string {
	ids, err := GetInstanceIdsForAsgE(t, asgName, awsRegion)
	if err != nil {
		t.Fatal(err)
	}
	return ids
}

// Get the IDs of EC2 Instances in the given ASG
func GetInstanceIdsForAsgE(t *testing.T, asgName string, awsRegion string) ([]string, error) {
	asgClient, err := NewAsgClientE(t, awsRegion)
	if err != nil {
		return nil, err
	}

	input := autoscaling.DescribeAutoScalingGroupsInput{AutoScalingGroupNames: []*string{aws.String(asgName)}}
	output, err := asgClient.DescribeAutoScalingGroups(&input)
	if err != nil {
		return nil, err
	}

	instanceIds := []string{}
	for _, asg := range output.AutoScalingGroups {
		for _, instance := range asg.Instances {
			instanceIds = append(instanceIds, aws.StringValue(instance.InstanceId))
		}
	}

	return instanceIds, nil
}

// Create an Auto Scaling Group client
func NewAsgClient(t *testing.T, region string) *autoscaling.AutoScaling {
	client, err := NewAsgClientE(t, region)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// Create an Auto Scaling Group client
func NewAsgClientE(t *testing.T, region string) (*autoscaling.AutoScaling, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return autoscaling.New(sess), nil
}

