package terraform

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// InitAndPlan runs terraform init and plan with the given options and return stdout/stderr from the apply command.
func InitAndPlan(t *testing.T, options *Options) int {
	exitCode, err := InitAndPlanE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return exitCode
}

// InitAndPlanE runs terraform init and plan with the given options and return stdout/stderr from the apply command.
func InitAndPlanE(t *testing.T, options *Options) (int, error) {
	if _, err := InitE(t, options); err != nil {
		return -1, err
	}

	if _, err := GetE(t, options); err != nil {
		return -1, err
	}

	return PlanExitCodeE(t, options)
}

// PlanExitCode runs terraform apply with the given options and returns the detailed exitcode.
func PlanExitCode(t *testing.T, options *Options) int {
	exitCode, err := PlanExitCodeE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return exitCode
}

// PlanExitCodeE runs terraform apply with the given options and returns the detailed exitcode.
func PlanExitCodeE(t *testing.T, options *Options) (int, error) {
	exitCode, err := GetExitCodeForTerraformCommandE(t, options, FormatArgs(options.Vars, "plan", "-input=false", "-lock=false", "-detailed-exitcode")...)
	if err != nil {
		return -1, err
	}
	return exitCode, nil
}

// AssertPlanHasNoChanges run terraform plan with detailed-exitcode and asserts for no changes
func AssertPlanHasNoChanges(t *testing.T, options *Options) {
	t.Parallel()

	exitCode := InitAndPlan(t, options)
	assert.Equal(t, 0, exitCode)
}

// AssertPlanHasChanges run terraform plan with detailed-exitcode and asserts for changes
func AssertPlanHasChanges(t *testing.T, options *Options) {
	t.Parallel()

	exitCode := InitAndPlan(t, options)
	assert.Equal(t, 2, exitCode)
}
