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
	timeout := 3 * time.Minute

	aws.WaitForSsmInstance(t, region, instanceID, timeout)

	stdout, stderr := aws.CheckSsmCommand(t, region, instanceID, "echo Hello, World", timeout)
	require.Equal(t, stdout, "Hello, World\n")
	require.Equal(t, stderr, "")

	stdout, stderr, err := aws.CheckSsmCommandE(t, region, instanceID, "cat /wrong/file", timeout)
	require.Error(t, err)
	require.Equal(t, err.Error(), "Failed")
	require.Equal(t, "cat: /wrong/file: No such file or directory\nfailed to run commands: exit status 1", stderr)
	require.Equal(t, "", stdout)
}
