package terraform

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terraform/jsonplan"
	"github.com/stretchr/testify/require"
)

func TestShow(t *testing.T) {
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
	require.Contains(t, out, "1 to add, 0 to change, 0 to destroy.")
	require.Contains(t, out, "This plan was saved to: "+options.Out)
	require.FileExistsf(t, options.Out, "Plan file was not created")

	options = &Options{
		TerraformDir: testFolder,
		Out:          testFolder + "/plan.out",
	}

	// Test the JSON string
	planJSON := Show(t, options)
	require.Contains(t, planJSON, "null_resource.test[0]")

	// Unmarshal the plan into golang types for deeper inspection
	planObject := jsonplan.Unmarshal(t, planJSON)
	resourceChanges := planObject.ResourceChanges

	for _, resourceChange := range resourceChanges {
		require.Contains(t, resourceChange.Change.Actions, "create")
	}
}
