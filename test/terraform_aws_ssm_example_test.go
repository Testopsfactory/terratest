package test

import (
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

func TestTerraformAwsSsmExample(t *testing.T) {
	t.Parallel()
	region := aws.GetRandomStableRegion(t, []string{}, []string{})

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/terraform-aws-ssm-example",
		Vars: map[string]interface{}{
			"region": region,
		},
	}
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	instanceID := terraform.Output(t, terraformOptions, "instance_id")
	timeout, _ := time.ParseDuration("3m")

	aws.WaitForSsmInstance(t, region, instanceID, timeout)

	stdout, stderr := aws.CheckSsmCommand(t, region, instanceID, "echo Hello, World", timeout)
	require.Equal(t, stdout, "Hello, World\n")
	require.Equal(t, stderr, "")

	_, _, err := aws.CheckSsmCommandE(t, region, instanceID, "false", timeout)
	require.Error(t, err)
	require.Equal(t, err.Error(), "Failed")
}
