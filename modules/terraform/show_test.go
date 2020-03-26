package terraform

import (
	"io/ioutil"
	"net"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terraform/jsonplan"
	"github.com/stretchr/testify/require"
)

func TestShowWithInlinePlan(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-basic-configuration", t.Name())
	require.NoError(t, err)

	options := &Options{
		TerraformDir: testFolder,
		Out:          testFolder + "/plan.out",
		Vars: map[string]interface{}{
			"cnt": 1,
		},
	}

	out := InitAndPlan(t, options)
	require.Contains(t, out, "This plan was saved to: "+options.Out)
	require.FileExistsf(t, options.Out, "Plan file was not created")

	// show command does not accept Vars
	options = &Options{
		TerraformDir: testFolder,
		Out:          testFolder + "/plan.out",
	}

	// Test the JSON string
	planJSON := Show(t, options)
	require.Contains(t, planJSON, "null_resource.test[0]")

	// Unmarshal the plan into golang types for deeper inspection
	planObject := jsonplan.Unmarshal(t, planJSON)

	for _, resourceChange := range planObject.ResourceChanges {
		require.Contains(t, resourceChange.Change.Actions, "create")
	}
}

func TestShowBasicJSON(t *testing.T) {
	t.Parallel()

	planJSON, err := ioutil.ReadFile("../../test/fixtures/terraform-basic-json/plan.json")
	require.NoError(t, err)

	// Unmarshal the plan into golang types for deeper inspection
	planObject := jsonplan.Unmarshal(t, string(planJSON))

	plannedResources := planObject.PlannedValues.RootModule.Resources

	for _, plannedResource := range plannedResources {
		if plannedResource.Type == "local_file" {
			require.Equal(t, plannedResource.AttributeValues["file_permission"], "0777", "%s file_permission incorrect", plannedResource.Address)
		}
		if plannedResource.Address == "local_file.example" {
			require.Equal(t, plannedResource.AttributeValues["content"], "example + test")
		} else if plannedResource.Address == "local_file.example2" {
			require.Equal(t, plannedResource.AttributeValues["content"], "test")
		}
	}
}

func TestShowAwsJSON(t *testing.T) {
	t.Parallel()

	planJSON, err := ioutil.ReadFile("../../test/fixtures/terraform-aws-json/plan.json")
	require.NoError(t, err)

	// Unmarshal the plan into golang types for deeper inspection
	planObject := jsonplan.Unmarshal(t, string(planJSON))

	plannedResources := planObject.PlannedValues.RootModule.Resources

	allowedInstanceTypes := []string{"t2.micro", "m5.large", "m5.xlarge"}

	for _, plannedResource := range plannedResources {
		if plannedResource.Type == "aws_instance" {
			require.Contains(t, plannedResource.AttributeValues["tags"], "Name")
			require.Contains(t, plannedResource.AttributeValues["ami"], "ami-")
			require.Contains(t, allowedInstanceTypes, plannedResource.AttributeValues["instance_type"])
		}
	}
}

func TestShowAwsEcsJSON(t *testing.T) {
	t.Parallel()

	planJSON, err := ioutil.ReadFile("../../test/fixtures/terraform-aws-ecs-json/plan.json")
	require.NoError(t, err)

	// Unmarshal the plan into golang types for deeper inspection
	planObject := jsonplan.Unmarshal(t, string(planJSON))

	plannedResources := planObject.PlannedValues.RootModule.Resources

	plannedResourceTypes := []string{}

	for _, plannedResource := range plannedResources {
		plannedResourceTypes = append(plannedResourceTypes, plannedResource.Type)
		if plannedResource.Type == "aws_iam_role" {
			require.Contains(t, plannedResource.AttributeValues["assume_role_policy"], "sts:AssumeRole")
			require.NotContains(t, plannedResource.AttributeValues["assume_role_policy"], "iam:")
		}
	}

	require.Subset(t, plannedResourceTypes, []string{"aws_ecs_cluster", "aws_iam_role", "aws_iam_role_policy_attachment"})
}

func TestShowAwsNetworkJSON(t *testing.T) {
	t.Parallel()

	planJSON, err := ioutil.ReadFile("../../test/fixtures/terraform-aws-network-json/plan.json")
	require.NoError(t, err)

	// Unmarshal the plan into golang types for deeper inspection
	planObject := jsonplan.Unmarshal(t, string(planJSON))

	plannedResources := planObject.PlannedValues.RootModule.Resources

	plannedResourceTypes := []string{}

	_, vpcSubnet, _ := net.ParseCIDR("10.10.0.0/16")
	var subnetIP net.IP

	for _, plannedResource := range plannedResources {
		plannedResourceTypes = append(plannedResourceTypes, plannedResource.Type)

		if plannedResource.Type == "aws_subnet" {

			subnetIP, _, _ = net.ParseCIDR(plannedResource.AttributeValues["cidr_block"].(string))
			require.True(t, vpcSubnet.Contains(subnetIP))

			require.Contains(t, plannedResource.AttributeValues["tags"], "Name")

			if plannedResource.Name == "private" {
				require.Equal(t, plannedResource.AttributeValues["map_public_ip_on_launch"], false)
			} else if plannedResource.Name == "public" {
				require.Equal(t, plannedResource.AttributeValues["map_public_ip_on_launch"], true)
			}
		}
	}

	require.Subset(t, plannedResourceTypes, []string{"aws_vpc", "aws_subnet", "aws_nat_gateway", "aws_internet_gateway"})
}
