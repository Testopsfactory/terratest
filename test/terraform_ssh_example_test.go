package test

import (
	"testing"
	"github.com/gruntwork-io/terratest/terraform"
	"fmt"
	"github.com/gruntwork-io/terratest/util"
	"github.com/gruntwork-io/terratest/aws"
	"github.com/gruntwork-io/terratest/http"
	"time"
	"github.com/gruntwork-io/terratest/ssh"
)

// An example of how to test the Terraform module in examples/terraform-ssh-example using Terratest.
func TerraformSshExampleTest(t *testing.T) {
	t.Parallel()

	// A unique ID we can use to namespace resources so we don't clash with anything already in the AWS account or
	// tests running in parallel
	uniqueId := util.UniqueId()

	// Give this EC2 Instance and other resources in the Terraform code a name with a unique ID so it doesn't clash
	// with anything else in the AWS account.
	instanceName := fmt.Sprintf("terratest-ssh-example-%s", uniqueId)

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := aws.PickRandomRegion(t)

	// Create an EC2 KeyPair that we can use for SSH access
	keyPairName := fmt.Sprintf("terratest-ssh-example-%s", uniqueId)
	keyPair := aws.CreateEC2KeyPair(t, awsRegion, keyPairName)

	terraformOptions := terraform.Options {
		// The path to where our Terraform code is located
		TerraformDir: "../examples/terraform-ssh-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]string {
			"aws_region":    awsRegion,
			"instance_name": instanceName,
			"key_pair_name": keyPairName,
		},
	}

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.Apply(t, terraformOptions)

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform output` to get the value of an output variable
	instanceIp := terraform.Output(t, terraformOptions, "instance_ip")

	// We're going to try to SSH to the instance IP, using the Key Pair we created earlier, and the user "ubuntu",
	// as we know the Instance is running an Ubuntu AMI that has such a user
	host := ssh.Host{
		Hostname: instanceIp,
		SshKeyPair: keyPair,
		SshUserName: "ubuntu",
	}

	// It can take a minute or so for the Instance to boot up, so retry a few times
	maxRetries := 15
	timeBetweenRetries := 5 * time.Second
	description := fmt.Sprintf("SSH to %s", instanceIp)

	// Run a simple echo command on the server
	expectedText := "Hello, World"
	command := fmt.Sprintf("echo '%s'", expectedText)

	// Verify that we can SSH to the Instance and run commands
	util.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() error {
		actualText := ssh.CheckSshCommand(t, host, command)

		if actualText != command {
			return fmt.Errorf("Expected SSH command to return '%s' but got '%s'", expectedText, actualText)
		}

		return nil
	})
}


