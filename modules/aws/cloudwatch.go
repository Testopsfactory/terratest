package aws

import (
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/aws"
	"testing"
)

// Return the CloudWatch log messages in the given region for the given log stream and log group
func GetCloudWatchLogEntries(t *testing.T, awsRegion string, logStreamName string, logGroupName string) []string {
	out, err := GetCloudWatchLogEntriesE(t, awsRegion, logStreamName, logGroupName)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// Return the CloudWatch log messages in the given region for the given log stream and log group
func GetCloudWatchLogEntriesE(t *testing.T, awsRegion string, logStreamName string, logGroupName string) ([]string, error) {
	client, err := NewCloudWatchLogsClient(awsRegion)
	if err != nil {
		return nil, err
	}

	output, err := client.GetLogEvents(&cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreamName),
	})

	if err != nil {
		return nil, err
	}

	entries := []string{}
	for _, event := range output.Events {
		entries = append(entries, *event.Message)
	}

	return entries, nil
}

// Create a new CloudWatch Logs client
func NewCloudWatchLogsClient(region string) (*cloudwatchlogs.CloudWatchLogs, error) {
	sess, err := GetAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return cloudwatchlogs.New(sess), nil
}
