package test

import (
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestTerraformAwsSsmExample(t *testing.T) {
	t.Parallel()
	region := "us-east-2"

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/terraform-aws-ssm-example",
	}
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	instanceID := terraform.Output(t, terraformOptions, "instance_id")
	timeout, _ := time.ParseDuration("3m")

	// website::tag::1:: Wait for the instance to appear in the SSM catalog
	aws.WaitForSsmInstance(t, region, instanceID, timeout)

	// website::tag::2:: Run a command and check its result
	stdout, stderr := aws.CheckSsmCommand(t, region, instanceID, "echo Hello, World", timeout)
	if stdout != "Hello, World\n" {
		t.Fatalf("Wrong value for stdout: %q", stdout)
	}
	if stderr != "" {
		t.Fatalf("Wrong value for stderr: %q", stderr)
	}

	// website::tag::3:: Run a command and get the error
	_, _, err := aws.CheckSsmCommandE(t, region, instanceID, "false", timeout)
	if err.Error() != "Failed" {
		t.Fatalf("Wrong value for error: %q", err.Error())
	}
}
