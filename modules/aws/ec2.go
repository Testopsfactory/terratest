package aws

import (
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Create an EC2 client
func NewEc2Client(region string) (*ec2.EC2, error) {
	sess, err := GetAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return ec2.New(sess), nil
}
