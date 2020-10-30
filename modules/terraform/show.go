package terraform

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// Show runs terraform show with the given options and returns stdout/stderr.
// This will fail the test if there is an error in the command.
func Show(t testing.TestingT, options *Options) string {
	out, err := ShowE(t, options)
	require.NoError(t, err)
	return out
}

// ShowE calls terraform show for the given plan output file and returns it.
func ShowE(t testing.TestingT, options *Options) (string, error) {
	// We can only run show if PlanFilePath is set.
	if options.PlanFilePath == "" {
		return "", PlanFilePathRequired
	}

	// We manually construct the args here instead of using `FormatArgs`, because show only accepts a limited set of
	// args.
	return RunTerraformCommandE(t, options, "show", "-no-color", "-json", options.PlanFilePath)
}
