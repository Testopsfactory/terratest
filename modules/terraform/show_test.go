package terraform

import (
	"io/ioutil"
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

func TestShowBasicPlanJSON(t *testing.T) {
	t.Parallel()

	planJSON, err := ioutil.ReadFile("../../test/fixtures/terraform-basic-json/plan.json")
	require.NoError(t, err)

	// Unmarshal the plan into golang types for deeper inspection
	planObject := jsonplan.Unmarshal(t, string(planJSON))
	require.Contains(t, planObject, "create")
}
